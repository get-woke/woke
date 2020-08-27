package util

import "fmt"

// MarkdownCodify returns a markdown code string
// https://www.markdownguide.org/basic-syntax/#code
func MarkdownCodify(s string) string {
	return fmt.Sprintf("`%s`", s)
}
