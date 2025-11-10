package processor

import (
	"fmt"
	"math"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
)

//assuming right handed golfer

//TODO: Be able to pass around Point or cv.Keypoint by using an interface

type Point struct {
	xPos float64
	yPos float64
}

type Line struct {
	slope       float64
	yIntercept  float64
	pointOnLine *Point
}

// Info about intersection between 2 lines
type Intersection struct {
	line1            *Line
	line2            *Line
	intersectPoint   *Point
	angleAtIntersect float64
}

type FeetLineInfo struct {
	feetLineMethod cv.FeetLineMethod
	lPoint         *cv.Keypoint
	rPoint         *cv.Keypoint
	feetLine       *Line
}

type CalibrationInfo struct {
	feetLineInfo   *FeetLineInfo
	horAxisLine    *Line
	vertAxisLine   *Line
	vanishingPoint *Point
}

//util
func getLengthBetweenTwoPoints(point1 *Point, point2 *Point) float64 {
	term1 := math.Pow(point2.xPos-point1.xPos, 2)
	term2 := math.Pow(point2.yPos-point1.yPos, 2)
	return math.Sqrt(term1 + term2)
}

func getSlope(point1 *Point, point2 *Point) float64 {
	rise := point2.yPos - point1.yPos
	run := point2.xPos - point1.xPos

	// TODO: handle 0 on denominator

	return rise / run
}

func getSlopeRecipricol(point1 *Point, point2 *Point) float64 {
	rise := point2.yPos - point1.yPos
	run := point2.xPos - point1.xPos

	// TODO: handle 0 denominator

	return float64(-1) * (run / rise)
}

func getYIntercept(point *Point, slope float64) float64 {
	return point.yPos - (slope * point.xPos)
}

func getMidpoint(point1 *Point, point2 Point) *Point {
	xMid := (point1.xPos + point2.xPos) / float64(2)
	yMid := (point1.yPos + point2.yPos) / float64(2)
	return &Point{xPos: xMid, yPos: yMid}
}

// point1 will be the pointOnLine
func getLine(point1 *Point, point2 *Point) *Line {
	slope := getSlope(point1, point2)
	return getLineWithSlope(point1, slope)
}

func getLineWithSlope(point1 *Point, slope float64) *Line {
	yIntercept := getYIntercept(point1, slope)
	return &Line{slope: slope, yIntercept: yIntercept, pointOnLine: point1}
}

// law of cosines
func getAngleAtIntersection(point1 *Point, intersectPoint *Point, point2 *Point) float64 {
	lenLineOppIntersect := getLengthBetweenTwoPoints(point1, point2)
	lenLineBetweenIntersectAnd1 := getLengthBetweenTwoPoints(point1, intersectPoint)
	lenLineBetweenIntersectAnd2 := getLengthBetweenTwoPoints(point2, intersectPoint)
	numerator := math.Pow(lenLineBetweenIntersectAnd1, 2) + math.Pow(lenLineBetweenIntersectAnd2, 2) - math.Pow(lenLineOppIntersect, 2)
	denominator := 2 * lenLineBetweenIntersectAnd1 * lenLineBetweenIntersectAnd2
	radAngle := math.Acos(numerator / denominator)
	return convertRadToDegrees(radAngle)
}

func getIntersection(line1 *Line, line2 *Line) *Intersection {
	xIntersect := (line2.yIntercept - line1.yIntercept) / (line1.slope - line2.slope)
	yIntersect := line1.slope*xIntersect + line1.yIntercept
	intersectPoint := &Point{xPos: xIntersect, yPos: yIntersect}
	angleAtIntersect := getAngleAtIntersection(line1.pointOnLine, intersectPoint, line2.pointOnLine)
	return &Intersection{line1: line1, line2: line2, intersectPoint: intersectPoint, angleAtIntersect: angleAtIntersect}
}

// keep from 0-180
func convertSlopeToDegrees(slope float64) float64 {
	rad := math.Atan(slope)
	deg := convertRadToDegrees(rad)
	if deg > 180 {
		return deg - 180
	} else {
		return deg
	}
}

func convertRadToDegrees(rad float64) float64 {
	return rad * 180 / math.Pi
}

func getDegreesOfLineAlwaysPositive(deg float64) float64 {
	if deg < 0 {
		return deg + float64(180)
	} else if deg > 180 {
		return deg - float64(180)
	} else {
		return deg
	}
}

func convertCvKeypointToPoint(cvKeypoint *cv.Keypoint) *Point {
	return &Point{xPos: cvKeypoint.X, yPos: cvKeypoint.Y}
}

func convertPointToCvKeypoint(point *Point) *cv.Keypoint {
	return &cv.Keypoint{X: point.xPos, Y: point.yPos}
}

func verifyKeypoint(keypoint *cv.Keypoint, keypointName string, threshold float64) error {
	if keypoint.Confidence < threshold {
		return fmt.Errorf("uncertain where %s is, confidence is %f. please make sure %s is visible in image", keypointName, keypoint.Confidence, keypointName)
	}
	return nil
}

func verifyFeetLineInfo(feetLineInfo *FeetLineInfo, threshold float64) error {
	var err error
	if e := verifyKeypoint(feetLineInfo.lPoint, "left point", threshold); e != nil {
		err = appendError(err, fmt.Errorf("%w, please set a different FeetLineMethod", e))
	}
	if e := verifyKeypoint(feetLineInfo.rPoint, "right point", threshold); e != nil {
		err = appendError(err, fmt.Errorf("%w, please set a different FeetLineMethod", e))
	}
	return err
}

func getFeetLineInfo(keypoints *cv.Body25PoseKeypoints, feetLineMethod cv.FeetLineMethod) (*FeetLineInfo, error) {
	feetLineInfo := &FeetLineInfo{feetLineMethod: feetLineMethod}
	if feetLineMethod == cv.FeetLineMethod_USE_TOE_LINE {
		feetLineInfo.lPoint = keypoints.LBigToe
		feetLineInfo.rPoint = keypoints.RBigToe
	} else { // default is USE_HEEL_LINE
		feetLineInfo.lPoint = keypoints.LHeel
		feetLineInfo.rPoint = keypoints.RHeel
	}
	feetLineInfo.feetLine = getLine(convertCvKeypointToPoint(feetLineInfo.lPoint), convertCvKeypointToPoint(feetLineInfo.rPoint))
	return feetLineInfo, nil
}

// TODO: Configure confidence level
func verifyCalibrationImageAxes(keypoints *cv.Body25PoseKeypoints, feetLineMethod cv.FeetLineMethod) (*CalibrationInfo, error) {
	// Get horizontal axis
	feetLineInfo, _ := getFeetLineInfo(keypoints, feetLineMethod)
	if err := verifyFeetLineInfo(feetLineInfo, 0.5); err != nil {
		return nil, err
	}
	// Get vertical axis
	if err := verifyKeypoint(keypoints.Midhip, "midhip", 0.5); err != nil {
		return nil, err
	}
	if err := verifyKeypoint(keypoints.Neck, "neck", 0.5); err != nil {
		return nil, err
	}
	horAxisLine := feetLineInfo.feetLine
	vertAxisLine := getLine(convertCvKeypointToPoint(keypoints.Midhip), convertCvKeypointToPoint(keypoints.Neck))
	// Check if angle between axes is around 90 degrees
	horDeg := convertSlopeToDegrees(horAxisLine.slope)
	vertDeg := convertSlopeToDegrees(vertAxisLine.slope)
	diff := math.Abs(vertDeg) + math.Abs(horDeg) - 90
	if math.Abs(diff) > 10 { // make this 5 or less after better test images
		return nil, fmt.Errorf("axes calibration image off. horizontal axis between heels is %f degrees. vertical axis between midhip and neck is %f degrees. difference of %f degrees is too large. please adjust camera, stance, or posture. recommend using alignment sticks to help calibration", horDeg, vertDeg, diff)
	}
	fmt.Printf("Good axes calibration. Horizontal axis between heels is %f degrees. vertical axis between midhip and neck is %f degrees\n", horDeg, vertDeg)
	return &CalibrationInfo{feetLineInfo: feetLineInfo, horAxisLine: horAxisLine, vertAxisLine: vertAxisLine}, nil
}

func appendError(err1 error, err2 error) error {
	if err1 != nil && err2 != nil {
		return fmt.Errorf("%w, %w", err1, err2)
	} else if err1 == nil {
		return err2
	} else if err2 == nil {
		return err1
	} else {
		return nil
	}
}
