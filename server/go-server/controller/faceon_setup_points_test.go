package controller

import (
	"math"
	"testing"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

var golfSpecificDatapointsFaceOn *skp.GolfSpecificDatapoints = &skp.GolfSpecificDatapoints{
	GolfBall: &skp.Datapoint{
		X:          528,
		Y:          1944,
		Confidence: 1.0,
	},
	ClubButt: &skp.Datapoint{
		X:          528,
		Y:          1480,
		Confidence: 1.0,
	},
	ClubHead: &skp.Datapoint{
		X:          508,
		Y:          1940,
		Confidence: 1.0,
	},
}

var calibrationInfoFaceOn *util.CalibrationInfo = &util.CalibrationInfo{
	CalibrationType: skp.CalibrationType_FULL_CALIBRATION,
	FeetLineMethod:  skp.FeetLineMethod_USE_HEEL_LINE,
	HorAxisLine: util.Line{
		Slope:      0.0000352,
		YIntercept: 1767.522,
		PointOnLine: util.Point{
			XPos: 619.025,
			YPos: 1767.543,
		},
	},
	VertAxisLine: util.Line{
		Slope:      1349.033,
		YIntercept: -686797.244,
		PointOnLine: util.Point{
			XPos: 510.171,
			YPos: 1441.491,
		},
	},
}

func TestGetSideBend(t *testing.T) {
	// neutral side bend
	bodyDatapoints := &skp.Body25PoseDatapoints{
		Midhip: &skp.Datapoint{
			X:          510.483,
			Y:          1441.562,
			Confidence: 1.0,
		},
		Neck: &skp.Datapoint{
			X:          510.222,
			Y:          1216.812,
			Confidence: 1.0,
		},
	}
	neutralExpected := 0.024
	neutralActual, warning := GetSideBend(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetSideBend(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetSideBend(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, neutralActual, neutralExpected)
	}
	// right side bend -> shift neck farther right
	bodyDatapoints.Neck = &skp.Datapoint{
		X:          410.222,
		Y:          1216.812,
		Confidence: 1.0,
	}
	rightExpected := 23.999
	rightActual, warning := GetSideBend(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetSideBend(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if rightActual <= neutralActual {
		t.Errorf("GetSideBend(%+v, %+v) right bend result %f is supposed to be greater than the neutral result %f", bodyDatapoints, calibrationInfoFaceOn, rightActual, neutralActual)
	}
	if math.Abs(rightActual-rightExpected) > 0.01 {
		t.Errorf("GetSideBend(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, rightActual, rightExpected)
	}
	// left side bend -> shift neck farther left
	bodyDatapoints.Neck = &skp.Datapoint{
		X:          610.222,
		Y:          1216.812,
		Confidence: 1.0,
	}
	leftExpected := -23.973
	leftActual, warning := GetSideBend(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetSideBend(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if leftActual >= neutralActual {
		t.Errorf("GetSideBend(%+v, %+v) left bend result %f is supposed to be less than the neutral result %f", bodyDatapoints, calibrationInfoFaceOn, leftActual, neutralActual)
	}
	if math.Abs(leftActual-leftExpected) > 0.01 {
		t.Errorf("GetSideBend(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, leftActual, leftExpected)
	}
}

func TestGetLeftFootFlare(t *testing.T) {
	// normal flare
	bodyDatapoints := &skp.Body25PoseDatapoints{
		LBigToe: &skp.Datapoint{
			X:          619.099,
			Y:          1814.907,
			Confidence: 1.0,
		},
		LHeel: &skp.Datapoint{
			X:          618.136,
			Y:          1760.744,
			Confidence: 1.0,
		},
	}
	normalExpected := 0.976
	normalActual, warning := GetLeftFootFlare(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetLeftFootFlare(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(normalActual-normalExpected) > 0.01 {
		t.Errorf("GetLeftFootFlare(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, normalActual, normalExpected)
	}
	// external -> shift left big toe out
	bodyDatapoints.LBigToe = &skp.Datapoint{
		X:          645.099,
		Y:          1814.907,
		Confidence: 1.0,
	}
	externalExpected := 26.422
	externalActual, warning := GetLeftFootFlare(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetLeftFootFlare(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if externalActual <= normalActual {
		t.Errorf("GetLeftFootFlare(%+v, %+v) external result %f is supposed to be greater than the normal result %f", bodyDatapoints, calibrationInfoFaceOn, externalActual, normalActual)
	}
	if math.Abs(externalActual-externalExpected) > 0.01 {
		t.Errorf("GetLeftFootFlare(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, externalActual, externalExpected)
	}
	// internal -> shift left big toe in
	bodyDatapoints.LBigToe = &skp.Datapoint{
		X:          601.099,
		Y:          1814.907,
		Confidence: 1.0,
	}
	internalExpected := -17.503
	internalActual, warning := GetLeftFootFlare(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetLeftFootFlare(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if internalActual >= normalActual {
		t.Errorf("GetLeftFootFlare(%+v, %+v) internal result %f is supposed to be less than the normal result %f", bodyDatapoints, calibrationInfoFaceOn, internalActual, normalActual)
	}
	if math.Abs(internalActual-internalExpected) > 0.01 {
		t.Errorf("GetLeftFootFlare(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, internalActual, internalExpected)
	}
}

func TestGetRightFootFlare(t *testing.T) {
	// normal flare
	bodyDatapoints := &skp.Body25PoseDatapoints{
		RBigToe: &skp.Datapoint{
			X:          415.270,
			Y:          1821.722,
			Confidence: 1.0,
		},
		RHeel: &skp.Datapoint{
			X:          415.120,
			Y:          1767.229,
			Confidence: 1.0,
		},
	}
	normalExpected := -0.115
	normalActual, warning := GetRightFootFlare(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetRightFootFlare(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(normalActual-normalExpected) > 0.01 {
		t.Errorf("GetRightFootFlare(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, normalActual, normalExpected)
	}
	// external -> shift right big toe out
	bodyDatapoints.RBigToe = &skp.Datapoint{
		X:          385.270,
		Y:          1821.722,
		Confidence: 1.0,
	}
	externalExpected := 28.755
	externalActual, warning := GetRightFootFlare(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetRightFootFlare(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if externalActual <= normalActual {
		t.Errorf("GetRightFootFlare(%+v, %+v) external result %f is supposed to be greater than the normal result %f", bodyDatapoints, calibrationInfoFaceOn, externalActual, normalActual)
	}
	if math.Abs(externalActual-externalExpected) > 0.01 {
		t.Errorf("GetRightFootFlare(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, externalActual, externalExpected)
	}
	// internal -> shift right big toe in
	bodyDatapoints.RBigToe = &skp.Datapoint{
		X:          435.270,
		Y:          1821.722,
		Confidence: 1.0,
	}
	internalExpected := -20.250
	internalActual, warning := GetRightFootFlare(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetRightFootFlare(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if internalActual >= normalActual {
		t.Errorf("GetRightFootFlare(%+v, %+v) internal result %f is supposed to be less than the normal result %f", bodyDatapoints, calibrationInfoFaceOn, internalActual, normalActual)
	}
	if math.Abs(internalActual-internalExpected) > 0.01 {
		t.Errorf("GetRightFootFlare(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, internalActual, internalExpected)
	}
}

func TestGetStanceWidth(t *testing.T) {
	// normal width
	bodyDatapoints := &skp.Body25PoseDatapoints{
		LHeel: &skp.Datapoint{
			X:          605.136,
			Y:          1760.744,
			Confidence: 1.0,
		},
		RHeel: &skp.Datapoint{
			X:          415.120,
			Y:          1767.229,
			Confidence: 1.0,
		},
		Midhip: &skp.Datapoint{
			X:          510.483,
			Y:          1441.562,
			Confidence: 1.0,
		},
		Neck: &skp.Datapoint{
			X:          510.222,
			Y:          1216.812,
			Confidence: 1.0,
		},
	}
	normalExpected := 0.845
	normalActual, warning := GetStanceWidth(bodyDatapoints)
	if warning != nil {
		t.Errorf("GetStanceWidth(%+v) has an unexpected warning: %v", bodyDatapoints, warning)
	}
	if math.Abs(normalActual-normalExpected) > 0.01 {
		t.Errorf("GetStanceWidth(%+v) = %f; expected %f", bodyDatapoints, normalActual, normalExpected)
	}
	// wide -> shift left heel farther
	bodyDatapoints.LHeel = &skp.Datapoint{
		X:          705.136,
		Y:          1760.744,
		Confidence: 1.0,
	}
	wideExpected := 1.290
	wideActual, warning := GetStanceWidth(bodyDatapoints)
	if warning != nil {
		t.Errorf("GetStanceWidth(%+v) has an unexpected warning: %v", bodyDatapoints, warning)
	}
	if wideActual <= normalActual {
		t.Errorf("GetStanceWidth(%+v) wide result %f is supposed to be greater than the normal result %f", bodyDatapoints, wideActual, normalActual)
	}
	if math.Abs(wideActual-wideExpected) > 0.01 {
		t.Errorf("GetStanceWidth(%+v) = %f; expected %f", bodyDatapoints, wideActual, wideExpected)
	}
	// narrow -> shift left heel closer
	bodyDatapoints.LHeel = &skp.Datapoint{
		X:          508.136,
		Y:          1760.744,
		Confidence: 1.0,
	}
	narrowExpected := 0.414
	narrowActual, warning := GetStanceWidth(bodyDatapoints)
	if warning != nil {
		t.Errorf("GetStanceWidth(%+v) has an unexpected warning: %v", bodyDatapoints, warning)
	}
	if narrowActual >= normalActual {
		t.Errorf("GetStanceWidth(%+v) narrow result %f is supposed to be less than the normal result %f", bodyDatapoints, narrowActual, normalActual)
	}
	if math.Abs(narrowActual-narrowExpected) > 0.01 {
		t.Errorf("GetStanceWidth(%+v) = %f; expected %f", bodyDatapoints, narrowActual, narrowExpected)
	}
}

func TestGetShoulderTilt(t *testing.T) {
	// neutral tilt
	bodyDatapoints := &skp.Body25PoseDatapoints{
		LShoulder: &skp.Datapoint{
			X:          578.036,
			Y:          1210.232,
			Confidence: 1.0,
		},
		RShoulder: &skp.Datapoint{
			X:          428.688,
			Y:          1217.123,
			Confidence: 1.0,
		},
	}
	neutralExpected := 2.643
	neutralActual, warning := GetShoulderTilt(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetShoulderTilt(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetShoulderTilt(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, neutralActual, neutralExpected)
	}
	// right tilt -> raise left shoulder
	bodyDatapoints.LShoulder = &skp.Datapoint{
		X:          578.036,
		Y:          1123.232,
		Confidence: 1.0,
	}
	rightTiltExpected := 32.158
	rightTiltActual, warning := GetShoulderTilt(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetShoulderTilt(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if rightTiltActual <= neutralActual {
		t.Errorf("GetShoulderTilt(%+v, %+v) right tilt result %f is supposed to be greater than the neutral result %f", bodyDatapoints, calibrationInfoFaceOn, rightTiltActual, neutralActual)
	}
	if math.Abs(rightTiltActual-rightTiltExpected) > 0.01 {
		t.Errorf("GetShoulderTilt(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, rightTiltActual, rightTiltExpected)
	}
	// left tilt -> lower left shoulder
	bodyDatapoints.LShoulder = &skp.Datapoint{
		X:          578.036,
		Y:          1324.232,
		Confidence: 1.0,
	}
	leftTiltExpected := -35.645
	leftTiltActual, warning := GetShoulderTilt(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetShoulderTilt(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if leftTiltActual >= neutralActual {
		t.Errorf("GetShoulderTilt(%+v, %+v) left tilt result %f is supposed to be less than the neutral result %f", bodyDatapoints, calibrationInfoFaceOn, leftTiltActual, neutralActual)
	}
	if math.Abs(leftTiltActual-leftTiltExpected) > 0.01 {
		t.Errorf("GetShoulderTilt(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, leftTiltActual, leftTiltExpected)
	}
}

func TestGetWaistTilt(t *testing.T) {
	// neutral tilt
	bodyDatapoints := &skp.Body25PoseDatapoints{
		LHip: &skp.Datapoint{
			X:          564.604,
			Y:          1441.614,
			Confidence: 1.0,
		},
		RHip: &skp.Datapoint{
			X:          462.534,
			Y:          1441.622,
			Confidence: 1.0,
		},
	}
	neutralExpected := 0.0065
	neutralActual, warning := GetWaistTilt(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetWaistTilt(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetWaistTilt(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, neutralActual, neutralExpected)
	}
	// right tilt -> raise left hip
	bodyDatapoints.LHip = &skp.Datapoint{
		X:          564.604,
		Y:          1341.614,
		Confidence: 1.0,
	}
	rightTiltExpected := 44.417
	rightTiltActual, warning := GetWaistTilt(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetWaistTilt(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if rightTiltActual <= neutralActual {
		t.Errorf("GetWaistTilt(%+v, %+v) right tilt result %f is supposed to be greater than the neutral result %f", bodyDatapoints, calibrationInfoFaceOn, rightTiltActual, neutralActual)
	}
	if math.Abs(rightTiltActual-rightTiltExpected) > 0.01 {
		t.Errorf("GetWaistTilt(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, rightTiltActual, rightTiltExpected)
	}
	// left tilt -> lower left hip
	bodyDatapoints.LHip = &skp.Datapoint{
		X:          564.604,
		Y:          1563.614,
		Confidence: 1.0,
	}
	leftTiltExpected := -50.078
	leftTiltActual, warning := GetWaistTilt(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetWaistTilt(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if leftTiltActual >= neutralActual {
		t.Errorf("GetWaistTilt(%+v, %+v) left tilt result %f is supposed to be less than the neutral result %f", bodyDatapoints, calibrationInfoFaceOn, leftTiltActual, neutralActual)
	}
	if math.Abs(leftTiltActual-leftTiltExpected) > 0.01 {
		t.Errorf("GetWaistTilt(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, leftTiltActual, leftTiltExpected)
	}
}

func TestGetShaftLean(t *testing.T) {
	// neutral shaft lean
	centralClubButt := golfSpecificDatapointsFaceOn.ClubButt
	centralExpected := 2.532
	centralActual, warning := GetShaftLean(golfSpecificDatapointsFaceOn, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetShaftLean(%+v, %+v) has an unexpected warning: %v", golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, warning)
	}
	if math.Abs(centralActual-centralExpected) > 0.01 {
		t.Errorf("GetShaftLean(%+v, %+v) = %f; expected %f", golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, centralActual, centralExpected)
	}
	// forward -> shift club butt forward
	golfSpecificDatapointsFaceOn.ClubButt = &skp.Datapoint{
		X:          632,
		Y:          1480,
		Confidence: 1.0,
	}
	forwardExpected := 15.128
	forwardActual, warning := GetShaftLean(golfSpecificDatapointsFaceOn, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetShaftLean(%+v, %+v) has an unexpected warning: %v", golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, warning)
	}
	if forwardActual <= centralActual {
		t.Errorf("GetShaftLean(%+v, %+v) forward result %f is supposed to be greater than the central result %f", golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, forwardActual, centralActual)
	}
	if math.Abs(forwardActual-forwardExpected) > 0.01 {
		t.Errorf("GetShaftLean(%+v, %+v) = %f; expected %f", golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, forwardActual, forwardExpected)
	}
	// backward -> shift club butt backwards
	golfSpecificDatapointsFaceOn.ClubButt = &skp.Datapoint{
		X:          422,
		Y:          1480,
		Confidence: 1.0,
	}
	backwardExpected := -10.547
	backwardActual, warning := GetShaftLean(golfSpecificDatapointsFaceOn, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetShaftLean(%+v, %+v) has an unexpected warning: %v", golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, warning)
	}
	if backwardActual >= centralActual {
		t.Errorf("GetShaftLean(%+v, %+v) backward result %f is supposed to be less than the central result %f", golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, backwardActual, centralActual)
	}
	if math.Abs(backwardActual-backwardExpected) > 0.01 {
		t.Errorf("GetShaftLean(%+v, %+v) = %f; expected %f", golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, backwardActual, backwardExpected)
	}
	// reset club head point
	golfSpecificDatapointsFaceOn.ClubButt = centralClubButt
}

func TestGetBallPosition(t *testing.T) {
	// central ball position
	centralGolfBall := golfSpecificDatapointsFaceOn.GolfBall
	bodyDatapoints := &skp.Body25PoseDatapoints{
		LHeel: &skp.Datapoint{
			X:          605.136,
			Y:          1760.744,
			Confidence: 1.0,
		},
		RHeel: &skp.Datapoint{
			X:          415.120,
			Y:          1767.229,
			Confidence: 1.0,
		},
	}
	centralExpected := 3.715
	centralActual, warning := GetBallPosition(bodyDatapoints, golfSpecificDatapointsFaceOn, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetBallPosition(%+v, %+v, %+v) has an unexpected warning: %v", bodyDatapoints, golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, warning)
	}
	if math.Abs(centralActual-centralExpected) > 0.01 {
		t.Errorf("GetBallPosition(%+v, %+v, %+v) = %f; expected %f", bodyDatapoints, golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, centralActual, centralExpected)
	}
	// forward -> shift golf ball forward
	golfSpecificDatapointsFaceOn.GolfBall = &skp.Datapoint{
		X:          632,
		Y:          1944,
		Confidence: 1.0,
	}
	forwardExpected := 32.143
	forwardActual, warning := GetBallPosition(bodyDatapoints, golfSpecificDatapointsFaceOn, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetBallPosition(%+v, %+v, %+v) has an unexpected warning: %v", bodyDatapoints, golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, warning)
	}
	if forwardActual <= centralActual {
		t.Errorf("GetBallPosition(%+v, %+v, %+v) forward result %f is supposed to be greater than the central result %f", bodyDatapoints, golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, forwardActual, centralActual)
	}
	if math.Abs(forwardActual-forwardExpected) > 0.01 {
		t.Errorf("GetBallPosition(%+v, %+v, %+v) = %f; expected %f", bodyDatapoints, golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, forwardActual, forwardExpected)
	}
	// back -> shift golf ball back
	golfSpecificDatapointsFaceOn.GolfBall = &skp.Datapoint{
		X:          445,
		Y:          1944,
		Confidence: 1.0,
	}
	backExpected := -21.844
	backActual, warning := GetBallPosition(bodyDatapoints, golfSpecificDatapointsFaceOn, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetBallPosition(%+v, %+v, %+v) has an unexpected warning: %v", bodyDatapoints, golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, warning)
	}
	if backActual >= centralActual {
		t.Errorf("GetBallPosition(%+v, %+v, %+v) back result %f is supposed to be less than the central result %f", bodyDatapoints, golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, backActual, centralActual)
	}
	if math.Abs(backActual-backExpected) > 0.01 {
		t.Errorf("GetBallPosition(%+v, %+v, %+v) = %f; expected %f", bodyDatapoints, golfSpecificDatapointsFaceOn, calibrationInfoFaceOn, backActual, backExpected)
	}
	// reset golf ball point
	golfSpecificDatapointsFaceOn.GolfBall = centralGolfBall
}

func TestGetHeadPosition(t *testing.T) {
	// central head position
	bodyDatapoints := &skp.Body25PoseDatapoints{
		Nose: &skp.Datapoint{
			X:          510.237,
			Y:          1155.832,
			Confidence: 1.0,
		},
		LHeel: &skp.Datapoint{
			X:          605.136,
			Y:          1760.744,
			Confidence: 1.0,
		},
		RHeel: &skp.Datapoint{
			X:          415.120,
			Y:          1767.229,
			Confidence: 1.0,
		},
	}
	centralExpected := 1.964
	centralActual, warning := GetHeadPosition(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetHeadPosition(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(centralActual-centralExpected) > 0.01 {
		t.Errorf("GetHeadPosition(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, centralActual, centralExpected)
	}
	// forward -> shift nose forward
	bodyDatapoints.Nose = &skp.Datapoint{
		X:          603.237,
		Y:          1155.832,
		Confidence: 1.0,
	}
	forwardExpected := 10.659
	forwardActual, warning := GetHeadPosition(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetHeadPosition(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if forwardActual <= centralActual {
		t.Errorf("GetHeadPosition(%+v, %+v) forward result %f is supposed to be greater than the central result %f", bodyDatapoints, calibrationInfoFaceOn, forwardActual, centralActual)
	}
	if math.Abs(forwardActual-forwardExpected) > 0.01 {
		t.Errorf("GetHeadPosition(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, forwardActual, forwardExpected)
	}
	// back -> shift nose backward
	bodyDatapoints.Nose = &skp.Datapoint{
		X:          412.237,
		Y:          1155.832,
		Confidence: 1.0,
	}
	backExpected := -7.189
	backActual, warning := GetHeadPosition(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetHeadPosition(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if backActual >= centralActual {
		t.Errorf("GetHeadPosition(%+v, %+v) back result %f is supposed to be less than the central result %f", bodyDatapoints, calibrationInfoFaceOn, backActual, centralActual)
	}
	if math.Abs(backActual-backExpected) > 0.01 {
		t.Errorf("GetHeadPosition(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, backActual, backExpected)
	}
}

func TestGetChestPosition(t *testing.T) {
	// central chest position
	bodyDatapoints := &skp.Body25PoseDatapoints{
		Neck: &skp.Datapoint{
			X:          510.222,
			Y:          1216.812,
			Confidence: 1.0,
		},
		LHeel: &skp.Datapoint{
			X:          605.136,
			Y:          1760.744,
			Confidence: 1.0,
		},
		RHeel: &skp.Datapoint{
			X:          415.120,
			Y:          1767.229,
			Confidence: 1.0,
		},
	}
	centralExpected := 1.964
	centralActual, warning := GetChestPosition(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetChestPosition(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(centralActual-centralExpected) > 0.01 {
		t.Errorf("GetChestPosition(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, centralActual, centralExpected)
	}
	// forward -> shift neck forward
	bodyDatapoints.Neck = &skp.Datapoint{
		X:          634.222,
		Y:          1216.812,
		Confidence: 1.0,
	}
	forwardExpected := 14.732
	forwardActual, warning := GetChestPosition(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetChestPosition(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if forwardActual <= centralActual {
		t.Errorf("GetChestPosition(%+v, %+v) forward result %f is supposed to be greater than the central result %f", bodyDatapoints, calibrationInfoFaceOn, forwardActual, centralActual)
	}
	if math.Abs(forwardActual-forwardExpected) > 0.01 {
		t.Errorf("GetChestPosition(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, forwardActual, forwardExpected)
	}
	// back -> shift neck backward
	bodyDatapoints.Neck = &skp.Datapoint{
		X:          422.222,
		Y:          1216.812,
		Confidence: 1.0,
	}
	backExpected := -7.172
	backActual, warning := GetChestPosition(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetChestPosition(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if backActual >= centralActual {
		t.Errorf("GetChestPosition(%+v, %+v) back result %f is supposed to be less than the central result %f", bodyDatapoints, calibrationInfoFaceOn, backActual, centralActual)
	}
	if math.Abs(backActual-backExpected) > 0.01 {
		t.Errorf("GetChestPosition(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, backActual, backExpected)
	}
}

func TestGetMidhipPosition(t *testing.T) {
	// central midhip position
	bodyDatapoints := &skp.Body25PoseDatapoints{
		Midhip: &skp.Datapoint{
			X:          510.483,
			Y:          1441.562,
			Confidence: 1.0,
		},
		LHeel: &skp.Datapoint{
			X:          605.136,
			Y:          1760.744,
			Confidence: 1.0,
		},
		RHeel: &skp.Datapoint{
			X:          415.120,
			Y:          1767.229,
			Confidence: 1.0,
		},
	}
	centralExpected := 2.017
	centralActual, warning := GetMidhipPosition(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetMidhipPosition(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(centralActual-centralExpected) > 0.01 {
		t.Errorf("GetMidhipPosition(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, centralActual, centralExpected)
	}
	// forward -> shift midhip forward
	bodyDatapoints.Midhip = &skp.Datapoint{
		X:          632.483,
		Y:          1441.562,
		Confidence: 1.0,
	}
	forwardExpected := 22.735
	forwardActual, warning := GetMidhipPosition(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetMidhipPosition(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if forwardActual <= centralActual {
		t.Errorf("GetMidhipPosition(%+v, %+v) forward result %f is supposed to be greater than the central result %f", bodyDatapoints, calibrationInfoFaceOn, forwardActual, centralActual)
	}
	if math.Abs(forwardActual-forwardExpected) > 0.01 {
		t.Errorf("GetMidhipPosition(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, forwardActual, forwardExpected)
	}
	// back -> shift midhip backward
	bodyDatapoints.Midhip = &skp.Datapoint{
		X:          413.483,
		Y:          1441.562,
		Confidence: 1.0,
	}
	backExpected := -14.731
	backActual, warning := GetMidhipPosition(bodyDatapoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetMidhipPosition(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoFaceOn, warning)
	}
	if backActual >= centralActual {
		t.Errorf("GetMidhipPosition(%+v, %+v) back result %f is supposed to be less than the central result %f", bodyDatapoints, calibrationInfoFaceOn, backActual, centralActual)
	}
	if math.Abs(backActual-backExpected) > 0.01 {
		t.Errorf("GetMidhipPosition(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoFaceOn, backActual, backExpected)
	}
}
