package util

import (
	"fmt"
)

// TODO: Change to Severity
type WarningType int

const (
	MINOR WarningType = iota
	SEVERE
)

func (wt WarningType) String() string {
	switch wt {
	case MINOR:
		return "Minor"
	case SEVERE:
		return "Severe"
	default:
		return "Unknown"
	}
}

type Warning interface {
	GetWarningType() WarningType
	Error() string
}

type WarningImpl struct {
	WarningType WarningType
	Message     string
}

func (w WarningImpl) Error() string {
	return fmt.Sprintf("warning type: %s. error is %s", w.WarningType.String(), w.Message)
}

func (w WarningImpl) GetWarningType() WarningType {
	return w.WarningType
}

func AppendMinorWarnings(warning1 Warning, warning2 Warning) Warning {
	if warning1 != nil && warning2 != nil {
		return WarningImpl{WarningType: MINOR, Message: fmt.Sprintf("%s, %s", warning1.Error(), warning2.Error())}
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
