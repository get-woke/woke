package result

import (
	"strings"

	"github.com/get-woke/woke/pkg/rule"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Results contains a list of Result
type Results struct {
	Results []Result
}

// Add adds a single Result object onto the Results stack
func (rs *Results) Add(r *Result) {
	rs.Results = append(rs.Results, *r)
}

// Push pushes a list of Result objects onto the Results stack
func (rs *Results) Push(r ...Result) {
	for _, result := range r {
		rs.Add(&result)
	}
}

func (rs *Results) String() string {
	s := []string{}
	for _, r := range rs.Results {
		s = append(s, r.String())
	}
	return strings.Join(s, "\n")
}

// Output is the logger output of results
func (rs *Results) Output() {
	var logger *zerolog.Event
	for _, r := range rs.Results {
		switch r.Rule.Severity {
		case rule.SevError:
			logger = log.Error()
		case rule.SevInfo:
			logger = log.Info()
		case rule.SevWarn:
			logger = log.Warn()
		}
		logger.Msg(r.String())
	}
}

// func (rs *Results) Pretty() string {
// 	// TODO: Implement this
// 	return rs.String()
// }

// func (rs *Results) Simple() string {
// 	// TODO: Implement this
// 	return rs.String()
// }

// func (rs *Results) JSON() string {
// 	b, _ := json.Marshal(rs)
// 	return string(b)
// }

// func (rs *Results) OutputString(t string) string {
// 	switch NewOutputType(t) {
// 	case OutputTypePretty:
// 		return rs.Pretty()
// 	case OutputTypeSimple:
// 		return rs.String()
// 	case OutputTypeJSON:
// 		b, _ := json.Marshal(rs)
// 		return string(b)
// 	}

// 	return fmt.Sprintf("Unsupported output type: %s", t)
// }
