package result

import (
	"path/filepath"
	"strings"

	"github.com/get-woke/woke/pkg/rule"
)

// PathResult is a ResultService meant for showing violations in a file path
type PathResult struct {
	Result
}

// Reason is the reason for the PathResult violation.
// It is similar to Result.Reason, but makes it clear that the violation is
// with the file path and not a line in the file
func (r PathResult) Reason() string {
	return "Filename violation: " + r.Rule.Reason(r.Result.Violation)
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
	dir, filename := filepath.Split(path)
	dirParts := append(filepath.SplitList(dir), strings.TrimSuffix(filename, filepath.Ext(filename)))

	for _, p := range dirParts {
		if r.MatchString(p, false) {
			rs = append(rs, PathResult{Result: NewResult(r, p, path, 1, 1, 1)})
		}
	}

	return
}
