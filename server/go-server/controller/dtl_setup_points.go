package controller

import (
	"context"
	"fmt"
	"math"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

// assuming right handed golfer

// 2 calibration images - 1 for perpendicular axes, 1 for vanishing points
// img 1: stand straddled, check to make sure heel horizontal and spine vertical are close to 90
// img 2: set alignment stick not centered, point alignment stick at target, set up with heels against alignment stick feet shoulder width or wider (check that heels are not centered in image)
// get vanishing point, intersection of vertaxis and heels axis

func CalculateDTLSetupPoints(ctx context.Context, bodyDatapoints *skp.Body25PoseDatapoints, golfSpecificDatapoints *skp.GolfSpecificDatapoints, calibrationInfo *util.CalibrationInfo) *skp.DTLGolfSetupPoints {
	fmt.Printf("Calculating Dtl setup points. BodyDatapoints: %+v\n GolfSpecificDatapoints: %+v\n CalibrationInfo: %+v\n", bodyDatapoints, golfSpecificDatapoints, calibrationInfo)
	spineAngle, warning := GetSpineAngle(bodyDatapoints, calibrationInfo)
	var spineAngleWarning string
	if warning != nil {
		spineAngleWarning = warning.Error()
	}
	fmt.Printf("Spine angle is %f\n", spineAngle)
	feetAlignment, warning := GetFeetAlignment(bodyDatapoints, calibrationInfo)
	var feetAlignmentWarning string
	if warning != nil {
		feetAlignmentWarning = warning.Error()
	}
	fmt.Printf("Feet alignment is %f\n", feetAlignment)
	heelAlignment, warning := GetHeelAlignment(bodyDatapoints, calibrationInfo)
	var heelAlignmentWarning string
	if warning != nil {
		heelAlignmentWarning = warning.Error()
	}
	fmt.Printf("Heel alignment is %f\n", heelAlignment)
	toeAlignment, warning := GetToeAlignment(bodyDatapoints, calibrationInfo)
	var toeAlignmentWarning string
	if warning != nil {
		toeAlignmentWarning = warning.Error()
	}
	fmt.Printf("Toe alignment is %f\n", toeAlignment)
	shoulderAlignment, warning := GetShoulderAlignment(bodyDatapoints, golfSpecificDatapoints, calibrationInfo)
	var shoulderAlignmentWarning string
	if warning != nil {
		shoulderAlignmentWarning = warning.Error()
	}
	fmt.Printf("Shoulder alignment is %f\n", shoulderAlignment)
	waistAlignment, warning := GetWaistAlignment(bodyDatapoints, calibrationInfo)
	var waistAlignmentWarning string
	if warning != nil {
		waistAlignmentWarning = warning.Error()
	}
	fmt.Printf("Waist alignment is %f\n", waistAlignment)
	kneeBend, warning := GetKneeBend(bodyDatapoints)
	var kneeBendWarning string
	if warning != nil {
		kneeBendWarning = warning.Error()
	}
	fmt.Printf("Knee bend is %f\n", kneeBend)
	distanceFromBall, warning := GetDistanceFromBall(bodyDatapoints, golfSpecificDatapoints)
	var distanceFromBallWarning string
	if warning != nil {
		distanceFromBallWarning = warning.Error()
	}
	fmt.Printf("Distance from ball is %f\n", distanceFromBall)
	ulnarDeviation, warning := GetUlnarDeviation(bodyDatapoints, golfSpecificDatapoints)
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
func GetSpineAngle(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate spine angle without axes calibration",
		}
	}
	var warning util.Warning
	if w := util.VerifyDatapoint(bodyDatapoints.Midhip, "midhip", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyDatapoint(bodyDatapoints.Neck, "neck", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	// convert datapoints to point
	midhip := util.ConvertDatapointToPoint(bodyDatapoints.Midhip)
	neck := util.ConvertDatapointToPoint(bodyDatapoints.Neck)
	// calculate spine angle
	vertAxisThroughMidhipLine := util.GetLineWithSlope(midhip, calibrationInfo.VertAxisLine.Slope)
	pointUpVertAxisSameHeightAsNeck := util.GetPointOnLineWithY(neck.YPos, vertAxisThroughMidhipLine)
	angleAtIntersect := util.GetAngleAtIntersection(neck, midhip, pointUpVertAxisSameHeightAsNeck)
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
func GetFeetAlignment(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
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
	currFeetLine, warning := util.GetFeetLine(bodyDatapoints, calibrationInfo.FeetLineMethod)
	if warning != nil && warning.GetSeverity() == util.SEVERE {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  fmt.Sprintf("Can't calculate feet alignment: %s", warning.Error()),
		}
	}
	return getFeetAlignmentHelper(currFeetLine, calibrationInfo), nil
}

func GetHeelAlignment(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
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
	heelLine, warning := util.GetFeetLine(bodyDatapoints, skp.FeetLineMethod_USE_HEEL_LINE)
	if warning != nil && warning.GetSeverity() == util.SEVERE {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  fmt.Sprintf("Can't calculate heel alignment: %s", warning.Error()),
		}
	}
	return getFeetAlignmentHelper(heelLine, calibrationInfo), nil
}

func GetToeAlignment(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
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
	toeLine, warning := util.GetFeetLine(bodyDatapoints, skp.FeetLineMethod_USE_TOE_LINE)
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
func GetShoulderAlignment(bodyDatapoints *skp.Body25PoseDatapoints, golfSpecificDatapoints *skp.GolfSpecificDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
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
	// add shoulder tilt for shoulder alignment calculation if provided
	if golfSpecificDatapoints.ShoulderTilt == nil {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate shoulder alignment without shoulder tilt in golf specific datapoints",
		}
	}
	var warning util.Warning
	if w := util.VerifyDatapoint(bodyDatapoints.LShoulder, "left shoulder", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyDatapoint(bodyDatapoints.RShoulder, "right shoulder", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	// convert datapoints to point
	lshoulder := util.ConvertDatapointToPoint(bodyDatapoints.LShoulder)
	rshoulder := util.ConvertDatapointToPoint(bodyDatapoints.RShoulder)
	// find auxillary vanishing point given shoulder tilt
	// shoulder tilt will raise or lower the vanishing point of lines parallel to the shoulder line on the vertical axis
	shoulderTiltRad := util.ConvertDegreesToRad(golfSpecificDatapoints.ShoulderTilt.Data)
	xtemp := calibrationInfo.VanishingPoint.XPos - rshoulder.XPos
	ytemp := calibrationInfo.VanishingPoint.YPos - rshoulder.YPos
	rotatedXTemp := (xtemp * math.Cos(shoulderTiltRad)) - (ytemp * math.Sin(shoulderTiltRad))
	rotatedYTemp := (xtemp * math.Sin(shoulderTiltRad)) - (ytemp * math.Cos(shoulderTiltRad))
	rotatedX := rotatedXTemp + rshoulder.XPos
	rotatedY := rotatedYTemp + rshoulder.YPos
	rotatedPoint := &util.Point{XPos: rotatedX, YPos: rotatedY}
	lineFromRShoulderToRotatedPoint := util.GetLine(rshoulder, rotatedPoint)
	intersection := util.GetIntersection(lineFromRShoulderToRotatedPoint, &calibrationInfo.VertAxisLine)
	avp := intersection.IntersectPoint
	// get the signed angle of rotation from line from rshoulder to auxillary vanishing point to the line from rshoulder to lshoulder
	vectFromRShoulderToAvp := util.GetVector(&avp, rshoulder)
	vectFromRShoulderToLShoulder := util.GetVector(lshoulder, rshoulder)
	return util.GetSignedAngleOfRotation(vectFromRShoulderToAvp, vectFromRShoulderToLShoulder), warning
}

// waist alignment
// find vanishing point
// any line that goes to vanishing point
// line from rhip to lhip
// find difference in angles
func GetWaistAlignment(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
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
	if w := util.VerifyDatapoint(bodyDatapoints.LHip, "left hip", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyDatapoint(bodyDatapoints.RHip, "right hip", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	// convert datapoints to point
	lhip := util.ConvertDatapointToPoint(bodyDatapoints.LHip)
	rhip := util.ConvertDatapointToPoint(bodyDatapoints.RHip)
	//get the signed angle of rotation from line from rhip to vanishing point to the line from rhip to lhip
	vectFromRHipToVp := util.GetVector(&calibrationInfo.VanishingPoint, rhip)
	vectFromRHipToLHip := util.GetVector(lhip, rhip)
	return util.GetSignedAngleOfRotation(vectFromRHipToVp, vectFromRHipToLHip), warning
}

// knee bend
// line from rhip to rknee
// line from rknee to rankle
// 180 - angle between those lines (ie. angle away from straight legs)
func GetKneeBend(bodyDatapoints *skp.Body25PoseDatapoints) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyDatapoint(bodyDatapoints.RHip, "right hip", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyDatapoint(bodyDatapoints.RKnee, "right knee", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyDatapoint(bodyDatapoints.RAnkle, "right ankle", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	// convert datapoints to point
	rhip := util.ConvertDatapointToPoint(bodyDatapoints.RHip)
	rknee := util.ConvertDatapointToPoint(bodyDatapoints.RKnee)
	rankle := util.ConvertDatapointToPoint(bodyDatapoints.RAnkle)
	// calculate knee bend
	kneeBend := util.GetAngleAtIntersection(rhip, rknee, rankle)
	return 180.0 - kneeBend, warning
}

// distance from ball
// line perpendicular to toeline that intersects ball
// line from midhip to neck
// ratio between two lengths
// the larger the number the farther away from ball
func GetDistanceFromBall(bodyDatapoints *skp.Body25PoseDatapoints, golfSpecificDatapoints *skp.GolfSpecificDatapoints) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyDatapoint(bodyDatapoints.Midhip, "midhip", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyDatapoint(bodyDatapoints.Neck, "neck", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyDatapoint(golfSpecificDatapoints.GolfBall, "golf ball", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	toeLine, w := util.GetFeetLine(bodyDatapoints, skp.FeetLineMethod_USE_TOE_LINE)
	if w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	// convert datapoints to point
	midhip := util.ConvertDatapointToPoint(bodyDatapoints.Midhip)
	neck := util.ConvertDatapointToPoint(bodyDatapoints.Neck)
	golfball := util.ConvertDatapointToPoint(golfSpecificDatapoints.GolfBall)
	// calculate distance from ball
	projection := util.GetProjectionOntoLine(&toeLine.Line, golfball)
	lengthFromBall := util.GetLengthBetweenTwoPoints(&projection.IntersectPoint, &projection.OriginalPoint)
	lengthOfSpine := util.GetLengthBetweenTwoPoints(midhip, neck)
	return lengthFromBall / lengthOfSpine, warning
}

// ulnar deviation
// line from right elbow to right wrist
// line from wrist to club head
// angle between those lines
// the larger the number the more ulnar deviation (ie. higher hands)
func GetUlnarDeviation(bodyDatapoints *skp.Body25PoseDatapoints, golfSpecificDatapoints *skp.GolfSpecificDatapoints) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyDatapoint(bodyDatapoints.RElbow, "right elbow", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyDatapoint(bodyDatapoints.RWrist, "right wrist", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyDatapoint(golfSpecificDatapoints.ClubHead, "club head", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	// convert datapoints to point
	relbow := util.ConvertDatapointToPoint(bodyDatapoints.RElbow)
	rwrist := util.ConvertDatapointToPoint(bodyDatapoints.RWrist)
	clubhead := util.ConvertDatapointToPoint(golfSpecificDatapoints.ClubHead)
	// calculate ulnar deviation
	angle := util.GetAngleAtIntersection(relbow, rwrist, clubhead)
	return angle, warning
}
