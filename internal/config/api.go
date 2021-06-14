package config

import (
	"os"
	"strings"
)

var PathbirdApiHost = (func() string {
	value, set := os.LookupEnv("PATHBIRD_API_HOST")
	if set {
		return strings.TrimRight(value, "/")
	}
	return "https://pathbird.com"
})()
