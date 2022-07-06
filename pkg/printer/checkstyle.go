package printer

import (
	"fmt"
	"io"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
)

// Checkstyle is a Checkstyle printer meant for use by a Checkstyle annotation
type Checkstyle struct {
	writer io.Writer
}

// NewGitHubActions returns a new GitHubActions printer
func NewCheckstyle(w io.Writer) *Checkstyle {
	return &Checkstyle{writer: w}
}

func (p *Checkstyle) PrintSuccessExitMessage() bool {
	return true
}

// Print prints in the format for Checkstyle
// https://docs.github.com/en/actions/reference/workflow-commands-for-github-actions#setting-an-error-message
func (p *Checkstyle) Print(fs *result.FileResults) error {
	fmt.Fprintf(p.writer, "  <file name=\"%s\">\n", fs.Filename)
	for _, r := range fs.Results {
		fmt.Fprintln(p.writer, formatResultForCheckstyle(r))
	}
	fmt.Fprintln(p.writer, `  </file>`)
	return nil
}

func (p *Checkstyle) Start() {
	fmt.Fprintln(p.writer, `<?xml version="1.0" encoding="UTF-8"?>
<checkstyle version="5.0">`)
}

func (p *Checkstyle) End() {
	fmt.Fprintln(p.writer, `</checkstyle>`)
}

func formatResultForCheckstyle(r result.Result) string {
	return fmt.Sprintf(`    <error column="%d" line="%d" message="%s" severity="%s" source="woke"/>`,
		r.GetStartPosition().Column,
		r.GetStartPosition().Line,
		r.Reason(),
		translateSeverityForCheckstyle(r.GetSeverity()),
	)
}

func translateSeverityForCheckstyle(s rule.Severity) string {
	if s == rule.SevError {
		return "error"
	}
	// treat everything else as a warning
	return "warning"
}
