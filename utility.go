package main

import (
	"strings"
	"unicode"
)

func lowerTrim(input string) string {
	return strings.ToLower(strings.TrimLeftFunc(input, unicode.IsSpace))
}
