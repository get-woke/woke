package printer

import (
	"fmt"

	"github.com/get-woke/woke/pkg/result"
)

// ErrorFormat is a ErrorFormat printer meant to be read by errorformat
// https://vim-jp.org/vimdoc-en/quickfix.html#error-file-format
type ErrorFormat struct {
	format string
}

// NewErrorFormat returns a new ErrorFormat printer
func NewErrorFormat() *ErrorFormat {
	return &ErrorFormat{
		format: "%s:%d:%d [%c] %s",
	}
}

// Print prints in the format 'filename:line:column: [s] message'
// where `s` is the first character of Severity
// The formatting is meant to line up with errorformat
// "%f:%l:%c [%t] %m"
// https://vim-jp.org/vimdoc-en/quickfix.html#error-file-format
func (e *ErrorFormat) Print(fs *result.FileResults) error {
	for _, r := range fs.Results {
		fmt.Printf(e.format+"\n",
			r.StartPosition.Filename,
			r.StartPosition.Line,
			r.StartPosition.Column,
			r.Rule.Severity.String()[0],
			r.Reason())
	}
	return nil
}
