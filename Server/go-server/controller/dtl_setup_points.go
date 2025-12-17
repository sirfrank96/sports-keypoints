package controller

import (
	"context"
	"fmt"
	"math"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

// assuming right handed golfer

// TODO: CONVERTKEYPOINTTOPOINT right away so dont have to convert everytime pass to func

// 2 calibration images? 1 for perpendicular axes, 1 for vanishing points
// img 1: stand straddled, check to make sure heel horizontal and spine vertical are close to 90
// img 2: set alignment stick not centered, point alignment stick at target, set up with heels against alignment stick feet shoulder width or wider (check that heels are not centered in image)
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
	fmt.Printf("Calculating Dtl setup points. Keypoints: %+v\n CalibrationInfo: %+v\n", keypoints, calibrationInfo)
	spineAngle, warning := GetSpineAngle(keypoints, calibrationInfo)
	var spineAngleWarning string
	if warning != nil {
		spineAngleWarning = warning.Error()
	}
	fmt.Printf("Spine angle is %f\n", spineAngle)
	feetAlignment, warning := GetFeetAlignment(keypoints, calibrationInfo)
	var feetAlignmentWarning string
	if warning != nil {
		feetAlignmentWarning = warning.Error()
	}
	fmt.Printf("Feet alignment is %f\n", feetAlignment)
	heelAlignment, warning := GetHeelAlignment(keypoints, calibrationInfo)
	var heelAlignmentWarning string
	if warning != nil {
		heelAlignmentWarning = warning.Error()
	}
	fmt.Printf("Heel alignment is %f\n", heelAlignment)
	toeAlignment, warning := GetToeAlignment(keypoints, calibrationInfo)
	var toeAlignmentWarning string
	if warning != nil {
		toeAlignmentWarning = warning.Error()
	}
	fmt.Printf("Toe alignment is %f\n", toeAlignment)
	shoulderAlignment, warning := GetShoulderAlignment(keypoints, calibrationInfo)
	var shoulderAlignmentWarning string
	if warning != nil {
		shoulderAlignmentWarning = warning.Error()
	}
	fmt.Printf("Shoulder alignment is %f\n", shoulderAlignment)
	waistAlignment, warning := GetWaistAlignment(keypoints, calibrationInfo)
	var waistAlignmentWarning string
	if warning != nil {
		waistAlignmentWarning = warning.Error()
	}
	fmt.Printf("Waist alignment is %f\n", waistAlignment)
	kneeBend, warning := GetKneeBend(keypoints)
	var kneeBendWarning string
	if warning != nil {
		kneeBendWarning = warning.Error()
	}
	fmt.Printf("Knee bend is %f\n", kneeBend)
	distanceFromBall, warning := GetDistanceFromBall(keypoints, calibrationInfo)
	var distanceFromBallWarning string
	if warning != nil {
		distanceFromBallWarning = warning.Error()
	}
	fmt.Printf("Distance from ball is %f\n", distanceFromBall)
	ulnarDeviation, warning := GetUlnarDeviation(keypoints, calibrationInfo)
	var ulnarDeviationWarning string
	if warning != nil {
		ulnarDeviationWarning = warning.Error()
	}
	fmt.Printf("Ulnar deviation is %f\n", ulnarDeviation)

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
		WaistAlignment: &skp.Double{
			Data:    waistAlignment,
			Warning: waistAlignmentWarning,
		},
		DistanceFromBall: &skp.Double{
			Data:    distanceFromBall,
			Warning: distanceFromBallWarning,
		},
		UlnarDeviation: &skp.Double{
			Data:    ulnarDeviation,
			Warning: ulnarDeviationWarning,
		},
	}
	return dtlGolfSetupPoints
}

// spine angle
// line from midhip to neck
// angle between that and vertical axis
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
	vertAxisThroughMidhipLine := util.GetLineWithSlope(util.ConvertKeypointToPoint(keypoints.Midhip), vertAxisLine.Slope)
	neckPoint := util.ConvertKeypointToPoint(keypoints.Neck)
	pointUpVertAxisSameHeightAsNeck := util.GetPointOnLineWithY(neckPoint.YPos, vertAxisThroughMidhipLine)
	midhipPoint := util.ConvertKeypointToPoint(keypoints.Midhip)
	angleAtIntersect := util.GetAngleAtIntersection(neckPoint, midhipPoint, pointUpVertAxisSameHeightAsNeck)
	return angleAtIntersect, warning
}

func getFeetAlignmentHelper(currFeetLine *util.FeetLine, calibrationInfo *util.CalibrationInfo) float64 {
	//get the signed angle of rotation from line from rpoint to vanishing point to the line from rpoint to lpoint
	vectFromRPointToVp := util.GetVector(&calibrationInfo.VanishingPoint, &currFeetLine.RPoint)
	vectFromRPointToLPoint := util.GetVector(&currFeetLine.LPoint, &currFeetLine.RPoint)
	return util.GetSignedAngleOfRotation(vectFromRPointToVp, vectFromRPointToLPoint)
}

// feet alignment
// assume feet are left of vert axis (maybe use toes? easier to see?)
// TODO: add other edge cases for feet crossing vertaxis
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
	currFeetLine, warning := util.GetFeetLine(keypoints, calibrationInfo.FeetLineMethod)
	if warning != nil && warning.GetSeverity() == util.SEVERE {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  fmt.Sprintf("Can't calculate feet alignment: %s", warning.Error()),
		}
	}
	return getFeetAlignmentHelper(currFeetLine, calibrationInfo), nil
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

// shoulder alignment
// find vanishing point
// any line that goes to vanishing point
// line from rshoulder to lshoulder
// find difference in angles
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
	if w := util.VerifyDouble(&calibrationInfo.ShoulderTilt); w != nil {
		return 0, w
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
	// find auxillary vanishing point given shoulder tilt
	// shoulder tilt will raise or lower the vanishing point of lines parallel to the shoulder line on the vertical axis
	shoulderTiltRad := util.ConvertDegreesToRad(calibrationInfo.ShoulderTilt.Data)
	xtemp := calibrationInfo.VanishingPoint.XPos - keypoints.RShoulder.X
	ytemp := calibrationInfo.VanishingPoint.YPos - keypoints.RShoulder.Y
	rotatedXTemp := (xtemp * math.Cos(shoulderTiltRad)) - (ytemp * math.Sin(shoulderTiltRad))
	rotatedYTemp := (xtemp * math.Sin(shoulderTiltRad)) - (ytemp * math.Cos(shoulderTiltRad))
	rotatedX := rotatedXTemp + keypoints.RShoulder.X
	rotatedY := rotatedYTemp + keypoints.RShoulder.Y
	rotatedPoint := &util.Point{XPos: rotatedX, YPos: rotatedY}
	lineFromRShoulderToRotatedPoint := util.GetLine(util.ConvertKeypointToPoint(keypoints.RShoulder), rotatedPoint)
	intersection := util.GetIntersection(lineFromRShoulderToRotatedPoint, &calibrationInfo.VertAxisLine)
	avp := intersection.IntersectPoint
	// get the signed angle of rotation from line from rshoulder to auxillary vanishing point to the line from rshoulder to lshoulder
	vectFromRShoulderToAvp := util.GetVector(&avp, util.ConvertKeypointToPoint(keypoints.RShoulder))
	vectFromRShoulderToLShoulder := util.GetVector(util.ConvertKeypointToPoint(keypoints.LShoulder), util.ConvertKeypointToPoint(keypoints.RShoulder))
	return util.GetSignedAngleOfRotation(vectFromRShoulderToAvp, vectFromRShoulderToLShoulder), warning
}

// waist alignment
// find vanishing point
// any line that goes to vanishing point
// line from rhip to lhip
// find difference in angles
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
	//get the signed angle of rotation from line from rhip to vanishing point to the line from rhip to lhip
	vectFromRHipToVp := util.GetVector(&calibrationInfo.VanishingPoint, util.ConvertKeypointToPoint(keypoints.RHip))
	vectFromRHipToLHip := util.GetVector(util.ConvertKeypointToPoint(keypoints.LHip), util.ConvertKeypointToPoint(keypoints.RHip))
	return util.GetSignedAngleOfRotation(vectFromRHipToVp, vectFromRHipToLHip), warning
}

// knee bend
// line from rhip to rknee
// line from rknee to rankle
// 180 - angle between those lines (ie. angle away from straight legs)
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
	kneeBend := util.GetAngleAtIntersection(util.ConvertKeypointToPoint(keypoints.RHip), util.ConvertKeypointToPoint(keypoints.RKnee), util.ConvertKeypointToPoint(keypoints.RAnkle))
	return 180.0 - kneeBend, warning
}

// distance from ball
// line perpendicular to toeline that intersects ball
// line from midhip to neck
// ratio between two lengths
// the larger the number the farther away from ball
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
	if w := util.VerifyKeypoint(&calibrationInfo.GolfBallPoint, "golf ball", 0.5); w != nil {
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
	projection := util.GetProjectionOntoLine(&toeLine.Line, util.ConvertKeypointToPoint(&calibrationInfo.GolfBallPoint))
	lengthFromBall := util.GetLengthBetweenTwoPoints(&projection.IntersectPoint, &projection.OriginalPoint)
	lengthOfSpine := util.GetLengthBetweenTwoPoints(util.ConvertKeypointToPoint(keypoints.Midhip), util.ConvertKeypointToPoint(keypoints.Neck))
	return lengthFromBall / lengthOfSpine, warning
}

// ulnar deviation
// line from right elbow to right wrist
// line from wrist to club head
// angle between those lines
// the larger the number the more ulnar deviation (ie. higher hands)
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
	if w := util.VerifyKeypoint(&calibrationInfo.ClubHeadPoint, "club head", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	angle := util.GetAngleAtIntersection(util.ConvertKeypointToPoint(keypoints.RElbow), util.ConvertKeypointToPoint(keypoints.RWrist), util.ConvertKeypointToPoint(&calibrationInfo.ClubHeadPoint))
	return angle, warning
}
