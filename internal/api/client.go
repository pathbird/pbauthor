package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	authToken  string
	host       string
	httpClient *http.Client
}

func New(authToken string) *Client {
	return &Client{authToken: authToken, host: "https://mynerva.io/api", httpClient: http.DefaultClient}
}

type request struct {
	route string
	body  interface{}
}

type response struct {
	route        string
	httpResponse *http.Response
}

func (r *response) ReadBytes() ([]byte, error) {
	respData, err := ioutil.ReadAll(r.httpResponse.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read api response")
	}
	return respData, nil
}

// Return an error if the HTTP status code indicates an error (i.e., 4xx or 5xx code).
func (r *response) StatusError() error {
	if r.httpResponse.StatusCode >= 400 {
		return errors.New(fmt.Sprintf("api endpoint (%s) returned error status: %s", r.route, r.httpResponse.Status))
	}
	return nil
}

var _ io.Closer = (*response)(nil)
func (r *response) Close() error {
	return r.httpResponse.Body.Close()
}

func (r *response) UnmarshalJson(target interface{}) error {
	data, err := r.ReadBytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, target)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal json response")
	}
	return nil
}

const userAgent = `mynerva-author-cli`

func (c *Client) postJson(r *request) (*response, error) {
	reqBody, err := json.Marshal(r.body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request")
	}

	endpoint := fmt.Sprintf("%s/%s", c.host, r.route)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	httpResponse, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request to api failed")
	}

	return &response{
		route:        r.route,
		httpResponse: httpResponse,
	}, nil
}

type multipartRequest struct {
	route  string
	fields map[string]string
	files  []FileRef
}

func (c *Client) postMultipart(r *multipartRequest) (*response, error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for fieldname, value := range r.fields {
		err := w.WriteField(fieldname, value)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create api request body")
		}
	}
	for _, file := range r.files {
		// We just always add files underneath the "files" key since that's what every
		// Mynerva API endpoint expects
		err := file.addToWriter("files", w)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create api request body")
		}
	}
	err := w.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create api request body")
	}

	endpoint := fmt.Sprintf("%s/%s", c.host, r.route)
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("User-Agent", userAgent)

	httpResponse, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request to api failed")
	}

	return &response{
		route:        r.route,
		httpResponse: httpResponse,
	}, nil
}
