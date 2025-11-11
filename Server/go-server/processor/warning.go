package processor

import (
	"fmt"
)

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

type warning interface {
	WarningType() WarningType
	Error() string
}

type Warning struct {
	warningType WarningType
	message     string
}

func (w Warning) Error() string {
	return fmt.Sprintf("warning type: %s. error is %s", w.warningType.String(), w.message)
}

func (w Warning) WarningType() WarningType {
	return w.warningType
}

func appendMinorWarnings(warning1 warning, warning2 warning) warning {
	if warning1 != nil && warning2 != nil {
		return Warning{warningType: MINOR, message: fmt.Sprintf("%s, %s", warning1.Error(), warning2.Error())}
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
