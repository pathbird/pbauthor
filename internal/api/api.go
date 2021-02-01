package api

import (
	"fmt"
	"time"
)

const dateTimeFormat = "2006-01-02T15:04:05Z"
func ParseDateTime(input string) (time.Time, error) {
	t, err := time.Parse(dateTimeFormat, input)
	if err != nil {
		return t, fmt.Errorf("failed to parse DateTime: %w", err)
	}
	return t, nil
}
