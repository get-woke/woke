package result

import (
	"go/token"

	"github.com/get-woke/woke/pkg/rule"
)

// Result is an interface for a finding of a rule
type Result interface {
	GetSeverity() rule.Severity
	GetRuleName() string
	GetStartPosition() *token.Position
	GetEndPosition() *token.Position
	GetTerms() []string
	Reason() string
	String() string
	GetLine() string
}
