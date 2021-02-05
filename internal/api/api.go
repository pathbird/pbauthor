package api

import (
	"github.com/pkg/errors"
	"time"
)

const dateTimeFormat = "2006-01-02T15:04:05Z"
func ParseDateTime(input string) (time.Time, error) {
	t, err := time.Parse(dateTimeFormat, input)
	if err != nil {
		return t, errors.Wrap(err, "failed to parse DateTime")
	}
	return t, nil
}
