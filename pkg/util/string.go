package util

import (
	"fmt"
	"strings"
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

// given a absolute path (example: /runner/folder/root/data/effx.yaml)
// and given a workingDir (example: root)
// it will return data/effx.yaml
func ParseRelativePathFromAbsolutePath(absoluteDir, workingDir string) string {
	// if working directory does not end with a slash, add it
	if strings.LastIndex(workingDir, "/") != len(workingDir)-1 {
		workingDir += "/"
	}

	res := strings.Split(absoluteDir, workingDir)
	if len(res) > 1 {
		return res[1]
	}
	return ""
}
