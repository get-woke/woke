package printer

import (
	"fmt"
	"go/token"

	"github.com/get-woke/woke/pkg/result"
)

// Simple is a simple printer meant for a machine to read
type Simple struct{}

// NewSimple returns a new simple printer
func NewSimple() *Simple {
	return &Simple{}
}

// Print prints in the format 'filename:line:column: message'
// based on golint's output: https://github.com/golang/lint/blob/738671d3881b9731cc63024d5d88cf28db875626/golint/golint.go#L121
func (p *Simple) Print(fs *result.FileResults) error {
	for _, r := range fs.Results {
		fmt.Printf("%v: [%s] %s\n",
			positionString(r.StartPosition),
			r.Rule.Severity,
			r.Reason())
	}
	return nil
}

// positionString is similar to Position.String, but includes the Column
// even if the column is 0
func positionString(pos *token.Position) string {
	s := pos.Filename
	if pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d", pos.Line)
		s += fmt.Sprintf(":%d", pos.Column)
	}
	return s
}
