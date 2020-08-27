package result

import (
	"fmt"
	"go/token"

	"github.com/caitlinelfring/woke/pkg/rule"
)

// Result contains data about the result of a broken rule
type Result struct {
	Rule          *rule.Rule
	Match         string
	StartPosition *token.Position
	EndPosition   *token.Position
}

// FindResults returns the results that match the rule for the given text.
// filename and line are only used for the Position
func FindResults(r *rule.Rule, filename, text string, line int) (rs []Result) {
	idxs := r.Regexp.FindAllStringIndex(text, -1)

	for _, idx := range idxs {
		start := idx[0]
		end := idx[1]
		newResult := Result{
			Rule:  r,
			Match: text[start:end],
			StartPosition: &token.Position{
				Filename: filename,
				Line:     line,
				Column:   start,
			},
			EndPosition: &token.Position{
				Filename: filename,
				Line:     line,
				Column:   end,
			},
		}

		rs = append(rs, newResult)
	}
	return
}

// Reason outputs the suggested alternatives for this rule
func (r *Result) Reason() string {
	return fmt.Sprintf("Instead of '%s', consider the following alternative(s): '%s'", r.Match, r.Rule.Alternatives)
}

func (r *Result) String() string {
	pos := fmt.Sprintf("%s-%s",
		r.StartPosition.String(),
		r.EndPosition.String())
	return fmt.Sprintf("    %-14s %-10s %s", pos, r.Rule.Severity, r.Reason())
}
