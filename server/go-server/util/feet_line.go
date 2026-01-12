package util

import (
	"fmt"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

type FeetLineInfo struct {
	FeetLineMethod skp.FeetLineMethod `bson:"feet_line_method,omitempty"`
	LDatapoint     skp.Datapoint      `bson:"l_datapoint,omitempty"`
	RDatapoint     skp.Datapoint      `bson:"r_datapoint,omitempty"`
	LDatapointName string             `bson:"l_datapoint_name,omitempty"`
	RDatapointName string             `bson:"r_datapoint_name,omitempty"`
	Threshold      float64            `bson:"threshold,omitempty"`
}

type FeetLine struct {
	FeetLineMethod skp.FeetLineMethod `bson:"feet_line_method,omitempty"`
	LPoint         Point              `bson:"l_point,omitempty"`
	RPoint         Point              `bson:"r_point,omitempty"`
	Line           Line               `bson:"line,omitempty"`
}

func GetFeetLine(datapoints *skp.Body25PoseDatapoints, feetLineMethod skp.FeetLineMethod) (*FeetLine, Warning) {
	feetLineInfo := GetFeetLineInfo(datapoints, feetLineMethod)
	warning := VerifyFeetLineInfo(feetLineInfo)
	if warning != nil && warning.GetSeverity() == SEVERE {
		return nil, warning
	}
	feetLine := GetFeetLineFromInfo(feetLineInfo)
	return feetLine, warning
}

// TODO: Configure threshold
func GetFeetLineInfo(datapoints *skp.Body25PoseDatapoints, feetLineMethod skp.FeetLineMethod) *FeetLineInfo {
	feetLineInfo := &FeetLineInfo{FeetLineMethod: feetLineMethod, Threshold: 0.5}
	lDatapoint, lDatapointName := GetLeftFootPoint(datapoints, feetLineMethod)
	feetLineInfo.LDatapoint = *lDatapoint
	feetLineInfo.LDatapointName = lDatapointName
	rDatapoint, rDatapointName := GetRightFootPoint(datapoints, feetLineMethod)
	feetLineInfo.RDatapoint = *rDatapoint
	feetLineInfo.RDatapointName = rDatapointName
	return feetLineInfo
}

func VerifyFeetLineInfo(feetLineInfo *FeetLineInfo) Warning {
	var warning Warning
	if w := VerifyDatapoint(&feetLineInfo.LDatapoint, feetLineInfo.LDatapointName, feetLineInfo.Threshold); w != nil {
		if w.GetSeverity() == SEVERE {
			return w
		}
		wStruct := WarningImpl{
			Severity: w.GetSeverity(),
			Message:  fmt.Sprintf("%s, please set a different FeetLineMethod", w.Error()),
		}
		warning = AppendMinorWarnings(warning, wStruct)
	}
	if w := VerifyDatapoint(&feetLineInfo.RDatapoint, feetLineInfo.RDatapointName, feetLineInfo.Threshold); w != nil {
		if w.GetSeverity() == SEVERE {
			return w
		}
		wStruct := WarningImpl{
			Severity: w.GetSeverity(),
			Message:  fmt.Sprintf("%s, please set a different FeetLineMethod", w.Error()),
		}
		warning = AppendMinorWarnings(warning, wStruct)
	}
	return warning
}

func GetFeetLineFromInfo(feetLineInfo *FeetLineInfo) *FeetLine {
	feetLine := &FeetLine{FeetLineMethod: feetLineInfo.FeetLineMethod}
	feetLine.LPoint = *ConvertDatapointToPoint(&feetLineInfo.LDatapoint)
	feetLine.RPoint = *ConvertDatapointToPoint(&feetLineInfo.RDatapoint)
	feetLine.Line = *GetLine(&feetLine.RPoint, &feetLine.LPoint)
	return feetLine
}

func GetLeftFootPoint(datapoints *skp.Body25PoseDatapoints, feetLineMethod skp.FeetLineMethod) (*skp.Datapoint, string) {
	if feetLineMethod == skp.FeetLineMethod_USE_TOE_LINE {
		return datapoints.LBigToe, "left big toe"
	} else { // default is USE_HEEL_LINE
		return datapoints.LHeel, "left heel"
	}
}

func GetRightFootPoint(datapoints *skp.Body25PoseDatapoints, feetLineMethod skp.FeetLineMethod) (*skp.Datapoint, string) {
	if feetLineMethod == skp.FeetLineMethod_USE_TOE_LINE {
		return datapoints.RBigToe, "right big toe"
	} else { // default is USE_HEEL_LINE
		return datapoints.RHeel, "right heel"
	}
}
