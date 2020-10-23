package result

import (
	"encoding/json"
	"fmt"
	"go/token"

	"github.com/get-woke/woke/pkg/rule"

	"github.com/rs/zerolog/log"
)

// MaxLineLength is the max line length that this printer
// will show the source of the violation and the location within the line of the violation.
// Helps avoid consuming the console when minified files contain violations.
const MaxLineLength = 200

// LineResult contains data about the result of a broken rule
type LineResult struct {
	Rule      *rule.Rule
	Violation string
	// Line is the full string of the line, unless it's over MaxLintLength,
	// where Line will be an empty string
	Line          string
	StartPosition *token.Position
	EndPosition   *token.Position
}

// NewLineResult returns a LineResult based on the metadata from a violation
func NewLineResult(r *rule.Rule, violation, filename string, line, startColumn, endColumn int) LineResult {
	return LineResult{
		Rule:      r,
		Violation: violation,
		StartPosition: &token.Position{
			Filename: filename,
			Line:     line,
			Column:   startColumn,
		},
		EndPosition: &token.Position{
			Filename: filename,
			Line:     line,
			Column:   endColumn,
		},
	}
}

// FindResults returns the results that match the rule for the given text.
// filename and line are only used for the Position
func FindResults(r *rule.Rule, filename, text string, line int) (rs []Result) {
	if r.CanIgnoreLine(text) {
		log.Debug().
			Str("rule", r.Name).
			Str("file", filename).
			Int("line", line).
			Msg("ignoring via in-line")
		return
	}

	idxs := r.FindMatchIndexes(text)

	for _, idx := range idxs {
		start := idx[0]
		end := idx[1]
		newResult := NewLineResult(r, text[start:end], filename, line, start, end)

		if len(text) < MaxLineLength {
			newResult.Line = text
		}

		rs = append(rs, newResult)
	}
	return
}

// Reason outputs the suggested alternatives for this rule
func (r LineResult) Reason() string {
	return r.Rule.Reason(r.Violation)
}

func (r LineResult) String() string {
	pos := fmt.Sprintf("%s-%s",
		r.StartPosition.String(),
		r.EndPosition.String())
	return fmt.Sprintf("    %-14s %-10s %s", pos, r.Rule.Severity, r.Reason())
}

// GetSeverity returns the rule severity for the Result
func (r LineResult) GetSeverity() rule.Severity { return r.Rule.Severity }

// GetStartPosition returns the start position for the Result
func (r LineResult) GetStartPosition() *token.Position { return r.StartPosition }

// GetEndPosition returns the start position for the Result
func (r LineResult) GetEndPosition() *token.Position { return r.EndPosition }

// GetLine returns the entire line for the LineResult
func (r LineResult) GetLine() string { return r.Line }

type jsonLineResult LineResult

// MarshalJSON override to include Reason in the json response
func (r LineResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		jsonLineResult
		Reason string
	}{
		jsonLineResult: jsonLineResult(r),
		Reason:         r.Reason(),
	})
}
