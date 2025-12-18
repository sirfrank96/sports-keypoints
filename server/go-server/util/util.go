package util

import (
	"math"
)

// assuming right handed golfer

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

// Find the y intercept of the line with the given point and slope
func GetYIntercept(point *Point, slope float64) float64 {
	return point.YPos - (slope * point.XPos)
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

func ConvertDegreesToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func GetVector(point *Point, origin *Point) *Point {
	vect := &Point{
		XPos: point.XPos - origin.XPos,
		YPos: point.YPos - origin.YPos,
	}
	return vect
}

func GetDotProduct(vect1 *Point, vect2 *Point) float64 {
	return (vect1.XPos * vect2.XPos) + (vect1.YPos * vect2.YPos)
}

func GetCrossProduct(vect1 *Point, vect2 *Point) float64 {
	return (vect1.XPos * vect2.YPos) - (vect1.YPos * vect2.XPos)
}

func GetSignedAngleOfRotation(vect1 *Point, vect2 *Point) float64 {
	det := GetCrossProduct(vect1, vect2)
	dot := GetDotProduct(vect1, vect2)
	rad := math.Atan2(det, dot)
	return ConvertRadToDegrees(rad)
}
