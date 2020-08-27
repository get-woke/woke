package printer

import (
	"fmt"

	"github.com/caitlinelfring/woke/pkg/result"
)

type Text struct{}

func NewText() *Text {
	return &Text{}
}

func (t *Text) Print(fs *result.FileResults) error {
	fmt.Println(fs.Filename)

	for _, r := range fs.Results {
		pos := fmt.Sprintf("%d:%d-%d:%d",
			r.StartPosition.Line,
			r.StartPosition.Column,
			r.EndPosition.Line,
			r.EndPosition.Column)

		fmt.Printf("\t%-14s %-10s %s\n", pos, r.Rule.Severity, r.Reason())
	}
	return nil
}
