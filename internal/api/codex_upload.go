package api

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"github.com/djherbis/buffer"
	"github.com/djherbis/nio"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type UploadCodexRequest struct {
	// required
	// the codex category that the uploaded codex should belong to
	CodexCategoryId string `json:"codexCategoryId"`

	// required
	// the files that comprise the codex bundle
	Files []FileRef `json:"-"`

	// optional
	// if specified, replace this codex instead of uploading a new codex
	CodexId string `json:"replaceCodexId,omitempty"`

	KernelOptions KernelOptions `json:"kernelOptions"`
}

type KernelOptions struct {
	SystemPackages []string `json:"systemPackages,omitempty"`
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
	// Do this first so we can bail out early
	codexFile, err := getCodexFile(r.Files)
	if err != nil {
		return nil, nil, err
	}

	buf := buffer.New(1024 * 16)
	pr, pw := nio.Pipe(buf)
	form := multipart.NewWriter(pw)

	// Write the request asynchronously
	writeRequest := func() error {
		defer func() {
			if err := form.Close(); err != nil {
				log.WithError(err).Debug("failed to finalize form for codex upload")
			}
		}()

		requestFormFile, err := form.CreateFormFile("request", "request.json")
		if err != nil {
			return errors.Wrap(err, "initializing upload codex request")
		}
		if err := json.NewEncoder(requestFormFile).Encode(r); err != nil {
			return errors.Wrap(err, "initializing upload codex request")
		}

		codexFormFile, err := form.CreateFormFile("codex", codexFile.Name)
		if err != nil {
			err = errors.Wrap(err, "failed to initialize upload codex request")
			return err
		}
		if err := codexFile.copyTo(codexFormFile); err != nil {
			err = errors.Wrap(err, "failed to initialize upload codex request")
			return err
		}

		tarFormFile, err := form.CreateFormFile("body", "body.tar")
		if err != nil {
			err = errors.Wrap(err, "failed to initialize upload codex request")
			return err
		}

		if err := writeCodexTar(tarFormFile, r.Files, codexFile); err != nil {
			err = errors.Wrap(err, "failed to upload codex files")
			return err
		}

		log.Debugf("wrote all request files for codex upload")
		return nil
	}
	go func() {
		if err := writeRequest(); err != nil {
			// This only a debug since we report the "primary" error that is returned by
			// c.do below.
			log.WithError(err).Debug("failed to write request for codex upload")
		}
	}()

	// Actually make the HTTP request
	contentType := fmt.Sprintf(
		"application/x-pb-multipart-request; boundary=%s",
		form.Boundary(),
	)

	httpReq, err := c.newRequest("POST", "author/upload-codex", contentType, pr)
	if err != nil {
		return nil, nil, errors.Wrap(err, "uploading codex (creating HTTP request)")
	}
	httpRes, err := c.do(httpReq)
	if err != nil {
		return nil, nil, errors.Wrap(err, "uploading codex (HTTP request)")
	}

	// Close the writer pipe.
	// We do this because sometimes the Saturn API returns a response before we've written the entire
	// request body (e.g., if it's returning an HTTP 403, it doesn't need to read the entire tarrball
	// to know that). If that happens we want to stop sending data.
	_ = pw.Close()

	res := &response{httpReq.URL.Path, httpRes}
	time.Sleep(2 * time.Second)
	return parseCodexUploadResponse(res)
}

func getCodexFile(fs []FileRef) (FileRef, error) {
	// TODO:
	// 		We should probably only search the root directory here, but oh well.
	var found FileRef
	for _, f := range fs {
		if filepath.Ext(f.Name) != ".ipynb" {
			continue
		}
		if found.Name != "" {
			return FileRef{}, errors.New("expected to find at most one .ipynb file")
		}
		found = f
	}
	if found.Name == "" {
		return FileRef{}, errors.New("no codex file found (expected one .ipynb file)")
	}
	return found, nil
}

func parseCodexUploadResponse(res *response) (*UploadCodexResponse, *CodexParseFailedError, error) {
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

		case "ErrClientUnsupported":
			log.WithError(statusError).Debug("got ErrClientUnsupported")
			return nil, nil, errors.New(statusError.error.Message)

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

func writeCodexTar(w io.Writer, files []FileRef, exclude FileRef) (retErr error) {
	// The format for this is slightly convoluted.
	// We upload two things (in this order) as form/multipart files:
	// 1. A "request" JSON blob which contains the metadata for the upload (e.g., codex category, etc).
	// 2. The actual body of the upload, which is a tar file.

	tarw := tar.NewWriter(w)
	defer func() {
		if err := tarw.Close(); err != nil {
			// Check if we've already set the error (and use that instead if so)
			if retErr != nil {
				return
			}
			retErr = errors.Wrap(err, "failed to create codex tar archive")
		}
	}()

	// subtract one since we exclude the codex file here
	nFiles := len(files) - 1

	// Upload all the remaining files
	log.Debugf("uploading %d files from codex directory", nFiles)
	for i, f := range files {
		if f.FsPath == exclude.FsPath {
			continue
		}

		log.Debugf("uploading file %q (%d of %d)", f.Name, i+1, nFiles)
		stat, err := os.Stat(f.FsPath)
		if err != nil {
			return errors.Wrap(err, "adding files to codex tar archive")
		}

		// Generate a tar header for! the file
		hdr, err := tar.FileInfoHeader(stat, "")
		if err != nil {
			return errors.Wrap(err, "adding files to codex upload")
		}
		// Use FormatPAX here for more accurate mtimes (otherwise they're truncated to the
		// nearest second with the default USTAR format which can sometimes cause weird
		// issues).
		hdr.Format = tar.FormatPAX

		if err := tarw.WriteHeader(hdr); err != nil {
			return errors.Wrap(err, "adding file to codex tar archive")
		}
		n, err := copyFile(f.FsPath, tarw)
		if err != nil {
			return errors.Wrap(err, "adding file to codex tar archive")
		}
		log.Debugf("wrote %d bytes for file %q", n, f.Name)
	}
	return nil
}

func copyFile(file string, w io.Writer) (int64, error) {
	fd, err := os.Open(file)
	if err != nil {
		return -1, err
	}
	defer fd.Close()
	return io.Copy(w, fd)
}
