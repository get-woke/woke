package printer

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/get-woke/woke/pkg/result"
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

	for _, r := range fs.Results {
		pos := fmt.Sprintf("%d:%d-%d",
			r.GetStartPosition().Line,
			r.GetStartPosition().Column,
			r.GetEndPosition().Column)
		sev := r.GetSeverity()
		fmt.Printf("%s:%s: %s (%s)\n",
			color.New(color.Bold, color.FgHiCyan).Sprint(fs.Filename),
			color.New(color.Bold).Sprint(pos),
			color.New(color.FgHiMagenta).Sprint(r.Reason()),
			sev.Colorize())

		// If the line empty, skip showing the source code
		// This could happen if the line is too long to be worth showing
		if len(r.GetLine()) > 0 {
			fmt.Println(r.GetLine())
			fmt.Printf("%s\n", t.arrowUnderLine(r))
		}
	}

	return nil
}

func (t *Text) arrowUnderLine(r result.ResultService) string {
	// if columns == 0 it means column is unknown
	if r.GetStartPosition().Column == 0 && r.GetEndPosition().Column == 0 {
		return ""
	}

	line := r.GetLine()
	prefix := make([]rune, 0, len(line))

	for i := 0; i < len(line) && i < r.GetStartPosition().Column; i++ {
		if line[i] == '\t' {
			prefix = append(prefix, '\t')
		} else {
			prefix = append(prefix, ' ')
		}
	}

	return fmt.Sprintf("%s%s", string(prefix), color.YellowString("^"))
}
