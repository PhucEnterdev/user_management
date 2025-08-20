package utils

import "strings"

func NormalizeString(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}
