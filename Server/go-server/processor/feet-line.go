package processor

import (
	"fmt"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
)

type FeetLineInfo struct {
	feetLineMethod cv.FeetLineMethod
	lKeypoint      *cv.Keypoint
	rKeypoint      *cv.Keypoint
	lKeypointName  string
	rKeypointName  string
	threshold      float64
}

type FeetLine struct {
	feetLineMethod cv.FeetLineMethod
	lPoint         *Point
	rPoint         *Point
	line           *Line
}

func GetFeetLine(keypoints *cv.Body25PoseKeypoints, feetLineMethod cv.FeetLineMethod) (*FeetLine, warning) {
	feetLineInfo := getFeetLineInfo(keypoints, feetLineMethod)
	warning := verifyFeetLineInfo(feetLineInfo)
	if warning != nil && warning.WarningType() == SEVERE {
		return nil, warning
	}
	feetLine := getFeetLine(feetLineInfo)
	return feetLine, warning
}

//TODO: CONFIGURE THRESHOLD
func getFeetLineInfo(keypoints *cv.Body25PoseKeypoints, feetLineMethod cv.FeetLineMethod) *FeetLineInfo {
	feetLineInfo := &FeetLineInfo{feetLineMethod: feetLineMethod, threshold: 0.6}
	if feetLineMethod == cv.FeetLineMethod_USE_TOE_LINE {
		feetLineInfo.lKeypoint = keypoints.LBigToe
		feetLineInfo.lKeypointName = "left big toe"
		feetLineInfo.rKeypoint = keypoints.RBigToe
		feetLineInfo.rKeypointName = "right big toe"
	} else { // default is USE_HEEL_LINE
		feetLineInfo.lKeypoint = keypoints.LHeel
		feetLineInfo.lKeypointName = "left heel"
		feetLineInfo.rKeypoint = keypoints.RHeel
		feetLineInfo.lKeypointName = "right heel"
	}
	return feetLineInfo
}

func verifyFeetLineInfo(feetLineInfo *FeetLineInfo) warning {
	var warning warning
	if w := verifyKeypoint(feetLineInfo.lKeypoint, feetLineInfo.lKeypointName, feetLineInfo.threshold); w != nil {
		if w.WarningType() == SEVERE {
			return w
		}
		wStruct := Warning{
			warningType: w.WarningType(),
			message:     fmt.Sprintf("%w, please set a different FeetLineMethod", w.Error()),
		}
		warning = appendMinorWarnings(warning, wStruct)
	}
	if w := verifyKeypoint(feetLineInfo.rKeypoint, feetLineInfo.rKeypointName, feetLineInfo.threshold); w != nil {
		if w.WarningType() == SEVERE {
			return w
		}
		wStruct := Warning{
			warningType: w.WarningType(),
			message:     fmt.Sprintf("%w, please set a different FeetLineMethod", w.Error()),
		}
		warning = appendMinorWarnings(warning, wStruct)
	}
	return warning
}

func getFeetLine(feetLineInfo *FeetLineInfo) *FeetLine {
	feetLine := &FeetLine{feetLineMethod: feetLineInfo.feetLineMethod}
	feetLine.lPoint = convertCvKeypointToPoint(feetLineInfo.lKeypoint)
	feetLine.rPoint = convertCvKeypointToPoint(feetLineInfo.rKeypoint)
	feetLine.line = getLine(feetLine.lPoint, feetLine.rPoint)
	return feetLine
}
