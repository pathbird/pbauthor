package auth

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mynerva-io/author-cli/internal/api"
	"github.com/mynerva-io/author-cli/internal/constants"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"
)

func AuthenticateFromUserInput() (*Auth, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	email = strings.TrimSpace(email)

	fmt.Printf("password: ")
	passwordBytes, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return nil, err
	}
	fmt.Println()
	password := strings.TrimSpace(string(passwordBytes))

	resp, err := apiAuthenticate(email, password)

	if err != nil {
		return nil, err
	}

	expiration, err := api.ParseDateTime(resp.Token.ExpiresAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse expiration time")
	}

	auth := Auth{
		ApiToken:   resp.Token.Token,
		Expiration: expiration,
	}

	err = SaveAuth(&auth)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save auth state")
	}

	return &auth, nil
}

const apiLoginEndpoint = constants.API_HOST + "/auth/login"

type apiLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type apiLoginErrorResponse struct {
	Type        string `json:"error"`
	Description string `json:"description"`
}
type apiLoginTokenResponse struct {
	Token struct {
		UserID    string `json:"userId"`
		IssuedAt  string `json:"issuedAt"`
		ExpiresAt string `json:"expiresAt"`

		// Token is filled in manually by inspecting the cookie
		Token string
	} `json:"token"`
}

func apiAuthenticate(email string, password string) (*apiLoginTokenResponse, error) {
	reqBody, err := json.Marshal(&apiLoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal apiLoginRequest")
	}
	req, err := http.NewRequest("POST", apiLoginEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "mynerva-author-cli")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request to api failed")
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read api response")
	}

	if resp.StatusCode != http.StatusOK {
		var payload apiLoginErrorResponse
		err := json.Unmarshal(respData, &payload)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal api error response")
		}
		return nil, fmt.Errorf("authentication failed: %s (%s)", payload.Type, payload.Description)
	}

	fmt.Println(resp.Header.Get("Set-Cookie"))

	var payload apiLoginTokenResponse
	err = json.Unmarshal(respData, &payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal api response")
	}

	var token string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "auth_token" {
			token = cookie.Value
			break
		}
	}
	if token == "" {
		return nil, fmt.Errorf("api login returned success, but didn't include a token cookie")
	}
	payload.Token.Token = token

	return &payload, nil
}
