package controller

import (
	"context"
	"fmt"
	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

// assuming right handed golfer

// TODO: Common utility funcs for similar funcs

func CalculateFaceOnSetupPoints(ctx context.Context, bodyDatapoints *skp.Body25PoseDatapoints, golfSpecificDatapoints *skp.GolfSpecificDatapoints, calibrationInfo *util.CalibrationInfo) *skp.FaceOnGolfSetupPoints {
	fmt.Printf("Calculating Face on setup points. BodyDatapoints: %#v\n GolfSpecificDatapoints: %#v\n CalibrationInfo: %#v\n", bodyDatapoints, golfSpecificDatapoints, calibrationInfo)
	sideBend, warning := GetSideBend(bodyDatapoints, calibrationInfo)
	var sideBendWarning string
	if warning != nil {
		sideBendWarning = warning.Error()
	}
	fmt.Printf("Side bend is %f\n", sideBend)
	lFootFlare, warning := GetLeftFootFlare(bodyDatapoints, calibrationInfo)
	var lFootFlareWarning string
	if warning != nil {
		lFootFlareWarning = warning.Error()
	}
	fmt.Printf("Left foot flare is %f\n", lFootFlare)
	rFootFlare, warning := GetRightFootFlare(bodyDatapoints, calibrationInfo)
	var rFootFlareWarning string
	if warning != nil {
		rFootFlareWarning = warning.Error()
	}
	fmt.Printf("Right foot flare is %f\n", rFootFlare)
	stanceWidth, warning := GetStanceWidth(bodyDatapoints)
	var stanceWidthWarning string
	if warning != nil {
		stanceWidthWarning = warning.Error()
	}
	fmt.Printf("Stance width is %f\n", stanceWidth)
	shoulderTilt, warning := GetShoulderTilt(bodyDatapoints, calibrationInfo)
	var shoulderTiltWarning string
	if warning != nil {
		shoulderTiltWarning = warning.Error()
	}
	fmt.Printf("Shoulder tilt is %f\n", shoulderTilt)
	waistTilt, warning := GetWaistTilt(bodyDatapoints, calibrationInfo)
	var waistTiltWarning string
	if warning != nil {
		waistTiltWarning = warning.Error()
	}
	fmt.Printf("Waist tilt is %f\n", waistTilt)
	shaftLean, warning := GetShaftLean(golfSpecificDatapoints, calibrationInfo)
	var shaftLeanWarning string
	if warning != nil {
		shaftLeanWarning = warning.Error()
	}
	fmt.Printf("Shaft lean is %f\n", shaftLean)
	ballPosition, warning := GetBallPosition(bodyDatapoints, golfSpecificDatapoints, calibrationInfo)
	var ballPositionWarning string
	if warning != nil {
		ballPositionWarning = warning.Error()
	}
	fmt.Printf("Ball position is %f\n", ballPosition)
	headPosition, warning := GetHeadPosition(bodyDatapoints, calibrationInfo)
	var headPositionWarning string
	if warning != nil {
		headPositionWarning = warning.Error()
	}
	fmt.Printf("Head position is %f\n", headPosition)
	chestPosition, warning := GetChestPosition(bodyDatapoints, calibrationInfo)
	var chestPositionWarning string
	if warning != nil {
		chestPositionWarning = warning.Error()
	}
	fmt.Printf("Chest position is %f\n", chestPosition)
	midHipPosition, warning := GetMidhipPosition(bodyDatapoints, calibrationInfo)
	var midHipPositionWarning string
	if warning != nil {
		midHipPositionWarning = warning.Error()
	}
	fmt.Printf("Mid hip position is %f\n", midHipPosition)

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

// side bend
// line from midhip to neck
// angle of intersect between that and vertical axis through midhip
func GetSideBend(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate side bend without axes calibration",
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
	// calculate side bend
	lineFromMidhipWithVertAxisSlope := util.GetLineWithSlope(midhip, calibrationInfo.VertAxisLine.Slope)
	pointOnLine := util.GetPointOnLineWithY(neck.YPos, lineFromMidhipWithVertAxisSlope)
	vectFromMidhipToPointOnLine := util.GetVector(pointOnLine, midhip)
	spineVect := util.GetVector(neck, midhip)
	return util.GetSignedAngleOfRotation(spineVect, vectFromMidhipToPointOnLine), warning
}

type Direction int

const (
	Left Direction = iota + 1
	Right
)

// foot flares
// line from heel to big toe
// relative to vert axis slope
func getFootFlare(heel *util.Point, toe *util.Point, calibrationInfo *util.CalibrationInfo, direction Direction) float64 {
	lineFromHeelWithVertAxisSlope := util.GetLineWithSlope(heel, calibrationInfo.VertAxisLine.Slope)
	pointOnLine := util.GetPointOnLineWithY(toe.YPos, lineFromHeelWithVertAxisSlope)
	vectFromHeelToPointOnLine := util.GetVector(pointOnLine, heel)
	heelToToeVect := util.GetVector(toe, heel)
	if direction == Right {
		return util.GetSignedAngleOfRotation(vectFromHeelToPointOnLine, heelToToeVect)
	} else {
		return util.GetSignedAngleOfRotation(heelToToeVect, vectFromHeelToPointOnLine)
	}
}

func GetLeftFootFlare(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate left foot flare without axes calibration",
		}
	}
	var warning util.Warning
	if w := util.VerifyDatapoint(bodyDatapoints.LHeel, "left heel", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyDatapoint(bodyDatapoints.LBigToe, "left big toe", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	return getFootFlare(util.ConvertDatapointToPoint(bodyDatapoints.LHeel), util.ConvertDatapointToPoint(bodyDatapoints.LBigToe), calibrationInfo, Left), warning
}

func GetRightFootFlare(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate right foot flare without axes calibration",
		}
	}
	var warning util.Warning
	if w := util.VerifyDatapoint(bodyDatapoints.RHeel, "right heel", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyDatapoint(bodyDatapoints.RBigToe, "right big toe", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	return getFootFlare(util.ConvertDatapointToPoint(bodyDatapoints.RHeel), util.ConvertDatapointToPoint(bodyDatapoints.RBigToe), calibrationInfo, Right), warning
}

// stance width
// line from left heel to right heel
// line from midhip to neck
// ratio between 2 lengths
// the larger the number the wider the stance
func GetStanceWidth(bodyDatapoints *skp.Body25PoseDatapoints) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyDatapoint(bodyDatapoints.LHeel, "left heel", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyDatapoint(bodyDatapoints.RHeel, "right heel", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
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
	lheel := util.ConvertDatapointToPoint(bodyDatapoints.LHeel)
	rheel := util.ConvertDatapointToPoint(bodyDatapoints.RHeel)
	// calculate stance width
	lengthOfSpine := util.GetLengthBetweenTwoPoints(midhip, neck)
	stanceWidth := util.GetLengthBetweenTwoPoints(lheel, rheel)
	return stanceWidth / lengthOfSpine, warning
}

// shoulder tilt
// relative to horizontal axis slope
// positive angle for right shoulder lower than left, negative angle if right shoulder is higher than left
func GetShoulderTilt(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate shoulder tilt without axes calibration",
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
	// calculate shoulder tilt
	lineFromRShoulderWithHorAxisSlope := util.GetLineWithSlope(rshoulder, calibrationInfo.HorAxisLine.Slope)
	pointOnLine := util.GetPointOnLineWithX(lshoulder.XPos, lineFromRShoulderWithHorAxisSlope)
	vectFromRShoulderToPointOnLine := util.GetVector(pointOnLine, rshoulder)
	shoulderVect := util.GetVector(lshoulder, rshoulder)
	return util.GetSignedAngleOfRotation(shoulderVect, vectFromRShoulderToPointOnLine), warning
}

// waist tilt
// relative to horizontal axis slope
// positive angle for right hip lower than left, negative angle if right hip is higher than left
func GetWaistTilt(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate waist tilt without axes calibration",
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
	// calculate waist tilt
	lineFromRHipWithHorAxisSlope := util.GetLineWithSlope(rhip, calibrationInfo.HorAxisLine.Slope)
	pointOnLine := util.GetPointOnLineWithX(lhip.XPos, lineFromRHipWithHorAxisSlope)
	vectFromRHipToPointOnLine := util.GetVector(pointOnLine, rhip)
	waistVect := util.GetVector(lhip, rhip)
	return util.GetSignedAngleOfRotation(waistVect, vectFromRHipToPointOnLine), warning
}

// shaft lean
// line from club head to club butt
// relative to vertical axis slope
// positive angle is forward shaft lean, negative angle is backwards shaft lean
func GetShaftLean(golfSpecificDatapoints *skp.GolfSpecificDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	if calibrationInfo.CalibrationType == skp.CalibrationType_NO_CALIBRATION {
		return 0, util.WarningImpl{
			Severity: util.MINOR,
			Message:  "Can't calculate waist tilt without axes calibration",
		}
	}
	var warning util.Warning
	if w := util.VerifyDatapoint(golfSpecificDatapoints.ClubButt, "club butt", 0.5); w != nil {
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
	clubbutt := util.ConvertDatapointToPoint(golfSpecificDatapoints.ClubButt)
	clubhead := util.ConvertDatapointToPoint(golfSpecificDatapoints.ClubHead)
	// calculate shaft lean
	lineFromClubheadWithVertAxisSlope := util.GetLineWithSlope(clubhead, calibrationInfo.VertAxisLine.Slope)
	pointOnLine := util.GetPointOnLineWithY(clubbutt.YPos, lineFromClubheadWithVertAxisSlope)
	vectFromClubheadToPointOnLine := util.GetVector(pointOnLine, clubhead)
	shaftVect := util.GetVector(clubbutt, clubhead)
	return util.GetSignedAngleOfRotation(vectFromClubheadToPointOnLine, shaftVect), warning
}

// ball position
// line perpendicular to feet line that goes through midpoint of feet
// line from midpoint of feet to ball
// angle between these lines
// positive angle means ball closer to lead side, negative angle means ball closer to trail side
func GetBallPosition(bodyDatapoints *skp.Body25PoseDatapoints, golfSpecificDatapoints *skp.GolfSpecificDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyDatapoint(golfSpecificDatapoints.GolfBall, "golf ball", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	// convert datapoints to point
	lFootPoint, _ := util.GetLeftFootPoint(bodyDatapoints, calibrationInfo.FeetLineMethod)
	rFootPoint, _ := util.GetRightFootPoint(bodyDatapoints, calibrationInfo.FeetLineMethod)
	lfoot := util.ConvertDatapointToPoint(lFootPoint)
	rfoot := util.ConvertDatapointToPoint(rFootPoint)
	golfball := util.ConvertDatapointToPoint(golfSpecificDatapoints.GolfBall)
	// calculate ball position
	feetLineMidpoint := util.GetMidpoint(lfoot, rfoot)
	feetLineSlopeRecipricol := util.GetSlopeRecipricol(lfoot, rfoot)
	linePerpendicularToFeetLine := util.GetLineWithSlope(feetLineMidpoint, feetLineSlopeRecipricol)
	pointOnPerpendicularLine := util.GetPointOnLineWithY(golfball.YPos, linePerpendicularToFeetLine)
	vectFromMidpointToPointOnPerpLine := util.GetVector(pointOnPerpendicularLine, feetLineMidpoint)
	vectFromMidpointToBall := util.GetVector(golfball, feetLineMidpoint)
	var ballPosition float64
	if golfball.YPos < feetLineMidpoint.YPos {
		ballPosition = util.GetSignedAngleOfRotation(vectFromMidpointToPointOnPerpLine, vectFromMidpointToBall)
	} else {
		ballPosition = util.GetSignedAngleOfRotation(vectFromMidpointToBall, vectFromMidpointToPointOnPerpLine)
	}
	return ballPosition, warning
}

// head position
// line perpendicular to feet line that goes through midpoint of feet
// line from midpoint of feet to nose
// angle between these lines
// positive angle means head is closer to lead side, negative angle means head is closer to trail side
func GetHeadPosition(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyDatapoint(bodyDatapoints.Nose, "nose", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	// convert datapoints to point
	lFootPoint, _ := util.GetLeftFootPoint(bodyDatapoints, calibrationInfo.FeetLineMethod)
	rFootPoint, _ := util.GetRightFootPoint(bodyDatapoints, calibrationInfo.FeetLineMethod)
	lfoot := util.ConvertDatapointToPoint(lFootPoint)
	rfoot := util.ConvertDatapointToPoint(rFootPoint)
	nose := util.ConvertDatapointToPoint(bodyDatapoints.Nose)
	// calculate head position
	feetLineMidpoint := util.GetMidpoint(lfoot, rfoot)
	feetLineSlopeRecipricol := util.GetSlopeRecipricol(lfoot, rfoot)
	linePerpendicularToFeetLine := util.GetLineWithSlope(feetLineMidpoint, feetLineSlopeRecipricol)
	pointOnPerpendicularLine := util.GetPointOnLineWithY(nose.YPos, linePerpendicularToFeetLine)
	vectFromMidpointToPointOnPerpLine := util.GetVector(pointOnPerpendicularLine, feetLineMidpoint)
	vectFromMidpointToHead := util.GetVector(nose, feetLineMidpoint)
	var headPosition float64
	if nose.YPos < feetLineMidpoint.YPos {
		headPosition = util.GetSignedAngleOfRotation(vectFromMidpointToPointOnPerpLine, vectFromMidpointToHead)
	} else {
		headPosition = util.GetSignedAngleOfRotation(vectFromMidpointToHead, vectFromMidpointToPointOnPerpLine)
	}
	return headPosition, warning
}

// chest position
// line perpendicular to feet line that goes through midpoint of feet
// line from midpoint of feet to neck
// angle between these lines
// positive angle means head is closer to lead side, negative angle means head is closer to trail side
func GetChestPosition(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyDatapoint(bodyDatapoints.Neck, "neck", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	// convert datapoints to point
	lFootPoint, _ := util.GetLeftFootPoint(bodyDatapoints, calibrationInfo.FeetLineMethod)
	rFootPoint, _ := util.GetRightFootPoint(bodyDatapoints, calibrationInfo.FeetLineMethod)
	lfoot := util.ConvertDatapointToPoint(lFootPoint)
	rfoot := util.ConvertDatapointToPoint(rFootPoint)
	neck := util.ConvertDatapointToPoint(bodyDatapoints.Neck)
	// calculate chest position
	feetLineMidpoint := util.GetMidpoint(lfoot, rfoot)
	feetLineSlopeRecipricol := util.GetSlopeRecipricol(lfoot, rfoot)
	linePerpendicularToFeetLine := util.GetLineWithSlope(feetLineMidpoint, feetLineSlopeRecipricol)
	pointOnPerpendicularLine := util.GetPointOnLineWithY(neck.YPos, linePerpendicularToFeetLine)
	vectFromMidpointToPointOnPerpLine := util.GetVector(pointOnPerpendicularLine, feetLineMidpoint)
	vectFromMidpointToChest := util.GetVector(neck, feetLineMidpoint)
	var chestPosition float64
	if neck.YPos < feetLineMidpoint.YPos {
		chestPosition = util.GetSignedAngleOfRotation(vectFromMidpointToPointOnPerpLine, vectFromMidpointToChest)
	} else {
		chestPosition = util.GetSignedAngleOfRotation(vectFromMidpointToChest, vectFromMidpointToPointOnPerpLine)
	}
	return chestPosition, warning
}

// midhip position
// line perpendicular to feet line that goes through midpoint of feet
// line from midpoint of feet to neck
// angle between these lines
// positive angle means head is closer to lead side, negative angle means head is closer to trail side
func GetMidhipPosition(bodyDatapoints *skp.Body25PoseDatapoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	var warning util.Warning
	if w := util.VerifyDatapoint(bodyDatapoints.Midhip, "mid hip", 0.5); w != nil {
		if w.GetSeverity() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	// convert datapoints to point
	lFootPoint, _ := util.GetLeftFootPoint(bodyDatapoints, calibrationInfo.FeetLineMethod)
	rFootPoint, _ := util.GetRightFootPoint(bodyDatapoints, calibrationInfo.FeetLineMethod)
	lfoot := util.ConvertDatapointToPoint(lFootPoint)
	rfoot := util.ConvertDatapointToPoint(rFootPoint)
	midhip := util.ConvertDatapointToPoint(bodyDatapoints.Midhip)
	// calculate midhip position
	feetLineMidpoint := util.GetMidpoint(lfoot, rfoot)
	feetLineSlopeRecipricol := util.GetSlopeRecipricol(lfoot, rfoot)
	linePerpendicularToFeetLine := util.GetLineWithSlope(feetLineMidpoint, feetLineSlopeRecipricol)
	pointOnPerpendicularLine := util.GetPointOnLineWithY(midhip.YPos, linePerpendicularToFeetLine)
	vectFromMidpointToPointOnPerpLine := util.GetVector(pointOnPerpendicularLine, feetLineMidpoint)
	vectFromMidpointToMidhip := util.GetVector(midhip, feetLineMidpoint)
	var midhipPosition float64
	if midhip.YPos < feetLineMidpoint.YPos {
		midhipPosition = util.GetSignedAngleOfRotation(vectFromMidpointToPointOnPerpLine, vectFromMidpointToMidhip)
	} else {
		midhipPosition = util.GetSignedAngleOfRotation(vectFromMidpointToMidhip, vectFromMidpointToPointOnPerpLine)
	}
	return midhipPosition, warning
}
