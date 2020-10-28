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
	Options      Options  `yaml:"options"`

	reWordBoundary *regexp.Regexp
	re             *regexp.Regexp
}

func (r *Rule) findAllStringSubmatchIndex(text string) [][]int {
	if r.Options.WordBoundary {
		return r.reWordBoundary.FindAllStringSubmatchIndex(text, -1)
	}
	return r.re.FindAllStringSubmatchIndex(text, -1)
}

// FindMatchIndexes returns the start and end indexes for all rule violations for the text supplied.
func (r *Rule) FindMatchIndexes(text string) [][]int {
	if r.Disabled() {
		return [][]int(nil)
	}

	r.SetRegexp()
	matches := r.findAllStringSubmatchIndex(text)
	if matches == nil {
		return [][]int(nil)
	}

	idx := [][]int{}

	// Need to return a list of int pairs, which are the start and end index
	// of all matches in all capture groups. For FindAllStringSubmatchIndex,
	// Submatch 0 is the match of the entire expression, submatch 1 the match
	// of the first parenthesized subexpression, and so on. We only care about Submatch 1+
	for _, m := range matches {
		if len(m) < 4 {
			continue
		}

		// Right now, assume there's only one capture group.
		// This should be updated to support more capture groups if necessary.
		start := m[2]
		end := m[3]

		if start == -1 || end == -1 {
			// something went wrong with the regex
			continue
		}

		idx = append(idx, []int{start, end})
	}

	return idx
}

// MatchString reports whether the string s
// contains any match of the regular expression re.
func (r *Rule) MatchString(s string, wordBoundary bool) bool {
	if r.Disabled() {
		return false
	}

	r.SetRegexp()

	if wordBoundary {
		return r.reWordBoundary.MatchString(s)
	}

	return r.re.MatchString(s)
}

// SetRegexp populates the regex for matching this rule
func (r *Rule) SetRegexp() {
	if r.re != nil && r.reWordBoundary != nil {
		return
	}
	group := strings.Join(escape(r.Terms), "|")
	r.reWordBoundary = regexp.MustCompile(fmt.Sprintf(`(?i)\b(%s)\b`, group))
	r.re = regexp.MustCompile(fmt.Sprintf(`(?i)(%s)`, group))
}

// Reason returns a human-readable reason for the rule violation
func (r *Rule) Reason(violation string) string {
	// fall back to the rule name if no violation was found
	// violation is mostly used for informational purposes
	if len(violation) == 0 {
		violation = r.Name
	}

	reason := new(strings.Builder)
	reason.WriteString(util.MarkdownCodify(violation) + " may be insensitive, ")

	if len(r.Alternatives) > 0 {
		alt := make([]string, len(r.Alternatives))
		for i, a := range r.Alternatives {
			alt[i] = util.MarkdownCodify(a)
		}
		reason.WriteString(fmt.Sprintf("use %s instead", strings.Join(alt, ", ")))
	} else {
		reason.WriteString("try not to use it")
	}

	return reason.String()
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
// For example, if a line has anywhere, wokeignore:rule=whitelist
// (should be commented out via whatever the language comment syntax is)
// it will not report that line in violation with the Rule with the name `whitelist` wokeignore:rule=whitelist
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

func escape(ss []string) []string {
	for i, s := range ss {
		ss[i] = regexp.QuoteMeta(s)
	}
	return ss
}

// Disabled denotes if the rule is disabled
// If no terms are provided, this essentially disables the rule
// which is helpful for disabling default rules. Eventually, there should be a better
// way to disable a default rule, and then, if a rule has no Terms, it falls back to the Name.
func (r *Rule) Disabled() bool {
	return len(r.Terms) == 0
}
