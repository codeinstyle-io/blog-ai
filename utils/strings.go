package utils

import (
	"strings"
)

func Slugify(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "-"))
}
