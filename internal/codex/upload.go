package codex

import (
	"github.com/pathbird/pbauthor/internal/api"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type UploadCodexOptions struct {
	// The codex directory
	Dir string
}

func UploadCodex(client *api.Client, opts *UploadCodexOptions) (*api.UploadCodexResponse, *api.CodexParseFailedError, error) {
	config, err := GetOrInitCodexConfig(opts.Dir)
	if err != nil {
		return nil, nil, err
	}

	files, err := getCodexFiles(config, opts.Dir)
	if err != nil {
		return nil, nil, err
	}
	log.Debugf("got %d codex files", len(files))

	// TODO:
	// 		We should request a confirmation before doing the upload.
	//		This will help make sure the author is aware of what course
	//		they're uploading to (e.g., in case they try to re-upload
	//		an old codex and intend to upload it to a new course but
	//		but the config still points to the old course).

	req := &api.UploadCodexRequest{
		CodexCategoryId: config.Upload.CodexCategory,
		Files:           files,
		CodexId:         config.Upload.CodexId,
	}

	res, parseErr, err := client.UploadCodex(req)
	if parseErr != nil {
		return nil, parseErr, nil
	}
	if err != nil {
		return nil, nil, err
	}

	config.Upload.CodexId = res.CodexId
	if err := config.Save(); err != nil {
		return nil, nil, errors.Wrap(err, "codex upload succeeded, but failed to save codex config file")
	}

	return res, nil, nil
}

const maxFiles = 20

// Get all the files associated with the codex.
// Recursively walks the filesystem starting at `dir`.
func getCodexFiles(_ *Config, dir string) ([]api.FileRef, error) {
	var files []api.FileRef
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		isHidden := strings.HasPrefix(info.Name(), ".")

		if info.IsDir() {
			// Don't recurse into hidden directories
			if isHidden {
				return filepath.SkipDir
			}
			// For non-hidden directories, we'll still recurse into all the files
			// but we don't need to do anything with the directory itself.
			return nil
		}

		if isHidden {
			return nil
		}
		relpath, err := filepath.Rel(dir, path)
		if err != nil {
			return errors.Wrapf(err, "couldn't determine relative file path: %s", path)
		}

		if len(files) > maxFiles {
			return errors.Errorf("too many codex files (exceeds limit: %d)", maxFiles)
		}
		files = append(files, api.FileRef{
			Name:   relpath,
			FsPath: path,
		})

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "unable to build codex file list")
	}

	return files, nil
}
