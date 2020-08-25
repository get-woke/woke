package parser

import (
	"github.com/caitlinelfring/woke/pkg/ignore"
	"github.com/caitlinelfring/woke/pkg/util"
)

const defaultPath = "."

// Parsable contains the list of files that can be parsed
type Parsable struct {
	Files []string
}

// ParsableFiles returns a list of files that can be parsed after the ignorer has
// excluded files that should be ignored
func ParsableFiles(fileGlobs []string, ignorer *ignore.Ignore) *Parsable {
	if len(fileGlobs) == 0 {
		fileGlobs = []string{defaultPath}
	}

	var p Parsable
	allFiles, _, _ := util.GetFilesInGlobs(fileGlobs)
	for _, f := range allFiles {
		if ignorer.Match(f) {
			continue
		}
		p.Files = append(p.Files, f)
	}

	return &p
}
