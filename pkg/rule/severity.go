package rule

type Severity int

const (
	SevInfo Severity = iota
	SevWarn
	SevError
)

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
