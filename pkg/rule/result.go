package rule

import (
	"fmt"
	"go/token"
)

type Result struct {
	Rule     *Rule
	Match    string
	Position *token.Position
}

// Reason outputs the suggested alternatives for this rule
func (r *Result) Reason() string {
	return fmt.Sprintf("Instead of '%s', consider the following alternative(s): '%s'", r.Match, r.Rule.Alternatives)
}

func (r *Result) String() string {
	return fmt.Sprintf("[%s] %s", r.Position.String(), r.Reason())
}
