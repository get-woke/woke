package printer

import (
	"fmt"
	"os"

	"github.com/caitlinelfring/woke/pkg/result"
)

type Text struct{}

func NewText() *Text {
	return &Text{}
}

func (t *Text) Print(fs *result.FileResults) error {
	var err error
	if _, err = fmt.Fprintln(os.Stdout, fs.Filename); err != nil {
		return err
	}
	for _, r := range fs.Results {
		pos := fmt.Sprintf("%s-%s",
			r.StartPosition.String(),
			r.EndPosition.String())

		if _, err = fmt.Fprintf(os.Stdout, "    %-14s %-10s %s\n", pos, r.Rule.Severity, r.Reason()); err != nil {
			return err
		}
	}
	return nil
}
