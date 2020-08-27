package rule

import "github.com/fatih/color"

// Severity is a log severity
type Severity int

const (
	// SevInfo translates to Info
	SevInfo Severity = iota
	// SevWarn translates to Warn
	SevWarn
	// SevError translates to Error
	SevError
)

// NewSeverity turns a string into a Severity
func NewSeverity(s string) Severity {
	switch s {
	case SevInfo.String():
		return SevInfo
	case SevWarn.String():
		return SevWarn
	case SevError.String():
		return SevError
	}
	return SevWarn
}

func (s Severity) String() string {
	return [...]string{"info", "warn", "error"}[s]
}

func (s *Severity) Colorize() string {
	switch *s {
	case SevInfo:
		return color.GreenString(s.String())
	case SevWarn:
		return color.YellowString(s.String())
	case SevError:
		return color.RedString(s.String())
	}
	return ""
}
