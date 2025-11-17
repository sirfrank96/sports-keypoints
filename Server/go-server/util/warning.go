package util

import (
	"fmt"
)

// TODO: Change to Severity
type Severity int

const (
	MINOR Severity = iota
	SEVERE
)

func (s Severity) String() string {
	switch s {
	case MINOR:
		return "Minor"
	case SEVERE:
		return "Severe"
	default:
		return "Unknown"
	}
}

type Warning interface {
	GetSeverity() Severity
	Error() string
}

type WarningImpl struct {
	Severity Severity
	Message  string
}

func (w WarningImpl) Error() string {
	return fmt.Sprintf("severity: %s. error is %s", w.Severity.String(), w.Message)
}

func (w WarningImpl) GetSeverity() Severity {
	return w.Severity
}

func AppendMinorWarnings(warning1 Warning, warning2 Warning) Warning {
	if warning1 != nil && warning2 != nil {
		return WarningImpl{Severity: MINOR, Message: fmt.Sprintf("%s, %s", warning1.Error(), warning2.Error())}
	} else if warning1 == nil {
		return warning2
	} else if warning2 == nil {
		return warning1
	} else {
		return nil
	}
}

/*
func convertErrorToWarning(err error, warningType WarningType) warning {
	return Warning{warningType: warningType, message: err.Error()}
}

func appendError(err1 error, err2 error) error {
	if err1 != nil && err2 != nil {
		return fmt.Errorf("%w, %w", err1, err2)
	} else if err1 == nil {
		return err2
	} else if err2 == nil {
		return err1
	} else {
		return nil
	}
}*/
