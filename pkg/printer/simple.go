package printer

import (
	"fmt"

	"github.com/get-woke/woke/pkg/result"
)

// Simple is a simple printer meant for a machine to read
type Simple struct {
	format string
}

// NewSimple returns a new simple printer
func NewSimple() *Simple {
	return &Simple{
		format: "%s:%d:%d: [%s] %s",
	}
}

// Print prints in the format 'filename:line:column: message'
// based on golint's output: https://github.com/golang/lint/blob/738671d3881b9731cc63024d5d88cf28db875626/golint/golint.go#L121
func (s *Simple) Print(fs *result.FileResults) error {
	for _, r := range fs.Results {
		fmt.Printf(s.format+"\n",
			r.StartPosition.Filename,
			r.StartPosition.Line,
			r.StartPosition.Column,
			r.Rule.Severity,
			r.Reason())
	}
	return nil
}
