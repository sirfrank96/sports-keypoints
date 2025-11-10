package processor

import (
	"fmt"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
)

//assuming right handed golfer

// TODO: Make how far off axes are configurable
func verifyFaceOnCalibrationImage(keypoints *cv.Body25PoseKeypoints, feetLineMethod cv.FeetLineMethod) (*CalibrationInfo, error) {
	return verifyCalibrationImageAxes(keypoints, feetLineMethod)
}

//side bend
//line from midhip to neck
//angle of intersect between that and vertical axis through midhip
func getSideBend(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) (float64, error) {
	var warning error
	if err := verifyKeypoint(keypoints.Midhip, "midhip", 0.5); err != nil {
		warning = appendError(warning, err)
	}
	if err := verifyKeypoint(keypoints.Neck, "neck", 0.5); err != nil {
		warning = appendError(warning, err)
	}
	vertAxisLine := calibrationInfo.vertAxisLine
	fmt.Printf("VertAxisLine object: %+v\n", vertAxisLine)
	vertAxisThroughMidhipLine := getLineWithSlope(convertCvKeypointToPoint(keypoints.Midhip), vertAxisLine.slope)
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
		return angleAtIntersect, warning
	} else { // left
		return float64(-1) * angleAtIntersect, warning
	}
}

//foot flares
//line from heel to big toe
//angle of intersect vert axis

//stance width
//relative to hip to neck length
//line from left heel to right heel
//ratio of that line to hip to neck length
