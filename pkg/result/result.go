package result

import (
	"go/token"

	"github.com/get-woke/woke/pkg/rule"
)

type Result interface {
	GetSeverity() rule.Severity
	GetStartPosition() *token.Position
	GetEndPosition() *token.Position
	Reason() string
	String() string
	GetLine() string
}
