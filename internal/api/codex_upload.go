package api

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type UploadCodexRequest struct {
	// required
	// the codex category that the uploaded codex should belong to
	CodexCategoryId string `json:"codexCategoryId"`

	// required
	// the files that comprise the codex bundle
	Files []FileRef

	// optional
	// if specified, replace this codex instead of uploading a new codex
	CodexId string `json:"codexId,omitempty"`
}

type UploadCodexResponse struct {
	CodexId string `json:"codexId"`
}

type CodexParseError struct {
	// The type of the error
	Error          string `json:"error"`
	Message        string `json:"message"`
	SourcePosition string `json:"sourcePosition"`
	SourceInfo     struct {
		SourceContext struct {
			Lines []string `json:"lines"`
		} `json:"sourceContext"`
	} `json:"sourceInfo"`
}

type CodexParseFailedError struct {
	Errors []CodexParseError `json:"errors"`
}

func (e *CodexParseFailedError) Error() string {
	return fmt.Sprintf("failed to parse codex: %d errors occurred", len(e.Errors))
}

var _ error = (*CodexParseFailedError)(nil)

func (c *Client) UploadCodex(
	r *UploadCodexRequest,
) (*UploadCodexResponse, *CodexParseFailedError, error) {
	fields := make(map[string]string)
	fields["codexCategoryId"] = r.CodexCategoryId
	if r.CodexId != "" {
		fields["replaceCodexId"] = r.CodexId
	}
	res, err := c.postMultipart(&multipartRequest{
		route:  "author/upload-codex",
		fields: fields,
		files:  r.Files,
	})
	if err != nil {
		return nil, nil, err
	}
	statusError, err := res.StatusError()
	if err != nil {
		return nil, nil, err
	}
	if statusError != nil {
		switch statusError.error.Error {
		case "MySTParseFailedErr", "CodexASTParseFailedErr":
			var parseError CodexParseFailedError
			if err := json.Unmarshal(statusError.error.Details, &parseError); err != nil {
				return nil, nil, errors.Wrap(err, "failed to unmarshal codex parse error details")
			}
			return nil, &parseError, nil

		case "ErrUnauthenticated":
			log.WithError(statusError).Debug("got ErrUnauthenticated")
			return nil, nil, errors.New("You are not logged in (try running `pbauthor auth login`)")

		case "":
			return nil, nil, errors.Errorf("API returned an unknown error")

		default:
			return nil, nil, errors.Errorf(
				"API returned an error: %s: %s",
				statusError.error.Error,
				statusError.error.Message,
			)
		}
	}
	resp := &UploadCodexResponse{}
	err = res.UnmarshalJson(resp)
	if err != nil {
		return nil, nil, err
	}
	return resp, nil, nil
}
