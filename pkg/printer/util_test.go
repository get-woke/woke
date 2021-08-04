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
	return []result.Result{
		result.LineResult{
			Rule:    &rule.TestRule,
			Finding: "whitelist",                  // wokeignore:rule=whitelist
			Line:    "this whitelist must change", // wokeignore:rule=whitelist
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

func generateSecondFileResult() *result.FileResults {
	r := result.FileResults{Filename: "bar.txt"}
	r.Results = generateSecondResults(r.Filename)
	return &r
}

func generateSecondResults(filename string) []result.Result {
	return []result.Result{
		result.LineResult{
			Rule:    &rule.TestErrorRule,
			Finding: "slave",                       // wokeignore:rule=slave
			Line:    "this slave term must change", // wokeignore:rule=slave
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
