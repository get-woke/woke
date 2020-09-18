package result

import (
	"go/token"

	"github.com/get-woke/woke/pkg/rule"
)

type ResultService interface {
	GetSeverity() rule.Severity
	GetStartPosition() *token.Position
	GetEndPosition() *token.Position
	Reason() string
	String() string
	GetLine() string
}
