package util

import (
	"fmt"
)

// MarkdownCodify returns a markdown code string
// https://www.markdownguide.org/basic-syntax/#code
func MarkdownCodify(s string) string {
	return fmt.Sprintf("`%s`", s)
}

// InSlice returns true if a string is found in the slice
func InSlice(s string, slice []string) bool {
	if len(slice) == 0 {
		return false
	}
	for _, el := range slice {
		if el == s {
			return true
		}
	}
	return false
}

// ContainsAlphanumeric returns true if alphanumeric chars are found in the string
func ContainsAlphanums(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		char := s[i]
		if ('a' <= char && char <= 'z') ||
			('A' <= char && char <= 'Z') ||
			('0' <= char && char <= '9') {
			return true
		}
	}
	return false
}
