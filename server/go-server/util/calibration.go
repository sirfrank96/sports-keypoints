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

func CheckIfDatapointExists(datapoint *skp.Datapoint) bool {
	if datapoint == nil {
		return false
	}
	return datapoint.X != 0 || datapoint.Y != 0
}

func VerifyDatapoint(datapoint *skp.Datapoint, datapointName string, threshold float64) Warning {
	if !CheckIfDatapointExists(datapoint) {
		return WarningImpl{
			Severity: SEVERE,
			Message:  fmt.Sprintf("could not find datapoint %s", datapointName),
		}
	}
	if datapoint.Confidence < threshold {
		return WarningImpl{
			Severity: MINOR,
			Message:  fmt.Sprintf("uncertain where %s is, confidence is %f. please make sure %s is visible in image", datapointName, datapoint.Confidence, datapointName),
		}
	}
	return nil
}
