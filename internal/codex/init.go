package codex

import (
	"context"
	"github.com/mynerva-io/author-cli/internal/auth"
	"github.com/mynerva-io/author-cli/internal/course"
	"github.com/mynerva-io/author-cli/internal/graphql"
)

func initConfig(configFile string) (*Config, error) {
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
