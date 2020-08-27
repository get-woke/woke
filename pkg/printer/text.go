package printer

import (
	"fmt"

	"github.com/caitlinelfring/woke/pkg/result"
	"github.com/fatih/color"
)

// Text is a text printer meant for humans to read
type Text struct {
	enableColor bool
}

// NewText returns a text Printer with color optionally disabled
func NewText(enableColor bool) *Text {
	return &Text{
		enableColor: enableColor,
	}
}

// Print prints the file results
func (t *Text) Print(fs *result.FileResults) error {
	color.NoColor = !t.enableColor
	color.New(color.Underline, color.Bold).Println(fs.Filename)

	for _, r := range fs.Results {
		pos := fmt.Sprintf("%d:%d-%d:%d",
			r.StartPosition.Line,
			r.StartPosition.Column,
			r.EndPosition.Line,
			r.EndPosition.Column)
		sev := r.Rule.Severity.Colorize()
		fmt.Printf("\t%-14s %-20s %s\n", pos, sev, r.Reason())
	}
	fmt.Println()
	return nil
}
