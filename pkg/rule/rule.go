package rule

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/get-woke/woke/pkg/util"
)

var ignoreRuleRegex = regexp.MustCompile(`wokeignore:rule=(\S+)`)

// Rule is a linter rule
type Rule struct {
	Name         string   `yaml:"name"`
	Terms        []string `yaml:"terms"`
	Alternatives []string `yaml:"alternatives"`
	Note         string   `yaml:"note"`
	Severity     Severity `yaml:"severity"`

	re *regexp.Regexp
}

func (r *Rule) FindAllStringIndex(text string) [][]int {
	if r.re == nil {
		r.re = regexp.MustCompile(fmt.Sprintf(`(?i)\b(%s)\b`, strings.Join(r.Terms, "|")))
	}
	return r.re.FindAllStringIndex(text, -1)
}

func (r *Rule) String() string {
	return r.Name
}

// Reason returns a human-readable reason for the rule violation
func (r *Rule) Reason(violation string) string {
	reason := fmt.Sprintf("`%s` may be insensitive, ", violation)
	if len(r.Alternatives) > 0 {
		alt := make([]string, len(r.Alternatives))
		for i, a := range r.Alternatives {
			alt[i] = util.MarkdownCodify(a)
		}
		reason += fmt.Sprintf("use %s instead", strings.Join(alt, ", "))
	} else {
		reason += "try not to use it"
	}

	return reason
}

// ReasonWithNote returns a human-readable reason for the rule violation
// with an additional note, if defined.
func (r *Rule) ReasonWithNote(violation string) string {
	if len(r.Note) == 0 {
		return r.Reason(violation)
	}
	return fmt.Sprintf("%s (%s)", r.Reason(violation), r.Note)
}

// CanIgnoreLine returns a boolean value if the line contains the ignore directive.
// For example, if a line has anywhere, `woke:disable=whitelist`
// (should be commented out via whatever the language comment syntax is)
// it will not report that line in violation with the Rule with the name `whitelist`
func (r *Rule) CanIgnoreLine(line string) bool {
	matches := ignoreRuleRegex.FindAllStringSubmatch(line, -1)
	if matches == nil {
		return false
	}

	for _, match := range matches {
		if len(match) < 1 {
			continue
		}

		for _, m := range strings.Split(match[1], ",") {
			if m == r.Name {
				return true
			}
		}
	}

	return false
}
