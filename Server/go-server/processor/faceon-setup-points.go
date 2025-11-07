package processor

import (
	"fmt"
	"math"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
)

//assuming right handed golfer

// TODO: Make how far off axes are configurable
func checkIfFaceOnCalibrationImageIsGood(keypoints *cv.Body25PoseKeypoints) (*CalibrationInfo, error) {
	horAxisLine := getHorizontalAxisLine(keypoints.LHeel, keypoints.RHeel)
	vertAxisLine := getVerticalAxisLine(keypoints.Midhip, keypoints.Neck)
	horDeg := convertSlopeToDegrees(horAxisLine.slope)
	vertDeg := convertSlopeToDegrees(vertAxisLine.slope)
	diff := math.Abs(vertDeg) + math.Abs(horDeg) - 90
	if math.Abs(diff) > 10 { // make this 5 or less after better test images
		return nil, fmt.Errorf("axes calibration image off. horizontal axis between heels is %f degrees. vertical axis between midhip and neck is %f degrees. difference of %f degrees is too large. please adjust camera, stance, or posture. recommend using alignment sticks to help calibration", horDeg, vertDeg, diff)
	}
	fmt.Printf("Good axes calibration. Horizontal axis between heels is %f degrees. vertical axis between midhip and neck is %f degrees\n", horDeg, vertDeg)
	return &CalibrationInfo{horAxisLine: horAxisLine, vertAxisLine: vertAxisLine}, nil
}

//side bend
//line from midhip to neck
//angle of intersect between that and vertical axis through midhip
func getSideBend(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) float64 {
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
	// determine if left or right side bend
	if keypoints.Neck.X < keypoints.Midhip.X { // right
		return angleAtIntersect
	} else { // left
		return float64(-1) * angleAtIntersect
	}
}

//foot flares
//line from heel to big toe
//angle of intersect vert axis

//stance width
//relative to hip to neck length
//line from left heel to right heel
//ratio of that line to hip to neck length
