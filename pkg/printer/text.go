package printer

import (
	"fmt"

	"github.com/caitlinelfring/woke/pkg/result"
	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

// Text is a text printer meant for humans to read
type Text struct {
	disableColor bool
}

// NewText returns a text Printer with color optionally disabled
func NewText(disableColor bool) *Text {
	return &Text{
		disableColor: disableColor,
	}
}

// Print prints the file results
func (t *Text) Print(fs *result.FileResults) error {
	if t.disableColor {
		color.NoColor = true
	}
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

	if len(fs.Results) == 0 {
		log.Info().Msg("👏 Great work using inclusive language in your code! Stay woke! 🙌")
	}

	return nil
}
