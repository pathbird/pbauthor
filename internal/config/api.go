package config

import (
	"os"
	"strings"
)

var MynervaApiHost = (func() string {
	value, set := os.LookupEnv("MYNERVA_API_HOST")
	if set {
		return strings.TrimRight(value, "/")
	}
	return "https://mynerva.io"
})()
