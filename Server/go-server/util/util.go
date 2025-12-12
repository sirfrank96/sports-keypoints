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

// TODO: Do i need lines? or just use angles somehow
type Line struct {
	Slope       float64 `bson:"slope,omitempty"`
	YIntercept  float64 `bson:"y_intercept,omitempty"`
	PointOnLine Point   `bson:"point_on_line,omitempty"`
	// TODO: Add infinite slope/vertical line parameter
}

// Info about intersection between 2 lines
type Intersection struct {
	Line1            Line    `bson:"line1,omitempty"`
	Line2            Line    `bson:"line2,omitempty"`
	IntersectPoint   Point   `bson:"intersect_point,omitempty"`
	AngleAtIntersect float64 `bson:"angle_at_intersect,omitempty"`
}

type Projection struct {
	IntersectPoint Point
	ProjectionLine Line
	OriginalPoint  Point
	OriginalLine   Line
}

// Find the length between point1 and point2
func GetLengthBetweenTwoPoints(point1 *Point, point2 *Point) float64 {
	term1 := math.Pow(point2.XPos-point1.XPos, 2)
	term2 := math.Pow(point2.YPos-point1.YPos, 2)
	return math.Sqrt(term1 + term2)
}

// Get the slope of the line that passes through point1 and point2
func GetSlope(point1 *Point, point2 *Point) float64 {
	rise := point2.YPos - point1.YPos
	run := point2.XPos - point1.XPos

	// TODO: handle 0 on denominator

	return rise / run
}

// Find the recipricol of the given fraction
func GetRecipricol(fraction float64) float64 {
	return float64(-1) * (1 / fraction)
}

// Find the recipricol of the slope of the line that passes through point1 and point2
func GetSlopeRecipricol(point1 *Point, point2 *Point) float64 {
	rise := point2.YPos - point1.YPos
	run := point2.XPos - point1.XPos

	// TODO: handle 0 denominator

	return float64(-1) * (run / rise)
}

// Find the y intercept of the line with the given point and slope
func GetYIntercept(point *Point, slope float64) float64 {
	return point.YPos - (slope * point.XPos)
}

// Find the midpoint between point1 and point2
func GetMidpoint(point1 *Point, point2 *Point) *Point {
	xMid := (point1.XPos + point2.XPos) / float64(2)
	yMid := (point1.YPos + point2.YPos) / float64(2)
	return &Point{XPos: xMid, YPos: yMid}
}

// Return the line intersects point1 and point2
// point1 will be the pointOnLine in the Line struct
func GetLine(point1 *Point, point2 *Point) *Line {
	slope := GetSlope(point1, point2)
	return GetLineWithSlope(point1, slope)
}

// Given a point and a slope, return the line
func GetLineWithSlope(point1 *Point, slope float64) *Line {
	yIntercept := GetYIntercept(point1, slope)
	return &Line{Slope: slope, YIntercept: yIntercept, PointOnLine: *point1}
}

// Find the point on the line with the given x coordinate
func GetPointOnLineWithX(x float64, line *Line) *Point {
	yOnLine := (line.Slope * x) + line.YIntercept
	return &Point{XPos: x, YPos: yOnLine}
}

// Find the point on the line with the given y coordinate
func GetPointOnLineWithY(y float64, line *Line) *Point {
	xOnLine := (y - line.YIntercept) / line.Slope
	return &Point{XPos: xOnLine, YPos: y}
}

// Given three points, find the angle at the intersection between the line from point1 to intersectPoint and point2 and intersectPoint
// Uses law of cosines
func GetAngleAtIntersection(point1 *Point, intersectPoint *Point, point2 *Point) float64 {
	lenLineOppIntersect := GetLengthBetweenTwoPoints(point1, point2)
	lenLineBetweenIntersectAnd1 := GetLengthBetweenTwoPoints(point1, intersectPoint)
	lenLineBetweenIntersectAnd2 := GetLengthBetweenTwoPoints(point2, intersectPoint)
	numerator := math.Pow(lenLineBetweenIntersectAnd1, 2) + math.Pow(lenLineBetweenIntersectAnd2, 2) - math.Pow(lenLineOppIntersect, 2)
	denominator := 2 * lenLineBetweenIntersectAnd1 * lenLineBetweenIntersectAnd2
	radAngle := math.Acos(numerator / denominator)
	return ConvertRadToDegrees(radAngle)
}

// Given two lines, find the intersection of the two lines
func GetIntersection(line1 *Line, line2 *Line) *Intersection {
	xIntersect := (line2.YIntercept - line1.YIntercept) / (line1.Slope - line2.Slope)
	yIntersect := line1.Slope*xIntersect + line1.YIntercept
	intersectPoint := Point{XPos: xIntersect, YPos: yIntersect}
	angleAtIntersect := GetAngleAtIntersection(&line1.PointOnLine, &intersectPoint, &line2.PointOnLine)
	return &Intersection{Line1: *line1, Line2: *line2, IntersectPoint: intersectPoint, AngleAtIntersect: angleAtIntersect}
}

// Given a line and a point, find the projection onto line from the point
// ie. the line perpendicular to the line that runs through the point
func GetProjectionOntoLine(line *Line, point *Point) *Projection {
	slopeOfProjection := GetRecipricol(line.Slope)
	projectionLine := GetLineWithSlope(point, slopeOfProjection)
	intersection := GetIntersection(line, projectionLine)
	return &Projection{IntersectPoint: intersection.IntersectPoint, ProjectionLine: *projectionLine, OriginalPoint: *point, OriginalLine: *line}
}

// TODO: keep from -90 to 90 or keep from 0-180?
func ConvertSlopeToDegrees(slope float64) float64 {
	rad := math.Atan(slope)
	deg := ConvertRadToDegrees(rad)
	if deg > 180.0 {
		return deg - float64(180)
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

// Converts *skp.Keypoint to *Point
func ConvertKeypointToPoint(cvKeypoint *skp.Keypoint) *Point {
	if cvKeypoint == nil {
		return nil
	}
	return &Point{XPos: cvKeypoint.X, YPos: cvKeypoint.Y}
}

// Converts *Point to *skp.Keypoint
func ConvertPointToKeypoint(point *Point) *skp.Keypoint {
	if point == nil {
		return nil
	}
	return &skp.Keypoint{X: point.XPos, Y: point.YPos}
}
