package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"time"
)

type Auth struct {
	ApiToken string `json:"apiToken"`
	Expiration time.Time `json:"expirationTime"`
}

func getApiTokenCacheFile() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to determine auth file: %w", err)
	}
	dir := path.Join(currentUser.HomeDir, ".mynerva")
	_ = os.MkdirAll(dir, 0700)
	return path.Join(currentUser.HomeDir, ".mynerva", "auth.json"), nil
}

var authCache = (*Auth)(nil)

func GetAuth() (*Auth, error) {
	if authCache != nil {
		return authCache, nil
	}

	file, err := getApiTokenCacheFile()
	if err != nil {
		return nil, fmt.Errorf("unable to determine auth file: %w", err)
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, nil
	}

	fp, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("unable to open auth file (%s): %w", file, err)
	}
	defer fp.Close()
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("unable to read auth file (%s): %w", file, err)
	}

	var auth Auth
	err = json.Unmarshal(data, &auth)
	if err != nil {
		return nil, fmt.Errorf("unable to parse auth file (%s): %w", file, err)
	}

	if auth.Expiration.After(time.Now()) {
		if err := os.Remove(file); err != nil {
			return nil, fmt.Errorf("failed to remove expired auth file (%s): %w", file, err)
		}
	}

	authCache = &auth
	return &auth, nil
}

func GetAuthApiToken() (string, error) {
	auth, err := GetAuth()
	if auth == nil || err != nil {
		return "", err
	}
	return auth.ApiToken, nil
}

func SaveAuth(auth *Auth) error {
	file, err := getApiTokenCacheFile()
	if err != nil {
		return err
	}

	b, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, b, 0600)
	if err != nil {
		return fmt.Errorf("failed to write auth file (%s): %w", file, err)
	}

	return nil
}
