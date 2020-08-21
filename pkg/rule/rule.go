package rule

import (
	"fmt"
	"regexp"

	"gopkg.in/yaml.v2"
)

// Rule is a linter rule
type Rule struct {
	Word         string         // `yaml:"word"`
	Regexp       *regexp.Regexp // `yaml:"regexp"`
	Alternatives string         // `yaml:"alternatives"`
}

func (r *Rule) String() string {
	return r.Word
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
	r.Word = a["word"]
	if re, ok := a["regexp"]; !ok {
		r.Regexp = regexp.MustCompile(fmt.Sprintf(`(?i)\b(%s)\b`, r.Word))
	} else {
		r.Regexp = regexp.MustCompile(fmt.Sprintf(`(?i)%s`, re))
	}
	return nil
}
