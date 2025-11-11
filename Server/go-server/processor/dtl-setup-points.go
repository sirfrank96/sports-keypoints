package processor

import (
	"fmt"
	"math"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
)

//assuming right handed golfer

//2 calibration images? 1 for perpendicular axes, 1 for vanishing points
//img 1: stand straddled, check to make sure heel horizontal and spine vertical are close to 90
//img 2: set alignment stick not centered, point alignment stick at target, set up with heels against alignment stick feet shoulder width or wider (check that heels are not centered in image)
// get vanishing point, intersection of vertaxis and heels axis
func VerifyDTLCalibrationImages(axesKeypoints *cv.Body25PoseKeypoints, vanishingPointKeypoints *cv.Body25PoseKeypoints, feetLineMethod cv.FeetLineMethod) (*CalibrationInfo, warning) {
	// verify axes image
	calibrationInfo, warning := VerifyCalibrationImageAxes(axesKeypoints, feetLineMethod)
	if warning != nil {
		return nil, warning
	}
	// verify vanishing point image
	feetLine, warning := GetFeetLine(vanishingPointKeypoints, feetLineMethod)
	if warning != nil {
		return nil, warning
	}
	slopeDiff := math.Abs(feetLine.line.slope - calibrationInfo.vertAxisLine.slope)
	if slopeDiff < float64(1) { // TODO: Configure how close slope is (and how to determine how close slope is)
		return nil, Warning{
			warningType: SEVERE,
			message:     fmt.Sprintf("vanishing point calibration image off. feet line slope %f and vertaxis line slope %f are too close (%f). make sure feet line is off centered or make sure alignment stick is pointed at target (parallel lines converge in distance)", feetLine.line.slope, calibrationInfo.vertAxisLine.slope, slopeDiff),
		}
	}
	fmt.Printf("Good vanishing point calibration. heel line slope is %f, and vertaxis line slope is %f", feetLine.line.slope, calibrationInfo.vertAxisLine.slope)
	intersection := getIntersection(feetLine.line, calibrationInfo.vertAxisLine)
	calibrationInfo.vanishingPoint = intersection.intersectPoint
	return calibrationInfo, nil
}

//spine angle
//line from midhip to neck
//angle between that and vertical axis
func GetSpineAngle(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) (float64, warning) {
	var warning warning
	if w := verifyKeypoint(keypoints.Midhip, "midhip", 0.5); w != nil {
		if w.WarningType() == SEVERE {
			return 0, w
		}
		warning = appendMinorWarnings(warning, w)
	}
	if w := verifyKeypoint(keypoints.Neck, "neck", 0.5); w != nil {
		if w.WarningType() == SEVERE {
			return 0, w
		}
		warning = appendMinorWarnings(warning, w)
	}
	vertAxisLine := calibrationInfo.vertAxisLine
	vertAxisThroughMidhipLine := getLineWithSlope(convertCvKeypointToPoint(keypoints.Midhip), vertAxisLine.slope)
	neckPoint := convertCvKeypointToPoint(keypoints.Neck)
	xOnVertAxis := (keypoints.Neck.Y - vertAxisThroughMidhipLine.yIntercept) / vertAxisThroughMidhipLine.slope
	pointUpVertAxisSameHeightAsNeck := &Point{xPos: xOnVertAxis, yPos: keypoints.Neck.Y}
	midhipPoint := convertCvKeypointToPoint(keypoints.Midhip)
	angleAtIntersect := getAngleAtIntersection(neckPoint, midhipPoint, pointUpVertAxisSameHeightAsNeck)
	return angleAtIntersect, warning
}

//feet alignment
//assume feet are left of vert axis (maybe use toes? easier to see?)
//TODO: add other edge cases for feet crossing vertaxis
func GetFeetAlignment(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) (float64, warning) {
	var warning warning
	currFeetLine, w := GetFeetLine(keypoints, calibrationInfo.feetLine.feetLineMethod)
	if w != nil {
		if w.WarningType() == SEVERE {
			return 0, w
		}
		warning = appendMinorWarnings(warning, w)
	}
	fmt.Printf("FeetLineInfo: %+v, LPoint: %+v, RPoint: %+v\n", currFeetLine, currFeetLine.lPoint, currFeetLine.rPoint)
	fmt.Printf("VanishingPoint: %+v\n", calibrationInfo.vanishingPoint)
	feetDegrees := getAngleAtIntersection(currFeetLine.lPoint, currFeetLine.rPoint, calibrationInfo.vanishingPoint)
	fmt.Printf("Angle at right point intersection: %f\n", feetDegrees)
	realParallelLine := getLine(calibrationInfo.vanishingPoint, currFeetLine.rPoint) // line that converges at vanishing point
	fmt.Printf("Real parallel is %+v\n", realParallelLine)
	fmt.Printf("Current feet line is %+v\n", currFeetLine.line)
	//determine if closed or open
	if currFeetLine.line.slope == realParallelLine.slope { // neutral
		return 0, warning
	} else if currFeetLine.line.slope > realParallelLine.slope && currFeetLine.line.slope <= 0 { // closed (will return positive number)
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
