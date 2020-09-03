package printer

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/get-woke/woke/pkg/result"
)

// MaxLineLengthForShowingViolationSource is the max line length that this printer
// will show the source of the violation and the location within the line of the violation.
// Helps avoid consuming the console when minified files contine violations.
const MaxLineLengthForShowingViolationSource = 200

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
			r.StartPosition.Line,
			r.StartPosition.Column,
			r.EndPosition.Column)
		sev := r.Rule.Severity.Colorize()
		fmt.Printf("%s:%s: %s (%s)\n",
			color.New(color.Bold, color.FgHiCyan).Sprint(fs.Filename),
			color.New(color.Bold).Sprint(pos),
			color.New(color.FgHiMagenta).Sprint(r.Reason()),
			sev)

		// If the line is too long, skip showing the source code
		if len(r.Line) > MaxLineLengthForShowingViolationSource {
			continue
		}

		fmt.Println(r.Line)
		fmt.Printf("%s\n", t.arrowUnderLine(&r))
	}

	return nil
}

func (t *Text) arrowUnderLine(r *result.Result) string {
	// if columns == 0 it means column is unknown
	if r.StartPosition.Column == 0 && r.EndPosition.Column == 0 {
		return ""
	}

	line := r.Line
	prefix := make([]rune, 0, len(line))

	for i := 0; i < len(line) && i < r.StartPosition.Column; i++ {
		if line[i] == '\t' {
			prefix = append(prefix, '\t')
		} else {
			prefix = append(prefix, ' ')
		}
	}

	return fmt.Sprintf("%s%s", string(prefix), color.YellowString("^"))
}
