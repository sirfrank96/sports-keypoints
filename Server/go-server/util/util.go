package util

import (
	"math"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

//assuming right handed golfer

//TODO: Be able to pass around Point or cv.Keypoint by using an interface

type Point struct {
	XPos float64 `bson:"x_pos,omitempty"`
	YPos float64 `bson:"y_pos,omitempty"`
}

type Line struct {
	Slope       float64 `bson:"slope,omitempty"`
	YIntercept  float64 `bson:"y_intercept,omitempty"`
	PointOnLine Point   `bson:"point_on_line,omitempty"`
}

// Info about intersection between 2 lines
type Intersection struct {
	Line1            Line    `bson:"line1,omitempty"`
	Line2            Line    `bson:"line2,omitempty"`
	IntersectPoint   Point   `bson:"intersect_point,omitempty"`
	AngleAtIntersect float64 `bson:"angle_at_intersect,omitempty"`
}

//util
func GetLengthBetweenTwoPoints(point1 *Point, point2 *Point) float64 {
	term1 := math.Pow(point2.XPos-point1.XPos, 2)
	term2 := math.Pow(point2.YPos-point1.YPos, 2)
	return math.Sqrt(term1 + term2)
}

func GetSlope(point1 *Point, point2 *Point) float64 {
	rise := point2.YPos - point1.YPos
	run := point2.XPos - point1.XPos

	// TODO: handle 0 on denominator

	return rise / run
}

func GetSlopeRecipricol(point1 *Point, point2 *Point) float64 {
	rise := point2.YPos - point1.YPos
	run := point2.XPos - point1.XPos

	// TODO: handle 0 denominator

	return float64(-1) * (run / rise)
}

func GetYIntercept(point *Point, slope float64) float64 {
	return point.YPos - (slope * point.XPos)
}

func GetMidpoint(point1 *Point, point2 *Point) *Point {
	xMid := (point1.XPos + point2.XPos) / float64(2)
	yMid := (point1.YPos + point2.YPos) / float64(2)
	return &Point{XPos: xMid, YPos: yMid}
}

// point1 will be the pointOnLine
func GetLine(point1 *Point, point2 *Point) *Line {
	slope := GetSlope(point1, point2)
	return GetLineWithSlope(point1, slope)
}

func GetLineWithSlope(point1 *Point, slope float64) *Line {
	yIntercept := GetYIntercept(point1, slope)
	return &Line{Slope: slope, YIntercept: yIntercept, PointOnLine: *point1}
}

// law of cosines
func GetAngleAtIntersection(point1 *Point, intersectPoint *Point, point2 *Point) float64 {
	lenLineOppIntersect := GetLengthBetweenTwoPoints(point1, point2)
	lenLineBetweenIntersectAnd1 := GetLengthBetweenTwoPoints(point1, intersectPoint)
	lenLineBetweenIntersectAnd2 := GetLengthBetweenTwoPoints(point2, intersectPoint)
	numerator := math.Pow(lenLineBetweenIntersectAnd1, 2) + math.Pow(lenLineBetweenIntersectAnd2, 2) - math.Pow(lenLineOppIntersect, 2)
	denominator := 2 * lenLineBetweenIntersectAnd1 * lenLineBetweenIntersectAnd2
	radAngle := math.Acos(numerator / denominator)
	return ConvertRadToDegrees(radAngle)
}

func GetIntersection(line1 *Line, line2 *Line) *Intersection {
	xIntersect := (line2.YIntercept - line1.YIntercept) / (line1.Slope - line2.Slope)
	yIntersect := line1.Slope*xIntersect + line1.YIntercept
	intersectPoint := Point{XPos: xIntersect, YPos: yIntersect}
	angleAtIntersect := GetAngleAtIntersection(&line1.PointOnLine, &intersectPoint, &line2.PointOnLine)
	return &Intersection{Line1: *line1, Line2: *line2, IntersectPoint: intersectPoint, AngleAtIntersect: angleAtIntersect}
}

// keep from 0-180
func ConvertSlopeToDegrees(slope float64) float64 {
	rad := math.Atan(slope)
	deg := ConvertRadToDegrees(rad)
	if deg > 180 {
		return deg - 180
	} else {
		return deg
	}
}

func ConvertRadToDegrees(rad float64) float64 {
	return rad * 180 / math.Pi
}

func GetDegreesOfLineAlwaysPositive(deg float64) float64 {
	if deg < 0 {
		return deg + float64(180)
	} else if deg > 180 {
		return deg - float64(180)
	} else {
		return deg
	}
}

func ConvertCvKeypointToPoint(cvKeypoint *skp.Keypoint) *Point {
	return &Point{XPos: cvKeypoint.X, YPos: cvKeypoint.Y}
}

func ConvertPointToCvKeypoint(point *Point) *skp.Keypoint {
	return &skp.Keypoint{X: point.XPos, Y: point.YPos}
}
