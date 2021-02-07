package graphql

import (
	"context"
	"fmt"
	"github.com/mynerva-io/author-cli/internal/auth"
	"github.com/mynerva-io/author-cli/internal/config"
	"github.com/mynerva-io/author-cli/internal/graphql/graphql_reflect"
	"github.com/mynerva-io/author-cli/internal/graphql/transport"
	"github.com/pkg/errors"
	"reflect"
)

type Client struct {
	*transport.Client
	auth *auth.Auth
}

func NewClient(auth *auth.Auth) Client {
	client := transport.NewClient(
		fmt.Sprintf("%s/graphql", config.MynervaApiHost),
	)
	client.Log = func(s string) {
		fmt.Printf("graphql client: %s\n", s)
	}
	return Client{
		Client: client,
		auth: auth,
	}
}

func (c *Client) Run(ctx context.Context, req *transport.Request, res interface{}) error {
	if c.auth != nil && c.auth.ApiToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.auth.ApiToken))
	} else {
		return errors.New("authorization not set")
	}
	return c.Client.Run(ctx, req, res)
}

func (c *Client) queryAndUnmarshall(ctx context.Context, v interface{}) error {
	frag, err := graphql_reflect.BuildQuery(reflect.TypeOf(v))
	if err != nil {
		return errors.Wrap(err, "couldn't build GraphQL fragment")
	}
	req := transport.NewRequest( frag)
	return c.Run(ctx, req, v)
}

type runOption func (req *transport.Request)

func withVariable(name string, value interface{}) runOption {
	return func(req *transport.Request) {
		req.Var(name, value)
	}
}

func (c *Client) runAndUnmarshall(ctx context.Context, v interface{}, opts ...runOption) error {
	queryString, err := graphql_reflect.BuildQuery(reflect.TypeOf(v))
	if err != nil {
		return err
	}
	req := transport.NewRequest(queryString)
	for _, opt := range opts {
		opt(req)
	}
	return c.Run(ctx, req, v)
}