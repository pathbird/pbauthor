package graphql

import (
	"context"
	"github.com/pkg/errors"
)

type queryViewerUserRes struct {
	Viewer struct {
		User User `json:"user"`
	} `json:"viewer"`
}

func (c *Client) QueryViewerUser(ctx context.Context) (*User, error) {
	var res queryViewerUserRes
	err := c.queryAndUnmarshall(ctx, &res)
	if err != nil {
		return nil, err
	}
	return &res.Viewer.User, nil
}

type queryCoursesRes struct {
	Viewer struct {
		User struct {
			ID string `json:"id"`
			Courses []CourseEdge `json:"courses" args:"roles [CourseRole!]"`
		} `json:"user"`
	} `json:"viewer"`
}

func (c *Client) QueryCourses(ctx context.Context) ([]CourseEdge, error) {
	var res queryCoursesRes
	err := c.runAndUnmarshall(ctx, &res, withVariable("roles", []string{"OWNER"}))
	if err != nil {
		return nil, err
	}
	if res.Viewer.User.ID == "" {
		// This probably(?) indicates that the user is not authenticated.
		return nil, errors.New("not authenticated")
	}
	return res.Viewer.User.Courses, nil
}