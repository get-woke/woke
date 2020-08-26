package printer

import (
	"fmt"

	"github.com/caitlinelfring/woke/pkg/result"
)

// Simple
type Simple struct{}

func NewSimple() *Simple {
	return &Simple{}
}

// Print prints in the format 'filename:line:column: message'
func (p *Simple) Print(fs *result.FileResults) error {
	var err error
	if _, err = fmt.Println(fs.Filename); err != nil {
		return err
	}

	for _, r := range fs.Results {
		out := fmt.Sprintf("%s:%d: [%s] %s",
			r.Filename,
			r.StartPosition.Line,
			r.Rule.Severity,
			r.Reason())

		if _, err = fmt.Print(out); err != nil {
			return err
		}
	}
	return nil
}
