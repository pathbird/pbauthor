package codex

import (
	"context"
	"github.com/pathbird/pbauthor/internal/auth"
	"github.com/pathbird/pbauthor/internal/course"
	"github.com/pathbird/pbauthor/internal/graphql"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Initialize a new codex config file
func InitConfig(dirname string) (*Config, error) {
	// Look for a codex file before initializing
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list files")
	}

	var codexSourceFile string
	for _, file := range files {
		if isCodexSourceFile(file) {
			codexSourceFile = file.Name()
			break
		}
	}

	if codexSourceFile == "" {
		return nil, errors.Errorf("directory (%s) does not contain a codex source file", dirname)
	}

	configFile := filepath.Join(dirname, ConfigFileName)
	conf := &Config{
		configFile: configFile,
	}

	// TODO: shouldn't create a new client here, but oh well
	authn, err := auth.GetAuth()
	if err != nil {
		return nil, err
	}
	g := graphql.NewClient(authn)
	courses, err := g.QueryCourses(context.Background())
	if err != nil {
		return nil, err
	}

	cour, err := course.PromptCourse(courses)
	if err != nil {
		return nil, err
	}

	cat, err := course.PromptCodexCategory(cour.CodexCategories)
	if err != nil {
		return nil, err
	}

	conf.Upload.CodexCategory = cat.ID

	if err := conf.Save(); err != nil {
		return nil, err
	}
	return conf, nil
}

func isCodexSourceFile(file os.FileInfo) bool {
	ext := filepath.Ext(file.Name())
	return !file.IsDir() && ext == ".ipynb"
}
