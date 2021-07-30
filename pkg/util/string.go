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

func FilterEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
