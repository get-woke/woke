package printer

import (
	"go/token"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
)

func generateFileResult() *result.FileResults {
	r := result.FileResults{Filename: "foo.txt"}
	r.Results = generateResults(r.Filename)
	return &r
}

func generateResults(filename string) []result.Result {
	r := rule.NewTestRule()
	return []result.Result{
		result.LineResult{
			Rule:      r,
			Violation: "testrule",
			Line:      "this testrule must change",
			StartPosition: &token.Position{
				Filename: filename,
				Offset:   0,
				Line:     1,
				Column:   6,
			},
			EndPosition: &token.Position{
				Filename: filename,
				Offset:   0,
				Line:     1,
				Column:   15,
			},
		},
	}
}

func newPosition(f string, l, c int) *token.Position {
	return &token.Position{
		Filename: f,
		Offset:   0,
		Line:     l,
		Column:   c,
	}
}
