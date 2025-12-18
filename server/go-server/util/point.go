package util

import (
	"math"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

type Point struct {
	XPos float64 `bson:"x_pos,omitempty"`
	YPos float64 `bson:"y_pos,omitempty"`
}

// Find the length between point1 and point2
func GetLengthBetweenTwoPoints(point1 *Point, point2 *Point) float64 {
	term1 := math.Pow(point2.XPos-point1.XPos, 2)
	term2 := math.Pow(point2.YPos-point1.YPos, 2)
	return math.Sqrt(term1 + term2)
}

// Find the midpoint between point1 and point2
func GetMidpoint(point1 *Point, point2 *Point) *Point {
	xMid := (point1.XPos + point2.XPos) / float64(2)
	yMid := (point1.YPos + point2.YPos) / float64(2)
	return &Point{XPos: xMid, YPos: yMid}
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
