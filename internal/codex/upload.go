package codex

import (
	"github.com/mynerva-io/author-cli/internal/api"
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

func UploadCodex(client *api.Client, opts *UploadCodexOptions) (*api.UploadCodexResponse, error) {
	files, err := getCodexFiles(&Config{}, opts.Dir)
	if err != nil {
		return nil, err
	}
	log.Debugf("got %d codex files", len(files))

	configFilePath := filepath.Join(opts.Dir, "codex.toml")
	if _, err := os.Stat(configFilePath); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("no _codex.toml file found - please create one")
		}
		return nil, errors.Wrap(err, "unable to stat codex.toml")
	}
	config, err := ReadConfigFile(configFilePath)
	if err != nil {
		return nil, err
	}

	req := &api.UploadCodexRequest{
		CodexCategoryId: config.Upload.CodexCategory,
		Files:           files,
	}

	res, err := client.UploadCodex(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

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
