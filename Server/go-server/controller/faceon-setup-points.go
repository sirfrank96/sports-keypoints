package controller

import (
	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	"github.com/sirfrank96/go-server/util"
)

//assuming right handed golfer

func VerifyFaceOnCalibrationImage(keypoints *cv.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (*util.CalibrationInfo, util.Warning) {
	return util.VerifyCalibrationImageAxes(keypoints, calibrationInfo)
}

//side bend
//line from midhip to neck
//angle of intersect between that and vertical axis through midhip
func GetSideBend(keypoints *cv.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	// TODO: IF calibrationInfo.AxesWarning is not nil return that
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.Midhip, "midhip", 0.5); w != nil {
		if w.GetWarningType() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.Neck, "neck", 0.5); w != nil {
		if w.GetWarningType() == util.SEVERE {
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
func getFootFlare(heel *cv.Keypoint, toe *cv.Keypoint, calibrationInfo *util.CalibrationInfo, midpoint *util.Point) float64 {
	vertAxisThroughMidpoint := util.GetLineWithSlope(midpoint, calibrationInfo.VertAxisLine.Slope)
	toeToHeelLine := util.GetLine(util.ConvertCvKeypointToPoint(toe), util.ConvertCvKeypointToPoint(heel))
	intersection := util.GetIntersection(toeToHeelLine, vertAxisThroughMidpoint)
	if intersection.IntersectPoint.YPos > toe.Y { // internal foot
		return float64(-1) * intersection.AngleAtIntersect
	} else { // external foot
		return intersection.AngleAtIntersect
	}
}

func GetLeftFootFlare(keypoints *cv.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	// TODO: IF calibrationInfo.AxesWarning is not nil return that
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.LHeel, "left heel", 0.5); w != nil {
		if w.GetWarningType() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.RHeel, "right heel", 0.5); w != nil {
		if w.GetWarningType() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.LBigToe, "left big toe", 0.5); w != nil {
		if w.GetWarningType() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	heelsMidpoint := util.GetMidpoint(util.ConvertCvKeypointToPoint(keypoints.LHeel), util.ConvertCvKeypointToPoint(keypoints.RHeel))
	return getFootFlare(keypoints.LHeel, keypoints.LBigToe, calibrationInfo, heelsMidpoint), warning
}

func GetRightFootFlare(keypoints *cv.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	// TODO: IF calibrationInfo.AxesWarning is not nil return that
	var warning util.Warning
	if w := util.VerifyKeypoint(keypoints.LHeel, "left heel", 0.5); w != nil {
		if w.GetWarningType() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.RHeel, "right heel", 0.5); w != nil {
		if w.GetWarningType() == util.SEVERE {
			return 0, w
		}
		warning = util.AppendMinorWarnings(warning, w)
	}
	if w := util.VerifyKeypoint(keypoints.RBigToe, "right big toe", 0.5); w != nil {
		if w.GetWarningType() == util.SEVERE {
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
//ratio of that line to hip to neck length
