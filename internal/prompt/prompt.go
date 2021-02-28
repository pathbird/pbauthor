package prompt

import (
	"fmt"
	"os"
	"strings"
)

func Confirm(msg string) bool {
	if _, err := fmt.Fprintf(os.Stderr, "%s [y/N] ", msg); err != nil {
		panic(err)
	}
	var res string
	_, _ = fmt.Scanln(&res)
	if strings.HasPrefix(strings.ToLower(res), "y") {
		return true
	}
	return false
}
