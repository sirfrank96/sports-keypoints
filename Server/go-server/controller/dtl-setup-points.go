package controller

import (
	"context"
	"fmt"
	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

//assuming right handed golfer

//2 calibration images? 1 for perpendicular axes, 1 for vanishing points
//img 1: stand straddled, check to make sure heel horizontal and spine vertical are close to 90
//img 2: set alignment stick not centered, point alignment stick at target, set up with heels against alignment stick feet shoulder width or wider (check that heels are not centered in image)
// get vanishing point, intersection of vertaxis and heels axis
func VerifyDTLCalibrationImages(axesKeypoints *skp.Body25PoseKeypoints, vanishingPointKeypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (*util.CalibrationInfo, util.Warning) {
	// verify axes image
	calibrationInfo, warning := util.VerifyCalibrationImageAxes(axesKeypoints, calibrationInfo)
	if warning != nil {
		return nil, warning
	}
	// verify vanishing point image
	calibrationInfo, warning = util.VerifyCalibrationImageVanishingPoint(vanishingPointKeypoints, calibrationInfo)
	if warning != nil {
		return nil, warning
	}
	return calibrationInfo, nil
}

func CalculateDTLSetupPoints(ctx context.Context, keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) *skp.DTLGolfSetupPoints {
	spineAngle, warning := GetSpineAngle(keypoints, calibrationInfo)
	var spineAngleWarning string
	if warning != nil {
		spineAngleWarning = warning.Error()
	}
	fmt.Printf("Spine angle is %f", spineAngle)
	feetAlignment, warning := GetFeetAlignment(keypoints, calibrationInfo)
	var feetAlignmentWarning string
	if warning != nil {
		feetAlignmentWarning = warning.Error()
	}
	fmt.Printf("Feet alignment is %f", feetAlignment)
	// TODO: Add heel and toe alignment based on FeetLineMethod
	kneeBend, warning := GetKneeBend(keypoints)
	var kneeBendWarning string
	if warning != nil {
		kneeBendWarning = warning.Error()
	}
	fmt.Printf("Knee bend is %f", kneeBend)
	shoulderAlignment, warning := GetShoulderAlignment(keypoints, calibrationInfo)
	var shoulderAlignmentWarning string
	if warning != nil {
		shoulderAlignmentWarning = warning.Error()
	}
	fmt.Printf("Shoulder alignment is %f", shoulderAlignment)

	dtlGolfSetupPoints := &skp.DTLGolfSetupPoints{
		SpineAngle: &skp.Double{
			Data:    spineAngle,
			Warning: spineAngleWarning,
		},
		FeetAlignment: &skp.Double{
			Data:    feetAlignment,
			Warning: feetAlignmentWarning,
		},
		KneeBend: &skp.Double{
			Data:    kneeBend,
			Warning: kneeBendWarning,
		},
		ShoulderAlignment: &skp.Double{
			Data:    shoulderAlignment,
			Warning: shoulderAlignmentWarning,
		},
	}
	return dtlGolfSetupPoints
}

//spine angle
//line from midhip to neck
//angle between that and vertical axis
func GetSpineAngle(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate spine angle without axes calibration",
		}
	}
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.Midhip, "midhip", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.Neck, "neck", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	vertAxisLine := calibrationInfo.VertAxisLine
	vertAxisThroughMidhipLine := util.GetLineWithSlope(util.ConvertCvKeypointToPoint(keypoints.Midhip), vertAxisLine.Slope)
	neckPoint := util.ConvertCvKeypointToPoint(keypoints.Neck)
	xOnVertAxis := (keypoints.Neck.Y - vertAxisThroughMidhipLine.YIntercept) / vertAxisThroughMidhipLine.Slope
	pointUpVertAxisSameHeightAsNeck := &util.Point{XPos: xOnVertAxis, YPos: keypoints.Neck.Y}
	midhipPoint := util.ConvertCvKeypointToPoint(keypoints.Midhip)
	angleAtIntersect := util.GetAngleAtIntersection(neckPoint, midhipPoint, pointUpVertAxisSameHeightAsNeck)
	return angleAtIntersect, warning
}

//feet alignment
//assume feet are left of vert axis (maybe use toes? easier to see?)
//TODO: add other edge cases for feet crossing vertaxis
func GetFeetAlignment(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate feet alignment without axes calibration",
		}
	}
	if calibrationInfo.CalibrationType == skp.CalibrationType_AXES_CALIBRATION_ONLY {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate feet alignment without vanishing point calibration",
		}
	}
	var warning util.Warning
	currFeetLine, w := util.GetFeetLine(keypoints, calibrationInfo.FeetLine.FeetLineMethod)
	if w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	feetDegrees := util.GetAngleAtIntersection(&currFeetLine.LPoint, &currFeetLine.RPoint, &calibrationInfo.VanishingPoint)
	realParallelLine := util.GetLine(&calibrationInfo.VanishingPoint, &currFeetLine.RPoint) // line that converges at vanishing point
	//determine if closed or open
	if currFeetLine.Line.Slope == realParallelLine.Slope { // neutral
		return 0, warning
	} else if currFeetLine.Line.Slope > realParallelLine.Slope && currFeetLine.Line.Slope <= 0 { // closed (will return positive number)
		return feetDegrees, warning
	} else { // open (will return negative number)
		return float64(-1) * feetDegrees, warning
	}
}

//alignents
//all relative to heel alignment (open to heels, closed to heels)

//knee bend
//line from rhip to rknee
//line from rknee to rankle
//angle between those lines
func GetKneeBend(keypoints *skp.Body25PoseKeypoints) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.RHip, "right hip", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.RKnee, "right knee", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.RAnkle, "right ankle", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	kneeBend := util.GetAngleAtIntersection(util.ConvertCvKeypointToPoint(keypoints.RHip), util.ConvertCvKeypointToPoint(keypoints.RKnee), util.ConvertCvKeypointToPoint(keypoints.RAnkle))
	return kneeBend, warning
}

//waist bend
//line between rhip to neck
//line between rhip to rknee
//angle between those lines

//shoulder alignment
//find vanishing point
//any line that goes to vanishing point
//line from rshoulder to lshoulder
//find difference in angles
func GetShoulderAlignment(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate feet alignment without axes calibration",
		}
	}
	if calibrationInfo.CalibrationType == skp.CalibrationType_AXES_CALIBRATION_ONLY {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate feet alignment without vanishing point calibration",
		}
	}
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.LShoulder, "left shoulder", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.RShoulder, "right shoulder", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	fmt.Printf("LShoulder: %+v, RSHoulder: %+v, VanishingPoint: %+v\n", *keypoints.LShoulder, *keypoints.RShoulder, calibrationInfo.VanishingPoint)
	currShoulderLine := util.GetLine(util.ConvertCvKeypointToPoint(keypoints.LShoulder), util.ConvertCvKeypointToPoint(keypoints.RShoulder))
	shoulderDegrees := util.GetAngleAtIntersection(util.ConvertCvKeypointToPoint(keypoints.LShoulder), util.ConvertCvKeypointToPoint(keypoints.RShoulder), &calibrationInfo.VanishingPoint)
	fmt.Printf("ShoulderDegrees: %f\n", shoulderDegrees)
	realParallelLine := util.GetLine(&calibrationInfo.VanishingPoint, util.ConvertCvKeypointToPoint(keypoints.RShoulder)) // line that converges at vanishing point
	//determine if closed or open
	if currShoulderLine.Slope == realParallelLine.Slope { // neutral
		return 0, warning
	} else if currShoulderLine.Slope > realParallelLine.Slope && currShoulderLine.Slope <= 0 { // closed (will return positive number)
		return shoulderDegrees, warning
	} else { // open (will return negative number)
		return float64(-1) * shoulderDegrees, warning
	}
}
