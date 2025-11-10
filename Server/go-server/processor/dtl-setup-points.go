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
func verifyDTLCalibrationImages(axesKeypoints *cv.Body25PoseKeypoints, vanishingPointKeypoints *cv.Body25PoseKeypoints, feetLineMethod cv.FeetLineMethod) (*CalibrationInfo, error) {
	// verify axes image
	calibrationInfo, err := verifyCalibrationImageAxes(axesKeypoints, feetLineMethod)
	if err != nil {
		return nil, err
	}
	// verify vanishing point image
	feetLineInfo, _ := getFeetLineInfo(vanishingPointKeypoints, feetLineMethod)
	if err := verifyFeetLineInfo(feetLineInfo, 0.5); err != nil {
		return nil, err
	}
	slopeDiff := math.Abs(feetLineInfo.feetLine.slope - calibrationInfo.vertAxisLine.slope)
	if slopeDiff < float64(1) { // TODO: Configure how close slope is (and how to determine how close slope is)
		return nil, fmt.Errorf("vanishing point calibration image off. feet line slope %f and vertaxis line slope %f are too close (%f). make sure feet line is off centered or make sure alignment stick is pointed at target (parallel lines converge in distance)", feetLineInfo.feetLine.slope, calibrationInfo.vertAxisLine.slope, slopeDiff)
	}
	fmt.Printf("Good vanishing point calibration. heel line slope is %f, and vertaxis line slope is %f", feetLineInfo.feetLine.slope, calibrationInfo.vertAxisLine.slope)
	intersection := getIntersection(feetLineInfo.feetLine, calibrationInfo.vertAxisLine)
	calibrationInfo.vanishingPoint = intersection.intersectPoint
	return calibrationInfo, nil
}

//spine angle
//line from midhip to neck
//angle between that and vertical axis
func getSpineAngle(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) (float64, error) {
	var warning error
	if err := verifyKeypoint(keypoints.Midhip, "midhip", 0.5); err != nil {
		warning = appendError(warning, err)
	}
	if err := verifyKeypoint(keypoints.Neck, "neck", 0.5); err != nil {
		warning = appendError(warning, err)
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
func getFeetAlignment(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) (float64, error) {
	var warning error
	currFeetLineInfo, _ := getFeetLineInfo(keypoints, calibrationInfo.feetLineInfo.feetLineMethod)
	if err := verifyFeetLineInfo(currFeetLineInfo, 0.6); err != nil {
		warning = appendError(warning, err)
	}
	fmt.Printf("FeetLineInfo: %+v, LPoint: %+v, RPoint: %+v\n", currFeetLineInfo, currFeetLineInfo.lPoint, currFeetLineInfo.rPoint)
	fmt.Printf("VanishingPoint: %+v\n", calibrationInfo.vanishingPoint)
	feetDegrees := getAngleAtIntersection(convertCvKeypointToPoint(currFeetLineInfo.lPoint), convertCvKeypointToPoint(currFeetLineInfo.rPoint), calibrationInfo.vanishingPoint)
	fmt.Printf("Angle at right point intersection: %f\n", feetDegrees)
	realParallelLine := getLine(calibrationInfo.vanishingPoint, convertCvKeypointToPoint(currFeetLineInfo.rPoint)) // line that converges at vanishing point
	fmt.Printf("Real parallel is %+v\n", realParallelLine)
	fmt.Printf("Current feet line is %+v\n", currFeetLineInfo.feetLine)
	//determine if closed or open
	if currFeetLineInfo.feetLine.slope == realParallelLine.slope { // neutral
		return 0, warning
	} else if currFeetLineInfo.feetLine.slope > realParallelLine.slope && currFeetLineInfo.feetLine.slope <= 0 { // closed (will return positive number)
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
