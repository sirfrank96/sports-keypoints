package controller

import (
	"context"
	"fmt"
	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

//assuming right handed golfer

func VerifyFaceOnCalibrationImage(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (*util.CalibrationInfo, util.Warning) {
	return util.VerifyCalibrationImageAxes(keypoints, calibrationInfo)
}

func CalculateFaceOnSetupPoints(ctx context.Context, keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) *skp.FaceOnGolfSetupPoints {
	sideBend, warning := GetSideBend(keypoints, calibrationInfo)
	var sideBendWarning string
	if warning != nil {
		sideBendWarning = warning.Error()
	}
	fmt.Printf("Side bend is %f", sideBend)
	lFootFlare, warning := GetLeftFootFlare(keypoints, calibrationInfo)
	var lFootFlareWarning string
	if warning != nil {
		lFootFlareWarning = warning.Error()
	}
	fmt.Printf("Left foot flare is %f", lFootFlare)
	rFootFlare, warning := GetRightFootFlare(keypoints, calibrationInfo)
	var rFootFlareWarning string
	if warning != nil {
		rFootFlareWarning = warning.Error()
	}
	fmt.Printf("Right foot flare is %f", rFootFlare)
	stanceWidth, warning := GetStanceWidth(keypoints)
	var stanceWidthWarning string
	if warning != nil {
		stanceWidthWarning = warning.Error()
	}
	fmt.Printf("Stance width is %f", stanceWidth)
	shoulderTilt, warning := GetShoulderTilt(keypoints, calibrationInfo)
	var shoulderTiltWarning string
	if warning != nil {
		shoulderTiltWarning = warning.Error()
	}
	fmt.Printf("Shoulder tilt is %f", shoulderTilt)
	waistTilt, warning := GetWaistTilt(keypoints, calibrationInfo)
	var waistTiltWarning string
	if warning != nil {
		waistTiltWarning = warning.Error()
	}
	fmt.Printf("Waist tilt is %f", waistTilt)
	shaftLean, warning := GetShaftLean(calibrationInfo)
	var shaftLeanWarning string
	if warning != nil {
		shaftLeanWarning = warning.Error()
	}
	fmt.Printf("Shaft lean is %f", shaftLean)
	ballPosition, warning := GetBallPosition(calibrationInfo)
	var ballPositionWarning string
	if warning != nil {
		ballPositionWarning = warning.Error()
	}
	fmt.Printf("Ball position is %f", ballPosition)
	headPosition, warning := GetHeadPosition(keypoints, calibrationInfo)
	var headPositionWarning string
	if warning != nil {
		headPositionWarning = warning.Error()
	}
	fmt.Printf("Head position is %f", headPosition)
	chestPosition, warning := GetChestPosition(keypoints, calibrationInfo)
	var chestPositionWarning string
	if warning != nil {
		chestPositionWarning = warning.Error()
	}
	fmt.Printf("Chest position is %f", chestPosition)
	midHipPosition, warning := GetMidhipPosition(keypoints, calibrationInfo)
	var midHipPositionWarning string
	if warning != nil {
		midHipPositionWarning = warning.Error()
	}
	fmt.Printf("Mid hip position is %f", midHipPosition)

	faceOnGolfSetupPoints := &skp.FaceOnGolfSetupPoints{
		SideBend: &skp.Double{
			Data:    sideBend,
			Warning: sideBendWarning,
		},
		LFootFlare: &skp.Double{
			Data:    lFootFlare,
			Warning: lFootFlareWarning,
		},
		RFootFlare: &skp.Double{
			Data:    rFootFlare,
			Warning: rFootFlareWarning,
		},
		StanceWidth: &skp.Double{
			Data:    stanceWidth,
			Warning: stanceWidthWarning,
		},
		ShoulderTilt: &skp.Double{
			Data:    shoulderTilt,
			Warning: shoulderTiltWarning,
		},
		WaistTilt: &skp.Double{
			Data:    waistTilt,
			Warning: waistTiltWarning,
		},
		ShaftLean: &skp.Double{
			Data:    shaftLean,
			Warning: shaftLeanWarning,
		},
		BallPosition: &skp.Double{
			Data:    ballPosition,
			Warning: ballPositionWarning,
		},
		HeadPosition: &skp.Double{
			Data:    headPosition,
			Warning: headPositionWarning,
		},
		ChestPosition: &skp.Double{
			Data:    chestPosition,
			Warning: chestPositionWarning,
		},
		MidHipPosition: &skp.Double{
			Data:    midHipPosition,
			Warning: midHipPositionWarning,
		},
	}
	return faceOnGolfSetupPoints
}

//side bend
//line from midhip to neck
//angle of intersect between that and vertical axis through midhip
func GetSideBend(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate side bend without axes calibration",
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
	// determine if left or right side bend
	if keypoints.Neck.X < keypoints.Midhip.X { // right
		return angleAtIntersect, warning
	} else { // left
		return float64(-1) * angleAtIntersect, warning
	}
}

//foot flares
//line from heel to big toe
//angle of intersect vert axis through midpoint of heels
func getFootFlare(heel *skp.Keypoint, toe *skp.Keypoint, calibrationInfo *util.CalibrationInfo, midpoint *util.Point) float64 {
	vertAxisThroughMidpoint := util.GetLineWithSlope(midpoint, calibrationInfo.VertAxisLine.Slope)
	toeToHeelLine := util.GetLine(util.ConvertCvKeypointToPoint(toe), util.ConvertCvKeypointToPoint(heel))
	intersection := util.GetIntersection(toeToHeelLine, vertAxisThroughMidpoint)
	if intersection.IntersectPoint.YPos > toe.Y { // internal foot
		return float64(-1) * intersection.AngleAtIntersect
	} else { // external foot
		return intersection.AngleAtIntersect
	}
}

func GetLeftFootFlare(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate left foot flare without axes calibration",
		}
	}
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.LHeel, "left heel", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.RHeel, "right heel", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.LBigToe, "left big toe", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	heelsMidpoint := util.GetMidpoint(util.ConvertCvKeypointToPoint(keypoints.LHeel), util.ConvertCvKeypointToPoint(keypoints.RHeel))
	return getFootFlare(keypoints.LHeel, keypoints.LBigToe, calibrationInfo, heelsMidpoint), warning
}

func GetRightFootFlare(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate right foot flare without axes calibration",
		}
	}
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.LHeel, "left heel", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.RHeel, "right heel", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.RBigToe, "right big toe", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	heelsMidpoint := util.GetMidpoint(util.ConvertCvKeypointToPoint(keypoints.LHeel), util.ConvertCvKeypointToPoint(keypoints.RHeel))
	return getFootFlare(keypoints.RHeel, keypoints.RBigToe, calibrationInfo, heelsMidpoint), warning
}

//stance width
//line from left heel to right heel
//line from midhip to neck
//ratio between 2 lengths
//the larger the number the wider the stance
func GetStanceWidth(keypoints *skp.Body25PoseKeypoints) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.LHeel, "left heel", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.RHeel, "right heel", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
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
	lengthOfSpine := util.GetLengthBetweenTwoPoints(util.ConvertCvKeypointToPoint(keypoints.Midhip), util.ConvertCvKeypointToPoint(keypoints.Neck))
	stanceWidth := util.GetLengthBetweenTwoPoints(util.ConvertCvKeypointToPoint(keypoints.LHeel), util.ConvertCvKeypointToPoint(keypoints.RHeel))
	return stanceWidth / lengthOfSpine, warning
}

//shoulder tilt
//relative to horizontal axis slope
//positive angle for right shoulder lower than left, negative angle if right shoulder is higher than left
func GetShoulderTilt(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate shoulder tilt without axes calibration",
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
	shoulderLineSlope := util.GetSlope(util.ConvertCvKeypointToPoint(keypoints.RShoulder), util.ConvertCvKeypointToPoint(keypoints.LShoulder))
	shoulderLineDegrees := util.ConvertSlopeToDegrees(shoulderLineSlope)
	horAxisLineDegrees := util.ConvertSlopeToDegrees(calibrationInfo.HorAxisLine.Slope)
	fmt.Printf("horAxisLineDegrees: %f, shoulderLineDegrees: %f", horAxisLineDegrees, shoulderLineDegrees)
	diff := horAxisLineDegrees - shoulderLineDegrees
	return diff, warning
}

//waist tilt
//relative to horizontal axis slope
//positive angle for right hip lower than left, negative angle if right hip is higher than left
func GetWaistTilt(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate waist tilt without axes calibration",
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
	waistLineSlope := util.GetSlope(util.ConvertCvKeypointToPoint(keypoints.RHip), util.ConvertCvKeypointToPoint(keypoints.LHip))
	waistLineDegrees := util.ConvertSlopeToDegrees(waistLineSlope)
	horAxisLineDegrees := util.ConvertSlopeToDegrees(calibrationInfo.HorAxisLine.Slope)
	fmt.Printf("horAxisLineDegrees: %f, waistLineDegrees: %f", horAxisLineDegrees, waistLineDegrees)
	diff := horAxisLineDegrees - waistLineDegrees
	return diff, warning
}

//shaft lean
//line from club head to club butt
//relative to vertical axis slope
//positive angle is forward shaft lean, negative angle is backwards shaft lean
func GetShaftLean(calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate waist tilt without axes calibration",
		}
	}
	var warning util.Warning
	if w := util.VerifyKeypoint(&calibrationInfo.ClubButtPoint, "club butt", 0.5); w != nil {
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
	shaftSlope := util.GetSlope(util.ConvertCvKeypointToPoint(&calibrationInfo.ClubButtPoint), util.ConvertCvKeypointToPoint(&calibrationInfo.ClubHeadPoint))
	shaftDegrees := util.ConvertSlopeToDegrees(shaftSlope)
	vertAxisLineDegrees := util.ConvertSlopeToDegrees(calibrationInfo.VertAxisLine.Slope)
	fmt.Printf("vertAxisLineDegrees: %f, shaftDegrees: %f", vertAxisLineDegrees, shaftDegrees)
	diff := vertAxisLineDegrees - shaftDegrees
	return diff, warning
}

//ball position
//line perpendicular to feet line that goes through midpoint of feet
//line from midpoint of feet to ball
//angle between these lines
//positive angle means ball closer to lead side, negative angle means ball closer to trail side
func GetBallPosition(calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyKeypoint(&calibrationInfo.GolfBallPoint, "golf ball", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	feetLine := calibrationInfo.FeetLine
	feetLineMidpoint := util.GetMidpoint(&feetLine.LPoint, &feetLine.RPoint)
	feetLineSlopeRecipricol := util.GetRecipricol(feetLine.Line.Slope)
	linePerpendicularToFeetLine := util.GetLineWithSlope(feetLineMidpoint, feetLineSlopeRecipricol)
	pointOnPerpendicularLine := util.GetPointOnLineWithX(calibrationInfo.GolfBallPoint.X, linePerpendicularToFeetLine)
	angle := util.GetAngleAtIntersection(util.ConvertCvKeypointToPoint(&calibrationInfo.GolfBallPoint), feetLineMidpoint, pointOnPerpendicularLine)
	return angle, warning
}

//head position
//line perpendicular to feet line that goes through midpoint of feet
//line from midpoint of feet to nose
//angle between these lines
//positive angle means head is closer to lead side, negative angle means head is closer to trail side
func GetHeadPosition(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.Nose, "nose", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	feetLine := calibrationInfo.FeetLine
	feetLineMidpoint := util.GetMidpoint(&feetLine.LPoint, &feetLine.RPoint)
	feetLineSlopeRecipricol := util.GetRecipricol(feetLine.Line.Slope)
	linePerpendicularToFeetLine := util.GetLineWithSlope(feetLineMidpoint, feetLineSlopeRecipricol)
	pointOnPerpendicularLine := util.GetPointOnLineWithX(calibrationInfo.GolfBallPoint.X, linePerpendicularToFeetLine)
	angle := util.GetAngleAtIntersection(util.ConvertCvKeypointToPoint(keypoints.Nose), feetLineMidpoint, pointOnPerpendicularLine)
	return angle, warning
}

//chest position
//line perpendicular to feet line that goes through midpoint of feet
//line from midpoint of feet to neck
//angle between these lines
//positive angle means head is closer to lead side, negative angle means head is closer to trail side
func GetChestPosition(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.Neck, "neck", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	feetLine := calibrationInfo.FeetLine
	feetLineMidpoint := util.GetMidpoint(&feetLine.LPoint, &feetLine.RPoint)
	feetLineSlopeRecipricol := util.GetRecipricol(feetLine.Line.Slope)
	linePerpendicularToFeetLine := util.GetLineWithSlope(feetLineMidpoint, feetLineSlopeRecipricol)
	pointOnPerpendicularLine := util.GetPointOnLineWithX(calibrationInfo.GolfBallPoint.X, linePerpendicularToFeetLine)
	angle := util.GetAngleAtIntersection(util.ConvertCvKeypointToPoint(keypoints.Neck), feetLineMidpoint, pointOnPerpendicularLine)
	return angle, warning
}

//midhip position
//line perpendicular to feet line that goes through midpoint of feet
//line from midpoint of feet to neck
//angle between these lines
//positive angle means head is closer to lead side, negative angle means head is closer to trail side
func GetMidhipPosition(keypoints *skp.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.Midhip, "mid hip", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	feetLine := calibrationInfo.FeetLine
	feetLineMidpoint := util.GetMidpoint(&feetLine.LPoint, &feetLine.RPoint)
	feetLineSlopeRecipricol := util.GetRecipricol(feetLine.Line.Slope)
	linePerpendicularToFeetLine := util.GetLineWithSlope(feetLineMidpoint, feetLineSlopeRecipricol)
	pointOnPerpendicularLine := util.GetPointOnLineWithX(calibrationInfo.GolfBallPoint.X, linePerpendicularToFeetLine)
	angle := util.GetAngleAtIntersection(util.ConvertCvKeypointToPoint(keypoints.Midhip), feetLineMidpoint, pointOnPerpendicularLine)
	return angle, warning
}
