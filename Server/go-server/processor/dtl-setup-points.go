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

//TODO: pull out axes calibration into common func
func verifyDTLCalibrationImages(axesKeypoints *cv.Body25PoseKeypoints, vanishingPointKeypoints *cv.Body25PoseKeypoints) (*CalibrationInfo, error) {
	horAxisLine := getHorizontalAxisLine(axesKeypoints.LHeel, axesKeypoints.RHeel)
	vertAxisLine := getVerticalAxisLine(axesKeypoints.Midhip, axesKeypoints.Neck)
	horDeg := convertSlopeToDegrees(horAxisLine.slope)
	vertDeg := convertSlopeToDegrees(vertAxisLine.slope)
	axesDiff := math.Abs(vertDeg) + math.Abs(horDeg) - 90
	if math.Abs(axesDiff) > 10 { // make this 5 or less after better test images
		return nil, fmt.Errorf("axes calibration image off. horizontal axis between heels is %f degrees. vertical axis between midhip and neck is %f degrees. difference of %f degrees is too large. please adjust camera, stance, or posture. recommend using alignment sticks to help calibration", horDeg, vertDeg, axesDiff)
	}
	fmt.Printf("Good axes calibration. Horizontal axis between heels is %f degrees. vertical axis between midhip and neck is %f degrees\n", horDeg, vertDeg)

	//TODO: Check that if heels are left of vertline slope is pos?or neg? and vice versa for heels right of vertline
	heelLine := getLine(vanishingPointKeypoints.RHeel, vanishingPointKeypoints.LHeel)
	slopeDiff := math.Abs(heelLine.slope - vertAxisLine.slope)
	if slopeDiff < float64(1) {
		return nil, fmt.Errorf("vanishing point calibration image off. heel line slope %f and vertaxis line slope %f are too close (%f). make sure heel line is off centered or make sure alignment stick is pointed at target (parallel lines converge in distance)", heelLine.slope, vertAxisLine.slope, slopeDiff)
	}
	fmt.Printf("Good vanishing point calibration. heel line slope is %f, and vertaxis line slope is %f", heelLine.slope, vertAxisLine.slope)
	intersection := getIntersection(heelLine, vertAxisLine)

	return &CalibrationInfo{horAxisLine: horAxisLine, vertAxisLine: vertAxisLine, vanishingPoint: intersection.intersectPoint}, nil
}

func verifyDTLImage(keypoints *cv.Body25PoseKeypoints) error {
	return nil
}

//spine angle
//line from midhip to neck
//angle between that and vertical axis
func getSpineAngle(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) (float64, error) {
	// TODO: CHECK CONFIDENCE LEVELS
	vertAxisLine := calibrationInfo.vertAxisLine
	vertAxisThroughMidhipLine := getLineWithSlope(keypoints.Midhip, vertAxisLine.slope)
	neckPoint := convertCvKeypointToPoint(keypoints.Neck)
	xOnVertAxis := (keypoints.Neck.Y - vertAxisThroughMidhipLine.yIntercept) / vertAxisThroughMidhipLine.slope
	pointUpVertAxisSameHeightAsNeck := &Point{xPos: xOnVertAxis, yPos: keypoints.Neck.Y}
	midhipPoint := convertCvKeypointToPoint(keypoints.Midhip)
	angleAtIntersect := getAngleAtIntersection(neckPoint, midhipPoint, pointUpVertAxisSameHeightAsNeck)
	return angleAtIntersect, nil
}

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

//feet alignment
//assume feet are left of vert axis (maybe use toes? easier to see?)
//TODO: add other edge cases for feet crossing vertaxis
//TODO: Use toes??? easier to see (If confidence is < 0.5, use toes)
/*func getFeetAlignment(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) float64 {
	fmt.Printf("LHeel: %+v\n, RHeel: %+v\n, vanishingPoint: %+v\n", keypoints.LHeel, keypoints.RHeel, calibrationInfo.vanishingPoint)
	feetDegrees := getAngleAtIntersection(convertCvKeypointToPoint(keypoints.LHeel), convertCvKeypointToPoint(keypoints.RHeel), calibrationInfo.vanishingPoint)
	fmt.Printf("Angle at right heel intersection: %f\n", feetDegrees)
	realParallelLine := getLine(convertPointToCvKeypoint(calibrationInfo.vanishingPoint), keypoints.RHeel) // line that converges at vanishing point
	fmt.Printf("Real parallel is %+v\n", realParallelLine)
	heelLine := getLine(keypoints.RHeel, keypoints.LHeel)
	fmt.Printf("Current heel line is %+v\n", heelLine)
	//determine if closed or open
	if heelLine.slope == realParallelLine.slope { // neutral
		return 0
	} else if heelLine.slope > realParallelLine.slope && heelLine.slope <= 0 { // closed (will return positive number)
		return feetDegrees
	} else { // open (will return negative number)
		return float64(-1) * feetDegrees
	}
}*/

// test toes
func getFeetAlignment(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) (float64, error) {
	// TODO: CHECK CONFIDENCE LEVELS
	fmt.Printf("LBigToe: %+v\n, RBigToe: %+v\n, vanishingPoint: %+v\n", keypoints.LBigToe, keypoints.RBigToe, calibrationInfo.vanishingPoint)
	feetDegrees := getAngleAtIntersection(convertCvKeypointToPoint(keypoints.LBigToe), convertCvKeypointToPoint(keypoints.RBigToe), calibrationInfo.vanishingPoint)
	fmt.Printf("Angle at right toe intersection: %f\n", feetDegrees)
	realParallelLine := getLine(convertPointToCvKeypoint(calibrationInfo.vanishingPoint), keypoints.RBigToe) // line that converges at vanishing point
	fmt.Printf("Real parallel is %+v\n", realParallelLine)
	toeLine := getLine(keypoints.RBigToe, keypoints.LBigToe)
	fmt.Printf("Current toe line is %+v\n", toeLine)
	//determine if closed or open
	if toeLine.slope == realParallelLine.slope { // neutral
		return 0, nil
	} else if toeLine.slope > realParallelLine.slope && toeLine.slope <= 0 { // closed (will return positive number)
		return feetDegrees, nil
	} else { // open (will return negative number)
		return float64(-1) * feetDegrees, nil
	}
}

//alignents
//all relative to heel alignment (open to heels, closed to heels)
