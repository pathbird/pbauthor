package api

import (
	"github.com/pkg/errors"
	"mime"
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

func isJsonContentType(contentType string) bool {
	// Ignore the possible error here, since error implies *not* JSON
	mediatype, _, _ := mime.ParseMediaType(contentType)
	return mediatype == "application/json"
}
