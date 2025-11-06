package processor

import (
	//"fmt"
	"math"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
)

//assuming right handed golfer

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

//util
func getSlope(point1 *cv.Keypoint, point2 *cv.Keypoint) float64 {
	rise := point2.Y - point1.Y
	run := point2.X - point1.X

	// TODO: handle 0 on denominator

	return rise / run
}

func getSlopeRecipricol(point1 *cv.Keypoint, point2 *cv.Keypoint) float64 {
	rise := point2.Y - point1.Y
	run := point2.X - point1.X

	// TODO: handle 0 denominator

	return float64(-1) * (run / rise)
}

func getYIntercept(point *Point, slope float64) float64 {
	return point.yPos - (slope * point.xPos)
}

func getMidpoint(point1 *cv.Keypoint, point2 *cv.Keypoint) *Point {
	xMid := (point1.X + point2.X) / float64(2)
	yMid := (point1.Y + point2.Y) / float64(2)
	return &Point{xPos: xMid, yPos: yMid}
}

// point1 will be the pointOnLine
func getLine(point1 *cv.Keypoint, point2 *cv.Keypoint) *Line {
	slope := getSlope(point1, point2)
	return getLineWithSlope(point1, slope)
}

func getLineWithSlope(point1 *cv.Keypoint, slope float64) *Line {
	yIntercept := getYIntercept(convertCvKeypointToPoint(point1), slope)
	return &Line{slope: slope, yIntercept: yIntercept, pointOnLine: convertCvKeypointToPoint(point1)}
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
