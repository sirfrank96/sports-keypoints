package processor

import (
	"fmt"
	"math"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
)

type CalibrationInfo struct {
	feetLine       *FeetLine
	horAxisLine    *Line
	vertAxisLine   *Line
	vanishingPoint *Point
}

func checkIfKeypointExists(keypoint *cv.Keypoint) bool {
	return keypoint.X != 0 || keypoint.Y != 0
}

func verifyKeypoint(keypoint *cv.Keypoint, keypointName string, threshold float64) warning {
	if !checkIfKeypointExists(keypoint) {
		return Warning{
			warningType: SEVERE,
			message:     fmt.Sprintf("could not find keypoint %s", keypointName),
		}
	}
	if keypoint.Confidence < threshold {
		return Warning{
			warningType: MINOR,
			message:     fmt.Sprintf("uncertain where %s is, confidence is %f. please make sure %s is visible in image", keypointName, keypoint.Confidence, keypointName),
		}
	}
	return nil
}

// TODO: Configure confidence level
func VerifyCalibrationImageAxes(keypoints *cv.Body25PoseKeypoints, feetLineMethod cv.FeetLineMethod) (*CalibrationInfo, warning) {
	// Get horizontal axis
	feetLine, warning := GetFeetLine(keypoints, feetLineMethod)
	if warning != nil {
		return nil, warning
	}
	horAxisLine := feetLine.line
	// Get vertical axis
	if warning := verifyKeypoint(keypoints.Midhip, "midhip", 0.5); warning != nil {
		return nil, warning
	}
	if warning := verifyKeypoint(keypoints.Neck, "neck", 0.5); warning != nil {
		return nil, warning
	}
	vertAxisLine := getLine(convertCvKeypointToPoint(keypoints.Midhip), convertCvKeypointToPoint(keypoints.Neck))
	// Check if angle between axes is around 90 degrees
	horDeg := convertSlopeToDegrees(horAxisLine.slope)
	vertDeg := convertSlopeToDegrees(vertAxisLine.slope)
	diff := math.Abs(vertDeg) + math.Abs(horDeg) - 90
	if math.Abs(diff) > 10 { // make this 5 or less after better test images
		return nil, Warning{
			warningType: SEVERE,
			message:     fmt.Sprintf("axes calibration image off. horizontal axis between heels is %f degrees. vertical axis between midhip and neck is %f degrees. difference of %f degrees is too large. please adjust camera, stance, or posture. recommend using alignment sticks to help calibration", horDeg, vertDeg, diff),
		}
	}
	fmt.Printf("Good axes calibration. Horizontal axis between heels is %f degrees. vertical axis between midhip and neck is %f degrees\n", horDeg, vertDeg)
	return &CalibrationInfo{feetLine: feetLine, horAxisLine: horAxisLine, vertAxisLine: vertAxisLine}, nil
}
