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
//relative to hip to neck length
//line from left heel to right heel
//ratio of that line to left shoulder to right shoulder line
//greater than 1 is wider than shoulder width, less than 1 is less wide than shoulder width
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
	shoulderWidth := util.GetLengthBetweenTwoPoints(util.ConvertCvKeypointToPoint(keypoints.LShoulder), util.ConvertCvKeypointToPoint(keypoints.RShoulder))
	stanceWidth := util.GetLengthBetweenTwoPoints(util.ConvertCvKeypointToPoint(keypoints.LHeel), util.ConvertCvKeypointToPoint(keypoints.RHeel))
	return stanceWidth / shoulderWidth, warning
}

//shoulder tilt
//relative to horizontal axes slope
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
