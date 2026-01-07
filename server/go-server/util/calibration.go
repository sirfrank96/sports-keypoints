package util

import (
	"fmt"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

type CalibrationInfo struct {
	CalibrationType              skp.CalibrationType `bson:"calibration_type,omitempty"`
	FeetLineMethod               skp.FeetLineMethod  `bson:"feet_line_method,omitempty"`
	HorAxisLine                  Line                `bson:"hor_axis_line,omitempty"`
	VertAxisLine                 Line                `bson:"vert_axis_line,omitempty"`
	VanishingPoint               Point               `bson:"vanishing_point,omitempty"`
	GolfBallPoint                skp.Keypoint        `bson:"golf_ball_point,omitempty"`
	ClubButtPoint                skp.Keypoint        `bson:"club_butt_point,omitempty"`
	ClubHeadPoint                skp.Keypoint        `bson:"club_head_point,omitempty"`
	ShoulderTilt                 skp.Double          `bson:"shoulder_tilt,omitempty"`
	ManualGenerated              bool                `bson:"manual_generated,omitempty"`
	CalibrationImgAxes           []byte              `bson:"calibration_img_axes,omitempty"`
	CalibrationImgVanishingPoint []byte              `bson:"calibration_img_vanishing_point,omitempty"`
}

func GetEmptyCalibrationInfo() *CalibrationInfo {
	return &CalibrationInfo{
		CalibrationType: skp.CalibrationType_NO_CALIBRATION,
	}
}

func VerifyDouble(double *skp.Double) Warning {
	if double == nil {
		return WarningImpl{
			Severity: SEVERE,
			Message:  "value does not exist",
		}
	}
	if double.Warning != "" {
		return WarningImpl{
			Severity: SEVERE,
			Message:  fmt.Sprintf("warning for value: %s", double.Warning),
		}
	}
	return nil
}

func CheckIfKeypointExists(keypoint *skp.Keypoint) bool {
	if keypoint == nil {
		return false
	}
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
