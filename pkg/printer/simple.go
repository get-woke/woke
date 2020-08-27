package printer

import (
	"fmt"

	"github.com/caitlinelfring/woke/pkg/result"
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
		fmt.Printf("%v: [%s] %s",
			r.StartPosition,
			r.Rule.Severity,
			r.Reason())
	}
	return nil
}
