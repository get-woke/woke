package rule

import (
	"fmt"
	"go/token"
	"regexp"

	"gopkg.in/yaml.v2"
)

// Rule is a linter rule
type Rule struct {
	Name         string         // `yaml:"name"`
	Regexp       *regexp.Regexp // `yaml:"regexp"`
	Alternatives string         // `yaml:"alternatives"`
	Severity     Severity
}

func (r *Rule) String() string {
	return r.Name
}

// compile-time check that Rule satisfies the yaml Unmarshaler
var _ yaml.Unmarshaler = (*Rule)(nil)

// UnmarshalYAML to enforce regexp at the unmarshal level
func (r *Rule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	a := make(map[string]string)
	if err := unmarshal(a); err != nil {
		return err
	}
	r.Alternatives = a["alternatives"]
	r.Name = a["name"]
	r.Severity = NewSeverity(a["severity"])
	if re, ok := a["regexp"]; !ok {
		r.Regexp = regexp.MustCompile(fmt.Sprintf(`(?i)\b(%s)\b`, r.Name))
	} else {
		r.Regexp = regexp.MustCompile(fmt.Sprintf(`(?i)%s`, re))
	}
	return nil
}

// FindResults returns the results that match the rule for the given text.
// filename and line are only used for the Position
func (r *Rule) FindResults(text, filename string, line int) (rs []Result) {
	idxs := r.Regexp.FindAllStringIndex(text, -1)

	for _, idx := range idxs {
		newResult := Result{
			Rule:  r,
			Match: text[idx[0]:idx[1]],
			Position: &token.Position{
				Filename: filename,
				Line:     line,
				Column:   idx[0],
			},
		}

		rs = append(rs, newResult)
	}
	return
}
