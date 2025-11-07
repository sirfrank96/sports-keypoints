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
func checkIfDTLCalibrationImagesAreGood(axesKeypoints *cv.Body25PoseKeypoints, vanishingPointKeypoints *cv.Body25PoseKeypoints) (*CalibrationInfo, error) {
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

//spine angle
//line from midhip to neck
//angle between that and vertical axis
func getSpineAngle(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) float64 {
	vertAxisLine := calibrationInfo.vertAxisLine
	fmt.Printf("VertAxisLine object: %+v\n", vertAxisLine)
	vertAxisThroughMidhipLine := getLineWithSlope(keypoints.Midhip, vertAxisLine.slope)
	fmt.Printf("VertAxisThroughMidhipLine object: %+v\n", vertAxisThroughMidhipLine)
	neckPoint := convertCvKeypointToPoint(keypoints.Neck)
	fmt.Printf("NeckPoint: %+v\n", neckPoint)
	xOnVertAxis := (keypoints.Neck.Y - vertAxisThroughMidhipLine.yIntercept) / vertAxisThroughMidhipLine.slope
	pointUpVertAxisSameHeightAsNeck := &Point{xPos: xOnVertAxis, yPos: keypoints.Neck.Y}
	fmt.Printf("PointUpVertAxisSameHeightAsNeck: %+v\n", pointUpVertAxisSameHeightAsNeck)
	midhipPoint := convertCvKeypointToPoint(keypoints.Midhip)
	fmt.Printf("MidhipPoint: %+v\n", midhipPoint)
	angleAtIntersect := getAngleAtIntersection(neckPoint, midhipPoint, pointUpVertAxisSameHeightAsNeck)
	return angleAtIntersect
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

//alignents
//all relative to heel alignment (open to heels, closed to heels)
