package main

import (
	"strings"
	"unicode"
)

func lowerTrim(input string) string {
	return strings.ToLower(strings.TrimLeftFunc(input, unicode.IsSpace))
}

func shiftLastCharacter(prefix string) string {
	if len(prefix) > 0 {
		runes := []rune(prefix)
		lastChar := len(runes) - 1
		runes[lastChar] = runes[lastChar] + 1
		return string(runes)
	}
	return ""
}
