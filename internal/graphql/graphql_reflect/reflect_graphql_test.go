package graphql_reflect

import (
	"reflect"
	"strings"
	"testing"
)

func buildGraphQLFragment(t reflect.Type) (string, error) {
	qb := &queryBuilder{}
	err := qb.addFragment(t)
	if err != nil {
		return "", err
	}
	return qb.body.String(), nil
}

func TestReflectGraphql(t *testing.T) {
	type Cat struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	type Person struct {
		ID   string `json:"id"`
		Cats []Cat  `json:"cats"`
	}

	frag, err := buildGraphQLFragment(reflect.TypeOf(&Cat{}))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if strings.TrimSpace(frag) != "id name" {
		t.Errorf("bad fragment: %s", frag)
	}

	frag, err = buildGraphQLFragment(reflect.TypeOf(&Person{}))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if strings.TrimSpace(frag) != "id cats { id name }" {
		t.Errorf("bad fragment: %s", frag)
	}
}

func TestReflectGraphQLWithArgs(t *testing.T) {
	type NodeQuery struct {
		Node struct {
			ID       string `json:"id"`
			TypeName string `json:"__typename"`
		} `json:"node" args:"id ID!"`
	}
	frag, err := BuildQuery(reflect.TypeOf(&NodeQuery{}))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	const expected = "query($id: ID!) { node(id: $id) { id __typename } }"
	actual := strings.TrimSpace(frag)
	if actual != expected {
		t.Errorf("fragment does not match:\n\texpected: %s\n\tgot:      %s\n", expected, actual)
	}
}
