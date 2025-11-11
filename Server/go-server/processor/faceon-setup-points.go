package processor

import (
	"fmt"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
)

//assuming right handed golfer

// TODO: Make how far off axes are configurable
func VerifyFaceOnCalibrationImage(keypoints *cv.Body25PoseKeypoints, feetLineMethod cv.FeetLineMethod) (*CalibrationInfo, warning) {
	return VerifyCalibrationImageAxes(keypoints, feetLineMethod)
}

//side bend
//line from midhip to neck
//angle of intersect between that and vertical axis through midhip
func GetSideBend(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) (float64, warning) {
	var warning warning
	if w := verifyKeypoint(keypoints.Midhip, "midhip", 0.5); w != nil {
		if w.WarningType() == SEVERE {
			return 0, w
		}
		warning = appendMinorWarnings(warning, w)
	}
	if w := verifyKeypoint(keypoints.Neck, "neck", 0.5); w != nil {
		if w.WarningType() == SEVERE {
			return 0, w
		}
		warning = appendMinorWarnings(warning, w)
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
//angle of intersect vert axis through midpoint of heels
func getFootFlare(heel *cv.Keypoint, toe *cv.Keypoint, calibrationInfo *CalibrationInfo, midpoint *Point) float64 {
	vertAxisThroughMidpoint := getLineWithSlope(midpoint, calibrationInfo.vertAxisLine.slope)
	toeToHeelLine := getLine(convertCvKeypointToPoint(toe), convertCvKeypointToPoint(heel))
	intersection := getIntersection(toeToHeelLine, vertAxisThroughMidpoint)
	if intersection.intersectPoint.yPos > toe.Y { // internal foot
		return float64(-1) * intersection.angleAtIntersect
	} else { // external foot
		return intersection.angleAtIntersect
	}
}

func GetLeftFootFlare(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) (float64, warning) {
	var warning warning
	if w := verifyKeypoint(keypoints.LHeel, "left heel", 0.5); w != nil {
		if w.WarningType() == SEVERE {
			return 0, w
		}
		warning = appendMinorWarnings(warning, w)
	}
	if w := verifyKeypoint(keypoints.RHeel, "right heel", 0.5); w != nil {
		if w.WarningType() == SEVERE {
			return 0, w
		}
		warning = appendMinorWarnings(warning, w)
	}
	if w := verifyKeypoint(keypoints.LBigToe, "left big toe", 0.5); w != nil {
		if w.WarningType() == SEVERE {
			return 0, w
		}
		warning = appendMinorWarnings(warning, w)
	}
	heelsMidpoint := getMidpoint(convertCvKeypointToPoint(keypoints.LHeel), convertCvKeypointToPoint(keypoints.RHeel))
	return getFootFlare(keypoints.LHeel, keypoints.LBigToe, calibrationInfo, heelsMidpoint), warning
}

func GetRightFootFlare(keypoints *cv.Body25PoseKeypoints, calibrationInfo *CalibrationInfo) (float64, warning) {
	var warning warning
	if w := verifyKeypoint(keypoints.LHeel, "left heel", 0.5); w != nil {
		if w.WarningType() == SEVERE {
			return 0, w
		}
		warning = appendMinorWarnings(warning, w)
	}
	if w := verifyKeypoint(keypoints.RHeel, "right heel", 0.5); w != nil {
		if w.WarningType() == SEVERE {
			return 0, w
		}
		warning = appendMinorWarnings(warning, w)
	}
	if w := verifyKeypoint(keypoints.RBigToe, "right big toe", 0.5); w != nil {
		if w.WarningType() == SEVERE {
			return 0, w
		}
		warning = appendMinorWarnings(warning, w)
	}
	heelsMidpoint := getMidpoint(convertCvKeypointToPoint(keypoints.LHeel), convertCvKeypointToPoint(keypoints.RHeel))
	return getFootFlare(keypoints.RHeel, keypoints.RBigToe, calibrationInfo, heelsMidpoint), warning
}

//stance width
//relative to hip to neck length
//line from left heel to right heel
//ratio of that line to hip to neck length
