package printer

import (
	"fmt"
	"io"

	"github.com/get-woke/woke/pkg/result"

	"github.com/fatih/color"
)

// Text is a text printer meant for humans to read
type Text struct {
	writer       io.Writer
	disableColor bool
}

// NewText returns a text Printer with color optionally disabled
func NewText(w io.Writer, disableColor bool) *Text {
	return &Text{
		writer:       w,
		disableColor: disableColor,
	}
}

func (t *Text) ShouldSkipExitMessage() bool {
	return false
}

// Print prints the file results
func (t *Text) Print(fs *result.FileResults) error {
	if t.disableColor {
		color.NoColor = true
	}

	var err error
	for _, r := range fs.Results {
		pos := fmt.Sprintf("%d:%d-%d",
			r.GetStartPosition().Line,
			r.GetStartPosition().Column,
			r.GetEndPosition().Column)

		sev := r.GetSeverity()

		_, err = fmt.Fprintf(t.writer, "%s:%s: %s (%s)\n",
			color.New(color.Bold, color.FgHiCyan).Sprint(fs.Filename),
			color.New(color.Bold).Sprint(pos),
			color.New(color.FgHiMagenta).Sprint(r.Reason()),
			sev.Colorize())

		if err != nil {
			return err
		}

		// If the line empty, skip showing the source code
		// This could happen if the line is too long to be worth showing
		if len(r.GetLine()) > 0 {
			_, err = fmt.Fprintln(t.writer, r.GetLine())
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(t.writer, t.arrowUnderLine(r))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Text) Start() {
}

func (t *Text) End() {
}

func (t *Text) arrowUnderLine(r result.Result) string {
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
