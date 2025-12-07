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
	heelAlignment, warning := GetHeelAlignment(keypoints, calibrationInfo)
	var heelAlignmentWarning string
	if warning != nil {
		heelAlignmentWarning = warning.Error()
	}
	fmt.Printf("Heel alignment is %f", heelAlignment)
	toeAlignment, warning := GetToeAlignment(keypoints, calibrationInfo)
	var toeAlignmentWarning string
	if warning != nil {
		toeAlignmentWarning = warning.Error()
	}
	fmt.Printf("Toe alignment is %f", toeAlignment)
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
		HeelAlignment: &skp.Double{
			Data:    heelAlignment,
			Warning: heelAlignmentWarning,
		},
		ToeAlignment: &skp.Double{
			Data:    toeAlignment,
			Warning: toeAlignmentWarning,
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
	pointUpVertAxisSameHeightAsNeck := util.GetPointOnLineWithY(neckPoint.YPos, vertAxisThroughMidhipLine)
	midhipPoint := util.ConvertCvKeypointToPoint(keypoints.Midhip)
	angleAtIntersect := util.GetAngleAtIntersection(neckPoint, midhipPoint, pointUpVertAxisSameHeightAsNeck)
	return angleAtIntersect, warning
}

func getFeetAlignmentHelper(currFeetLine *util.FeetLine, calibrationInfo *util.CalibrationInfo) float64 {
	feetDegrees := util.GetAngleAtIntersection(&currFeetLine.LPoint, &currFeetLine.RPoint, &calibrationInfo.VanishingPoint)
	realParallelLine := util.GetLine(&calibrationInfo.VanishingPoint, &currFeetLine.RPoint) // line that converges at vanishing point
	//determine if closed or open TODO: Make util functions
	if currFeetLine.Line.Slope == realParallelLine.Slope { // neutral
		return 0
	} else if currFeetLine.Line.Slope > realParallelLine.Slope && currFeetLine.Line.Slope <= 0 { // closed (will return positive number)
		return feetDegrees
	} else { // open (will return negative number)
		return float64(-1) * feetDegrees
	}
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
	currFeetLine := calibrationInfo.FeetLine
	return getFeetAlignmentHelper(&currFeetLine, calibrationInfo), nil
}

func GetHeelAlignment(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate heel alignment without axes calibration",
		}
	}
	if calibrationInfo.CalibrationType == skp.CalibrationType_AXES_CALIBRATION_ONLY {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate heel alignment without vanishing point calibration",
		}
	}
	heelLine, warning := util.GetFeetLine(keypoints, skp.FeetLineMethod_USE_HEEL_LINE)
	if warning != nil && warning.GetSeverity() == util.SEVERE {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  fmt.Sprintf("Can't calculate heel alignment: %s", warning.Error()),
		}
	}
	return getFeetAlignmentHelper(heelLine, calibrationInfo), nil
}

func GetToeAlignment(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
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
	toeLine, warning := util.GetFeetLine(keypoints, skp.FeetLineMethod_USE_TOE_LINE)
	if warning != nil && warning.GetSeverity() == util.SEVERE {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  fmt.Sprintf("Can't calculate toe alignment: %s", warning.Error()),
		}
	}
	return getFeetAlignmentHelper(toeLine, calibrationInfo), nil
}

//shoulder alignment
//find vanishing point
//any line that goes to vanishing point
//line from rshoulder to lshoulder
//find difference in angles
func GetShoulderAlignment(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate shoulder alignment without axes calibration",
		}
	}
	if calibrationInfo.CalibrationType == skp.CalibrationType_AXES_CALIBRATION_ONLY {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate shoulder alignment without vanishing point calibration",
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

//waist alignment
//find vanishing point
//any line that goes to vanishing point
//line from rhip to lhip
//find difference in angles
func GetWaistAlignment(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate waise alignment without axes calibration",
		}
	}
	if calibrationInfo.CalibrationType == skp.CalibrationType_AXES_CALIBRATION_ONLY {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate waist alignment without vanishing point calibration",
		}
	}
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.LHip, "left hip", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.RHip, "right hip", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	fmt.Printf("LHip: %+v, RHip: %+v, VanishingPoint: %+v\n", *keypoints.LHip, *keypoints.RHip, calibrationInfo.VanishingPoint)
	currWaistLine := util.GetLine(util.ConvertCvKeypointToPoint(keypoints.LHip), util.ConvertCvKeypointToPoint(keypoints.RHip))
	waistDegrees := util.GetAngleAtIntersection(util.ConvertCvKeypointToPoint(keypoints.LHip), util.ConvertCvKeypointToPoint(keypoints.RHip), &calibrationInfo.VanishingPoint)
	fmt.Printf("WaistDegrees: %f\n", waistDegrees)
	realParallelLine := util.GetLine(&calibrationInfo.VanishingPoint, util.ConvertCvKeypointToPoint(keypoints.RHip)) // line that converges at vanishing point
	//determine if closed or open
	if currWaistLine.Slope == realParallelLine.Slope { // neutral
		return 0, warning
	} else if currWaistLine.Slope > realParallelLine.Slope && currWaistLine.Slope <= 0 { // closed (will return positive number)
		return waistDegrees, warning
	} else { // open (will return negative number)
		return float64(-1) * waistDegrees, warning
	}
}

//knee bend
//line from rhip to rknee
//line from rknee to rankle
//180 - angle between those lines (ie. angle away from straight legs)
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
	return 180.0 - kneeBend, warning
}

//distance from ball
//line perpendicular to toeline that intersects ball
//line from midhip to neck
//ratio between two lengths
//the larger the number the farther away from ball
func GetDistanceFromBall(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
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
	if w := util.VerifyKeypoint(util.ConvertPointToCvKeypoint(&calibrationInfo.GolfBallPoint), "golf ball", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	toeLine, w := util.GetFeetLine(keypoints, skp.FeetLineMethod_USE_TOE_LINE)
	if w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	projection := util.GetProjectionOntoLine(&toeLine.Line, &calibrationInfo.GolfBallPoint)
	lengthFromBall := util.GetLengthBetweenTwoPoints(&projection.IntersectPoint, &projection.OriginalPoint)
	lengthOfSpine := util.GetLengthBetweenTwoPoints(util.ConvertCvKeypointToPoint(keypoints.Midhip), util.ConvertCvKeypointToPoint(keypoints.Neck))
	return lengthFromBall / lengthOfSpine, warning
}

//ulnar deviation
//line from right elbow to right wrist
//line from wrist to club head
//angle between those lines
//the larger the number the more ulnar deviation (ie. higher hands)
func GetUlnarDeviation(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.RElbow, "right elbow", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.RWrist, "right wrist", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(util.ConvertPointToCvKeypoint(&calibrationInfo.ClubHeadPoint), "club head", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	angle := util.GetAngleAtIntersection(util.ConvertCvKeypointToPoint(keypoints.RElbow), util.ConvertCvKeypointToPoint(keypoints.RWrist), &calibrationInfo.ClubHeadPoint)
	return angle, warning
}
