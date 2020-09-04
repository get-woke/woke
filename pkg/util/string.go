package util

import "fmt"

// MarkdownCodify returns a markdown code string
// https://www.markdownguide.org/basic-syntax/#code
func MarkdownCodify(s string) string {
	return fmt.Sprintf("`%s`", s)
}

// InSlice returns true if a string is found in the slice
func InSlice(s string, slice []string) bool {
	for _, el := range slice {
		if el == s {
			return true
		}
	}
	return false
}
