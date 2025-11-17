package util

import (
	"fmt"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
)

// TODO: Make FeetLineInfo internal struct
type FeetLineInfo struct {
	FeetLineMethod cv.FeetLineMethod `bson:"feet_line_method,omitempty"`
	LKeypoint      cv.Keypoint       `bson:"l_keypoint,omitempty"`
	RKeypoint      cv.Keypoint       `bson:"r_keypoint,omitempty"`
	LKeypointName  string            `bson:"l_keypoint_name,omitempty"`
	RKeypointName  string            `bson:"r_keypoint_name,omitempty"`
	Threshold      float64           `bson:"threshold,omitempty"`
}

type FeetLine struct {
	FeetLineMethod cv.FeetLineMethod `bson:"feet_line_method,omitempty"`
	LPoint         Point             `bson:"l_point,omitempty"`
	RPoint         Point             `bson:"r_point,omitempty"`
	Line           Line              `bson:"line,omitempty"`
}

func GetFeetLine(keypoints *cv.Body25PoseKeypoints, feetLineMethod cv.FeetLineMethod) (*FeetLine, Warning) {
	feetLineInfo := GetFeetLineInfo(keypoints, feetLineMethod)
	warning := VerifyFeetLineInfo(feetLineInfo)
	if warning != nil && warning.GetWarningType() == SEVERE {
		return nil, warning
	}
	feetLine := GetFeetLineFromInfo(feetLineInfo)
	return feetLine, warning
}

//TODO: CONFIGURE THRESHOLD
func GetFeetLineInfo(keypoints *cv.Body25PoseKeypoints, feetLineMethod cv.FeetLineMethod) *FeetLineInfo {
	feetLineInfo := &FeetLineInfo{FeetLineMethod: feetLineMethod, Threshold: 0.6}
	if feetLineMethod == cv.FeetLineMethod_USE_TOE_LINE {
		feetLineInfo.LKeypoint = *keypoints.LBigToe
		feetLineInfo.LKeypointName = "left big toe"
		feetLineInfo.RKeypoint = *keypoints.RBigToe
		feetLineInfo.RKeypointName = "right big toe"
	} else { // default is USE_HEEL_LINE
		feetLineInfo.LKeypoint = *keypoints.LHeel
		feetLineInfo.LKeypointName = "left heel"
		feetLineInfo.RKeypoint = *keypoints.RHeel
		feetLineInfo.LKeypointName = "right heel"
	}
	return feetLineInfo
}

func VerifyFeetLineInfo(feetLineInfo *FeetLineInfo) Warning {
	var warning Warning
	if w := VerifyKeypoint(&feetLineInfo.LKeypoint, feetLineInfo.LKeypointName, feetLineInfo.Threshold); w != nil {
		if w.GetWarningType() == SEVERE {
			return w
		}
		wStruct := WarningImpl{
			WarningType: w.GetWarningType(),
			Message:     fmt.Sprintf("%w, please set a different FeetLineMethod", w.Error()),
		}
		warning = AppendMinorWarnings(warning, wStruct)
	}
	if w := VerifyKeypoint(&feetLineInfo.RKeypoint, feetLineInfo.RKeypointName, feetLineInfo.Threshold); w != nil {
		if w.GetWarningType() == SEVERE {
			return w
		}
		wStruct := WarningImpl{
			WarningType: w.GetWarningType(),
			Message:     fmt.Sprintf("%w, please set a different FeetLineMethod", w.Error()),
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
