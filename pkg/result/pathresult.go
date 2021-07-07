package result

import (
	"path/filepath"
	"strings"

	"github.com/get-woke/woke/pkg/rule"
)

// PathResult is a Result meant for showing findings in a file path
type PathResult struct {
	LineResult
}

// Reason is the reason for the PathResult finding.
// It is similar to Result.Reason, but makes it clear that the finding is
// with the file path and not a line in the file
func (r PathResult) Reason() string {
	return "Filename finding: " + r.Rule.ReasonWithNote(r.LineResult.Finding)
}

// MatchPathRules will match the path against all the rules provided
func MatchPathRules(rules []*rule.Rule, path string) (rs []PathResult) {
	for _, r := range rules {
		rs = append(rs, MatchPath(r, path)...)
	}
	return
}

// MatchPath matches the path against the rule. If it is a match, it will
// return a PathResult with the line/start column/end column all at 1
func MatchPath(r *rule.Rule, path string) (rs []PathResult) {
	path = filepath.ToSlash(path)
	dir, filename := filepath.Split(path)
	dirParts := append(filepath.SplitList(dir), strings.TrimSuffix(filename, filepath.Ext(filename)))

	for _, p := range dirParts {
		p = filepath.ToSlash(p)
		if len(r.FindMatchIndexes(p)) > 0 {
			rs = append(rs, PathResult{LineResult: NewLineResult(r, p, path, 1, 1, 1)})
		}
	}

	return
}
