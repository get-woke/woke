package result

import (
	"strings"
)

type FileResults struct {
	Filename string
	Results  []Result
}

func (fr *FileResults) String() string {
	lines := []string{fr.Filename}
	for _, r := range fr.Results {
		lines = append(lines, r.String())
	}
	return strings.Join(lines, "\n")
}
