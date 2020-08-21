package rule

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
