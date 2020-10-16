package result

import (
	"go/token"

	"github.com/get-woke/woke/pkg/rule"
)

// Result is an interface for a violation of a rule
type Result interface {
	GetSeverity() rule.Severity
	GetStartPosition() *token.Position
	GetEndPosition() *token.Position
	Reason() string
	String() string
	GetLine() string
}
