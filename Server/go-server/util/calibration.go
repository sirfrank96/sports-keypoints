package util

import (
	"fmt"
	"math"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

// TODO: Do i need feetline in this struct?
type CalibrationInfo struct {
	CalibrationType                  skp.CalibrationType `bson:"calibration_type,omitempty"`
	FeetLineMethod                   skp.FeetLineMethod  `bson:"feet_line_method,omitempty"`
	AxesCalibrationWarning           WarningImpl         `bson:"axes_calibration_warning,omitempty"`
	VanishingPointCalibrationWarning WarningImpl         `bson:"vanishing_point_calibration_warning,omitempty"`
	GolfBallWarning                  WarningImpl         `bson:"golf_ball_warning,omitempty"`
	GolfClubWarning                  WarningImpl         `bson:"golf_club_warning,omitempty"`
	FeetLine                         FeetLine            `bson:"feet_line,omitempty"`
	HorAxisLine                      Line                `bson:"hor_axis_line,omitempty"`
	VertAxisLine                     Line                `bson:"vert_axis_line,omitempty"`
	VanishingPoint                   Point               `bson:"vanishing_point,omitempty"`
	GolfBallPoint                    Point               `bson:"golf_ball_point,omitempty"`
	ClubButtPoint                    Point               `bson:"club_butt_point,omitempty"`
	ClubHeadPoint                    Point               `bson:"club_head_point,omitempty"`
}

func GetEmptyCalibrationInfo() *CalibrationInfo {
	return &CalibrationInfo{
		CalibrationType: skp.CalibrationType_NO_CALIBRATION,
	}
}

func CheckIfKeypointExists(keypoint *skp.Keypoint) bool {
	return keypoint.X != 0 || keypoint.Y != 0
}

func VerifyKeypoint(keypoint *skp.Keypoint, keypointName string, threshold float64) Warning {
	if !CheckIfKeypointExists(keypoint) {
		return WarningImpl{
			Severity: SEVERE,
			Message:  fmt.Sprintf("could not find keypoint %s", keypointName),
		}
	}
	if keypoint.Confidence < threshold {
		return WarningImpl{
			Severity: MINOR,
			Message:  fmt.Sprintf("uncertain where %s is, confidence is %f. please make sure %s is visible in image", keypointName, keypoint.Confidence, keypointName),
		}
	}
	return nil
}

// TODO: Configure confidence level, configure how far off 90 degrees axes can be
func VerifyCalibrationImageAxes(keypoints *skp.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) (*CalibrationInfo, Warning) {
	// Get horizontal axis
	feetLine, warning := GetFeetLine(keypoints, calibrationInfo.FeetLineMethod)
	if warning != nil {
		return nil, warning
	}
	horAxisLine := feetLine.Line
	// Get vertical axis
	if warning := VerifyKeypoint(keypoints.Midhip, "midhip", 0.5); warning != nil {
		return nil, warning
	}
	if warning := VerifyKeypoint(keypoints.Neck, "neck", 0.5); warning != nil {
		return nil, warning
	}
	vertAxisLine := GetLine(ConvertCvKeypointToPoint(keypoints.Midhip), ConvertCvKeypointToPoint(keypoints.Neck))
	// Check if angle between axes is around 90 degrees
	horDeg := ConvertSlopeToDegrees(horAxisLine.Slope)
	vertDeg := ConvertSlopeToDegrees(vertAxisLine.Slope)
	diff := math.Abs(vertDeg) + math.Abs(horDeg) - 90
	if math.Abs(diff) > 10 { // make this 5 or less after better test images
		return nil, WarningImpl{
			Severity: SEVERE,
			Message:  fmt.Sprintf("axes calibration image off. horizontal axis between heels is %f degrees. vertical axis between midhip and neck is %f degrees. difference of %f degrees is too large. please adjust camera, stance, or posture. recommend using alignment sticks to help calibration", horDeg, vertDeg, diff),
		}
	}
	fmt.Printf("Good axes calibration. Horizontal axis between heels is %f degrees. vertical axis between midhip and neck is %f degrees\n", horDeg, vertDeg)
	calibrationInfo.FeetLine = *feetLine
	calibrationInfo.HorAxisLine = horAxisLine
	calibrationInfo.VertAxisLine = *vertAxisLine
	return calibrationInfo, nil
}

func VerifyCalibrationImageVanishingPoint(keypoints *skp.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) (*CalibrationInfo, Warning) {
	// verify vanishing point image
	feetLine, warning := GetFeetLine(keypoints, calibrationInfo.FeetLineMethod)
	if warning != nil {
		return nil, warning
	}
	slopeDiff := math.Abs(feetLine.Line.Slope - calibrationInfo.VertAxisLine.Slope)
	if slopeDiff < float64(1) { // TODO: Configure how close slope is (and how to determine how close slope is)
		return nil, WarningImpl{
			Severity: SEVERE,
			Message:  fmt.Sprintf("vanishing point calibration image off. feet line slope %f and vertaxis line slope %f are too close (%f). make sure feet line is off centered or make sure alignment stick is pointed at target (parallel lines converge in distance)", feetLine.Line.Slope, calibrationInfo.VertAxisLine.Slope, slopeDiff),
		}
	}
	fmt.Printf("Good vanishing point calibration. feet line is %+v, and vertaxis line is %+v", feetLine.Line, calibrationInfo.VertAxisLine)
	intersection := GetIntersection(&feetLine.Line, &calibrationInfo.VertAxisLine)
	calibrationInfo.VanishingPoint = intersection.IntersectPoint
	return calibrationInfo, nil
}
