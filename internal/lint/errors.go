package lint

import "fmt"

type ErrorCode string

type ErrorMessage string

type Error struct {
	ErrorCode       ErrorCode
	ErrorMessage    ErrorMessage
	Severity        Severity
	SourceKind      string
	SourceNamespace string
	SourceName      string
}

func (e *Error) Format() string {
	return fmt.Sprintf("%s | %s.%s.%s | %s | %s", e.Severity.String(), e.SourceNamespace, e.SourceKind, e.SourceName, string(e.ErrorCode), string(e.ErrorMessage))
}

type Severity int

func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarn:
		return "WARN"
	case SeverityError:
		return "ERROR"
	}

	panic("unrecognized severity")
	return "????"
}

const (
	SeverityInfo Severity = iota
	SeverityWarn
	SeverityError
)
