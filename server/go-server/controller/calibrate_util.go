package controller

import (
	"context"
	"fmt"
	"math"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

func (g *GolfKeypointsListener) calibrateInputImageHelper(ctx context.Context, calibrationInfo *util.CalibrationInfo, inputImageId string, golfBall *skp.Keypoint, clubButt *skp.Keypoint, clubHead *skp.Keypoint, horAxisLine *skp.Line, vertAxisLine *skp.Line, shoulderTilt *skp.Double, firstLineAtTarget *skp.Line, secondLineAtTarget *skp.Line) (*skp.CalibrateInputImageResponse, error) {
	// get inputimg with inputimgid from db
	inputImage, err := g.dbmgr.ReadInputImage(ctx, inputImageId)
	if err != nil {
		return nil, fmt.Errorf("could not read input image with id: %s: %w", inputImageId, err)
	}
	// put in golf ball/golf club points if available
	if golfBall != nil {
		calibrationInfo.GolfBallPoint = *golfBall
	}
	if clubButt != nil {
		calibrationInfo.ClubButtPoint = *clubButt
	}
	if clubHead != nil {
		calibrationInfo.ClubHeadPoint = *clubHead
	}
	// calibrate based on dtl or face on
	if inputImage.ImageType == skp.ImageType_DTL {
		calibrationInfo, err = g.calibrateDTLImage(calibrationInfo, horAxisLine, vertAxisLine, shoulderTilt, firstLineAtTarget, secondLineAtTarget)
		if err != nil {
			return nil, fmt.Errorf("could not calibrate dtl image: %w", err)
		}
	} else {
		calibrationInfo, err = g.calibrateFaceOnImage(calibrationInfo, horAxisLine, vertAxisLine)
		if err != nil {
			return nil, fmt.Errorf("could not calibrate face on image: %w", err)
		}
	}
	inputImage.CalibrationInfo = *calibrationInfo
	// update inputimg with inputimgid in db
	_, err = g.dbmgr.UpdateInputImage(ctx, inputImageId, inputImage)
	if err != nil {
		return nil, fmt.Errorf("could not update input image with id: %s with calibration info: %w", inputImageId, err)
	}
	// return response
	response := &skp.CalibrateInputImageResponse{
		Success: true,
	}
	return response, nil
}

func (g *GolfKeypointsListener) calibrateDTLImage(calibrationInfo *util.CalibrationInfo, horAxisLine *skp.Line, vertAxisLine *skp.Line, shoulderTilt *skp.Double, firstLineAtTarget *skp.Line, secondLineAtTarget *skp.Line) (*util.CalibrationInfo, error) {
	var err error
	// axes calibration
	if calibrationInfo.CalibrationType != skp.CalibrationType_NO_CALIBRATION {
		if calibrationInfo.ManualGenerated {
			calibrationInfo, err = g.generateAxesFromLines(calibrationInfo, horAxisLine, vertAxisLine)
			if err != nil {
				return nil, fmt.Errorf("could not generate axes from lines: %w", err)
			}
		} else {
			calibrationInfo, err = g.generateAxesFromImage(calibrationInfo)
			if err != nil {
				return nil, fmt.Errorf("could not generate axes from image: %w", err)
			}
		}
		// vanishing point calibration
		if calibrationInfo.CalibrationType != skp.CalibrationType_AXES_CALIBRATION_ONLY {
			// add shoulder tilt for shoulder alignment calculation if provided
			if shoulderTilt != nil {
				calibrationInfo.ShoulderTilt = *shoulderTilt
			} else {
				calibrationInfo.ShoulderTilt = skp.Double{Data: 0, Warning: "Shoulder tilt not provided"}
			}
			if calibrationInfo.ManualGenerated {
				calibrationInfo, err = g.generateVanishingPointFromLines(calibrationInfo, firstLineAtTarget, secondLineAtTarget)
				if err != nil {
					return nil, fmt.Errorf("could not generate vanishing point from lines: %w", err)
				}
			} else {
				calibrationInfo, err = g.generateVanishingPointFromImage(calibrationInfo)
				if err != nil {
					return nil, fmt.Errorf("could not generate vanishing point from image: %w", err)
				}
			}
		}
	}
	return calibrationInfo, nil
}

func (g *GolfKeypointsListener) calibrateFaceOnImage(calibrationInfo *util.CalibrationInfo, horAxisLine *skp.Line, vertAxisLine *skp.Line) (*util.CalibrationInfo, error) {
	var err error
	// axes calibration
	if calibrationInfo.CalibrationType != skp.CalibrationType_NO_CALIBRATION {
		if calibrationInfo.ManualGenerated {
			calibrationInfo, err = g.generateAxesFromLines(calibrationInfo, horAxisLine, vertAxisLine)
			if err != nil {
				return nil, fmt.Errorf("could not generate axes from lines: %w", err)
			}
		} else {
			calibrationInfo, err = g.generateAxesFromImage(calibrationInfo)
			if err != nil {
				return nil, fmt.Errorf("could not generate axes from image: %w", err)
			}
		}
	}
	return calibrationInfo, nil
}

func (g *GolfKeypointsListener) generateAxesFromImage(calibrationInfo *util.CalibrationInfo) (*util.CalibrationInfo, error) {
	if calibrationInfo.CalibrationImgAxes == nil {
		return nil, fmt.Errorf("calibration image axes is required")
	}
	getPoseDataResponse, err := g.cvmgr.GetPoseData(calibrationInfo.CalibrationImgAxes)
	if err != nil {
		return nil, fmt.Errorf("could not get pose data for calibration image axes %w", err)
	}
	fmt.Printf("Axes calibration image processed\n")
	keypoints := getPoseDataResponse.Keypoints
	var warning util.Warning
	// Get horizontal axis
	feetLine, warning := util.GetFeetLine(keypoints, calibrationInfo.FeetLineMethod)
	if warning != nil {
		return nil, warning
	}
	horAxisLine := feetLine.Line
	// Get vertical axis
	if warning := util.VerifyKeypoint(keypoints.Midhip, "midhip", 0.5); warning != nil {
		return nil, warning
	}
	if warning := util.VerifyKeypoint(keypoints.Neck, "neck", 0.5); warning != nil {
		return nil, warning
	}
	vertAxisLine := util.GetLine(util.ConvertKeypointToPoint(keypoints.Midhip), util.ConvertKeypointToPoint(keypoints.Neck))
	// validate axes angle
	err = validateAxesAngle(&horAxisLine, vertAxisLine)
	if err != nil {
		return nil, fmt.Errorf("could not validate axes angle: %w. please adjust camera, stance, or posture. recommend using alignment sticks to help calibration", err)
	}
	// set axes
	calibrationInfo.HorAxisLine = horAxisLine
	calibrationInfo.VertAxisLine = *vertAxisLine
	return calibrationInfo, nil
}

func (g *GolfKeypointsListener) generateAxesFromLines(calibrationInfo *util.CalibrationInfo, horAxisLine *skp.Line, vertAxisLine *skp.Line) (*util.CalibrationInfo, error) {
	if horAxisLine == nil || vertAxisLine == nil {
		return nil, fmt.Errorf("both horizontal and vertical axes are required")
	}
	calibrationInfo.HorAxisLine = *util.ConvertSkpLineToLine(horAxisLine)
	calibrationInfo.VertAxisLine = *util.ConvertSkpLineToLine(vertAxisLine)
	err := validateAxesAngle(&calibrationInfo.HorAxisLine, &calibrationInfo.VertAxisLine)
	if err != nil {
		return nil, fmt.Errorf("could not validate axes angle: %w", err)
	}
	return calibrationInfo, nil
}

func (g *GolfKeypointsListener) generateVanishingPointFromImage(calibrationInfo *util.CalibrationInfo) (*util.CalibrationInfo, error) {
	if calibrationInfo.CalibrationImgVanishingPoint == nil {
		return nil, fmt.Errorf("calibration image vanishing point is required")
	}
	getPoseDataResponse, err := g.cvmgr.GetPoseData(calibrationInfo.CalibrationImgVanishingPoint)
	if err != nil {
		return nil, fmt.Errorf("could not get pose data for calibration image vanishingpoint %w", err)
	}
	fmt.Printf("Vanishing point calibration image processed\n")
	keypoints := getPoseDataResponse.Keypoints
	var warning util.Warning
	//use feet line as one parallel line at target and vert line as other parallel line at target
	feetLine, warning := util.GetFeetLine(keypoints, calibrationInfo.FeetLineMethod)
	if warning != nil {
		return nil, warning
	}
	// validate lines
	err = validateVanishingPointLines(&feetLine.Line, &calibrationInfo.VertAxisLine)
	if err != nil {
		return nil, fmt.Errorf("could not validate vanishing point lines: %w. make sure feet line is off centered or make sure alignment stick is pointed at target", err)
	}
	// get intersection
	intersection := util.GetIntersection(&feetLine.Line, &calibrationInfo.VertAxisLine)
	calibrationInfo.VanishingPoint = intersection.IntersectPoint
	return calibrationInfo, nil
}

func (g *GolfKeypointsListener) generateVanishingPointFromLines(calibrationInfo *util.CalibrationInfo, firstLineAtTarget *skp.Line, secondLineAtTarget *skp.Line) (*util.CalibrationInfo, error) {
	if firstLineAtTarget == nil || secondLineAtTarget == nil {
		return nil, fmt.Errorf("both lines at target are required")
	}
	line1 := util.ConvertSkpLineToLine(firstLineAtTarget)
	line2 := util.ConvertSkpLineToLine(secondLineAtTarget)
	// validate lines
	err := validateVanishingPointLines(line1, line2)
	if err != nil {
		return nil, fmt.Errorf("could not validate vanishing point lines: %w", err)
	}
	// get intersection
	intersection := util.GetIntersection(line1, line2)
	calibrationInfo.VanishingPoint = intersection.IntersectPoint
	return calibrationInfo, nil
}

// Checks if angle between axes is around 90 degrees
func validateAxesAngle(horAxisLine *util.Line, vertAxisLine *util.Line) error {
	horDeg := util.ConvertSlopeToDegrees(horAxisLine.Slope)
	vertDeg := util.ConvertSlopeToDegrees(vertAxisLine.Slope)
	diff := math.Abs(vertDeg) + math.Abs(horDeg) - 90
	if math.Abs(diff) > 10 { // TODO: Configure confidence level, configure how far off 90 degrees axes can be
		return fmt.Errorf("axes are off. horizontal axis is %f degrees. vertical axis is %f degrees. offset from 90 degrees of %f is too large.", horDeg, vertDeg, diff)
	}
	fmt.Printf("Good axes calibration. Horizontal axis is %f degrees. vertical axis is %f degrees\n", horDeg, vertDeg)
	return nil
}

// Checks that vanishing point lines are not too close in slope
func validateVanishingPointLines(line1 *util.Line, line2 *util.Line) error {
	slopeDiff := math.Abs(line1.Slope - line2.Slope)
	if slopeDiff < float64(1) { // TODO: Configure how close slope is (and how to determine how close slope is)
		return fmt.Errorf("vanishing point lines are off. first line slope is %f, second line slope is %f, these slopes are too close (%f). make sure at least one line is off centered (parallel lines converge in distance)", line1.Slope, line2.Slope, slopeDiff)
	}
	fmt.Printf("Good vanishing point calibration. line1 is %+v, and line2 is %+v", line1, line2)
	return nil
}
