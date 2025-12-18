package util

type Line struct {
	Slope       float64 `bson:"slope,omitempty"`
	YIntercept  float64 `bson:"y_intercept,omitempty"`
	PointOnLine Point   `bson:"point_on_line,omitempty"`
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
