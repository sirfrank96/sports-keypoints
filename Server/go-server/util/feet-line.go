package util

import (
	"fmt"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

// TODO: Make FeetLineInfo internal struct
type FeetLineInfo struct {
	FeetLineMethod skp.FeetLineMethod `bson:"feet_line_method,omitempty"`
	LKeypoint      skp.Keypoint       `bson:"l_keypoint,omitempty"`
	RKeypoint      skp.Keypoint       `bson:"r_keypoint,omitempty"`
	LKeypointName  string             `bson:"l_keypoint_name,omitempty"`
	RKeypointName  string             `bson:"r_keypoint_name,omitempty"`
	Threshold      float64            `bson:"threshold,omitempty"`
}

type FeetLine struct {
	FeetLineMethod skp.FeetLineMethod `bson:"feet_line_method,omitempty"`
	LPoint         Point              `bson:"l_point,omitempty"`
	RPoint         Point              `bson:"r_point,omitempty"`
	Line           Line               `bson:"line,omitempty"`
}

func GetFeetLine(keypoints *skp.Body25PoseKeypoints, feetLineMethod skp.FeetLineMethod) (*FeetLine, Warning) {
	feetLineInfo := GetFeetLineInfo(keypoints, feetLineMethod)
	warning := VerifyFeetLineInfo(feetLineInfo)
	if warning != nil && warning.GetSeverity() == SEVERE {
		return nil, warning
	}
	feetLine := GetFeetLineFromInfo(feetLineInfo)
	return feetLine, warning
}

//TODO: CONFIGURE THRESHOLD
func GetFeetLineInfo(keypoints *skp.Body25PoseKeypoints, feetLineMethod skp.FeetLineMethod) *FeetLineInfo {
	feetLineInfo := &FeetLineInfo{FeetLineMethod: feetLineMethod, Threshold: 0.5}
	lKeypoint, lKeypointName := GetLeftFootPoint(keypoints, feetLineMethod)
	feetLineInfo.LKeypoint = *lKeypoint
	feetLineInfo.LKeypointName = lKeypointName
	rKeypoint, rKeypointName := GetRightFootPoint(keypoints, feetLineMethod)
	feetLineInfo.RKeypoint = *rKeypoint
	feetLineInfo.RKeypointName = rKeypointName
	return feetLineInfo
}

func VerifyFeetLineInfo(feetLineInfo *FeetLineInfo) Warning {
	var warning Warning
	if w := VerifyKeypoint(&feetLineInfo.LKeypoint, feetLineInfo.LKeypointName, feetLineInfo.Threshold); w != nil {
		if w.GetSeverity() == SEVERE {
			return w
		}
		wStruct := WarningImpl{
			Severity: w.GetSeverity(),
			Message:  fmt.Sprintf("%s, please set a different FeetLineMethod", w.Error()),
		}
		warning = AppendMinorWarnings(warning, wStruct)
	}
	if w := VerifyKeypoint(&feetLineInfo.RKeypoint, feetLineInfo.RKeypointName, feetLineInfo.Threshold); w != nil {
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
	feetLine.LPoint = *ConvertCvKeypointToPoint(&feetLineInfo.LKeypoint)
	feetLine.RPoint = *ConvertCvKeypointToPoint(&feetLineInfo.RKeypoint)
	feetLine.Line = *GetLine(&feetLine.LPoint, &feetLine.RPoint)
	return feetLine
}

func GetLeftFootPoint(keypoints *skp.Body25PoseKeypoints, feetLineMethod skp.FeetLineMethod) (*skp.Keypoint, string) {
	if feetLineMethod == skp.FeetLineMethod_USE_TOE_LINE {
		return keypoints.LBigToe, "left big toe"
	} else { // default is USE_HEEL_LINE
		return keypoints.LHeel, "left heel"
	}
}

func GetRightFootPoint(keypoints *skp.Body25PoseKeypoints, feetLineMethod skp.FeetLineMethod) (*skp.Keypoint, string) {
	if feetLineMethod == skp.FeetLineMethod_USE_TOE_LINE {
		return keypoints.RBigToe, "right big toe"
	} else { // default is USE_HEEL_LINE
		return keypoints.RHeel, "right heel"
	}
}
