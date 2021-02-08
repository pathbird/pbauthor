package graphql_reflect

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

var fragmentCache = make(map[reflect.Type]string)

type BuildQueryOpt func(builder *queryBuilder)

func WithQueryName(name string) BuildQueryOpt {
	return func(builder *queryBuilder) {

	}
}

func BuildQuery(t reflect.Type, opts ...BuildQueryOpt) (string, error) {
	qb := &queryBuilder{
		operation: Query,
	}
	for _, opt := range opts {
		opt(qb)
	}
	err := qb.addFragment(t)
	if err != nil {
		return "", err
	}
	outerSb := strings.Builder{}
	if qb.operation == Mutation {
		outerSb.WriteString("mutation")
	} else {
		outerSb.WriteString("query")
	}
	if qb.queryName != "" {
		outerSb.WriteString(" ")
		outerSb.WriteString(qb.queryName)
	}
	if len(qb.vars) > 0 {
		outerSb.WriteString("(")
		for i, v := range qb.vars {
			if i > 0 {
				outerSb.WriteString(", ")
			}
			outerSb.WriteString("$")
			outerSb.WriteString(v.name)
			outerSb.WriteString(": ")
			outerSb.WriteString(v.graphQLType)
		}
		outerSb.WriteString(")")
	}
	outerSb.WriteString(" {")
	outerSb.WriteString(qb.body.String())
	outerSb.WriteString(" }")
	return outerSb.String(), nil
}

type graphQLVariable struct {
	name        string
	graphQLType string
}

type OperationType string

const (
	Query    OperationType = "query"
	Mutation OperationType = "mutation"
)

type queryBuilder struct {
	operation OperationType
	queryName string
	body      strings.Builder
	vars      []graphQLVariable
}

func (qb *queryBuilder) addFragment(t reflect.Type) error {
	kind := t.Kind()
	if kind == reflect.Ptr || kind == reflect.Slice {
		t = t.Elem()
	}

	kind = t.Kind()
	if kind != reflect.Struct && kind != reflect.Array {
		return errors.Errorf("expected struct (not %s)", kind)
	}

	sb := &qb.body
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// Only add tags with json tags (otherwise we consider it "private")
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		sb.WriteString(" ")
		switch f.Type.Kind() {
		case reflect.Struct, reflect.Array, reflect.Ptr, reflect.Slice:
			// Maybe TODO:
			//		We don't currently support embedded structs. It's not a huge deal,
			//		and it's not clear how to actually do that, but it would be nice to
			//		be able to use them like spreading fragments in GraphQL
			sb.WriteString(jsonTag)
			if err := qb.writeArgs(sb, &f); err != nil {
				return err
			}
			sb.WriteString(" {")
			err := qb.addFragment(f.Type)
			if err != nil {
				return err
			}
			sb.WriteString(" }")
		case reflect.String, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Int, reflect.Int64:
			sb.WriteString(jsonTag)
		default:
			return errors.Errorf("invalid struct field type: %s.%s", f.Type.Name(), f.Type.Kind())
		}
	}

	return nil
}

func (qb *queryBuilder) writeArgs(sb *strings.Builder, field *reflect.StructField) error {
	tag := field.Tag.Get("args")
	if tag == "" {
		return nil
	}
	args := strings.Split(tag, ",")
	if len(args) == 0 {
		return nil
	}

	sb.WriteString("(")
	for argIndex, argString := range args {
		if argIndex > 0 {
			sb.WriteString(", ")
		}
		arg, err := parseArgString(argString)
		if err != nil {
			return err
		}
		log.Debugf("adding var: %s (%s)", arg.name, arg.graphQLType)
		qb.vars = append(qb.vars, *arg)
		sb.WriteString(fmt.Sprintf("%s: $%s", arg.name, arg.name))
	}
	sb.WriteString(")")
	return nil
}

func parseArgString(s string) (*graphQLVariable, error) {
	var name, vartype string
	_, err := fmt.Sscanf(s, "%s%s", &name, &vartype)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid GraphQL arg annotation: %s", s)
	}
	return &graphQLVariable{
		name:        name,
		graphQLType: vartype,
	}, nil
}
