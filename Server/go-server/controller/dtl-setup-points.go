package controller

import (
	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	"github.com/sirfrank96/go-server/util"
)

//assuming right handed golfer

//2 calibration images? 1 for perpendicular axes, 1 for vanishing points
//img 1: stand straddled, check to make sure heel horizontal and spine vertical are close to 90
//img 2: set alignment stick not centered, point alignment stick at target, set up with heels against alignment stick feet shoulder width or wider (check that heels are not centered in image)
// get vanishing point, intersection of vertaxis and heels axis
func VerifyDTLCalibrationImages(axesKeypoints *cv.Body25PoseKeypoints, vanishingPointKeypoints *cv.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (*util.CalibrationInfo, util.Warning) {
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

//spine angle
//line from midhip to neck
//angle between that and vertical axis
func GetSpineAngle(keypoints *cv.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
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
	return angleAtIntersect, warning
}

//feet alignment
//assume feet are left of vert axis (maybe use toes? easier to see?)
//TODO: add other edge cases for feet crossing vertaxis
func GetFeetAlignment(keypoints *cv.Body25PoseKeypoints, calibrationInfo *util.CalibrationInfo) (float64, util.Warning) {
	// TODO: IF calibrationInfo.AxesWarning or vanishing point is not nil return that
	var warning util.Warning
	currFeetLine, w := util.GetFeetLine(keypoints, calibrationInfo.FeetLine.FeetLineMethod)
	if w != nil {
		if w.GetWarningType() == util.SEVERE {
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

//waist bend
//line between rhip to neck
//line between rhip to rknee
//angle between those lines

//shoulder alignment
//find vanishing point
//any line that goes to vanishing point
//line from rshoulder to lshoulder
//find difference in angles
