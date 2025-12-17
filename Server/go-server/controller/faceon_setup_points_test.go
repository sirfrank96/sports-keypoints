package controller

import (
	"math"
	"testing"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

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
	GolfBallPoint: skp.Keypoint{
		X:          528,
		Y:          1944,
		Confidence: 1.0,
	},
	ClubButtPoint: skp.Keypoint{
		X:          528,
		Y:          1480,
		Confidence: 1.0,
	},
	ClubHeadPoint: skp.Keypoint{
		X:          508,
		Y:          1940,
		Confidence: 1.0,
	},
}

func TestGetSideBend(t *testing.T) {
	// neutral side bend
	keypoints := &skp.Body25PoseKeypoints{
		Midhip: &skp.Keypoint{
			X:          510.483,
			Y:          1441.562,
			Confidence: 1.0,
		},
		Neck: &skp.Keypoint{
			X:          510.222,
			Y:          1216.812,
			Confidence: 1.0,
		},
	}
	neutralExpected := 0.024
	neutralActual, warning := GetSideBend(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetSideBend(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetSideBend(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, neutralActual, neutralExpected)
	}
	// right side bend -> shift neck farther right
	keypoints.Neck = &skp.Keypoint{
		X:          410.222,
		Y:          1216.812,
		Confidence: 1.0,
	}
	rightExpected := 23.999
	rightActual, warning := GetSideBend(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetSideBend(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if rightActual <= neutralActual {
		t.Errorf("GetSideBend(%+v, %+v) right bend result %f is supposed to be greater than the neutral result %f", keypoints, calibrationInfoFaceOn, rightActual, neutralActual)
	}
	if math.Abs(rightActual-rightExpected) > 0.01 {
		t.Errorf("GetSideBend(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, rightActual, rightExpected)
	}
	// left side bend -> shift neck farther left
	keypoints.Neck = &skp.Keypoint{
		X:          610.222,
		Y:          1216.812,
		Confidence: 1.0,
	}
	leftExpected := -23.973
	leftActual, warning := GetSideBend(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetSideBend(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if leftActual >= neutralActual {
		t.Errorf("GetSideBend(%+v, %+v) left bend result %f is supposed to be less than the neutral result %f", keypoints, calibrationInfoFaceOn, leftActual, neutralActual)
	}
	if math.Abs(leftActual-leftExpected) > 0.01 {
		t.Errorf("GetSideBend(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, leftActual, leftExpected)
	}
}

func TestGetLeftFootFlare(t *testing.T) {
	// normal flare
	keypoints := &skp.Body25PoseKeypoints{
		LBigToe: &skp.Keypoint{
			X:          619.099,
			Y:          1814.907,
			Confidence: 1.0,
		},
		LHeel: &skp.Keypoint{
			X:          618.136,
			Y:          1760.744,
			Confidence: 1.0,
		},
	}
	normalExpected := 0.976
	normalActual, warning := GetLeftFootFlare(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetLeftFootFlare(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(normalActual-normalExpected) > 0.01 {
		t.Errorf("GetLeftFootFlare(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, normalActual, normalExpected)
	}
	// external -> shift left big toe out
	keypoints.LBigToe = &skp.Keypoint{
		X:          645.099,
		Y:          1814.907,
		Confidence: 1.0,
	}
	externalExpected := 26.422
	externalActual, warning := GetLeftFootFlare(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetLeftFootFlare(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if externalActual <= normalActual {
		t.Errorf("GetLeftFootFlare(%+v, %+v) external result %f is supposed to be greater than the normal result %f", keypoints, calibrationInfoFaceOn, externalActual, normalActual)
	}
	if math.Abs(externalActual-externalExpected) > 0.01 {
		t.Errorf("GetLeftFootFlare(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, externalActual, externalExpected)
	}
	// internal -> shift left big toe in
	keypoints.LBigToe = &skp.Keypoint{
		X:          601.099,
		Y:          1814.907,
		Confidence: 1.0,
	}
	internalExpected := -17.503
	internalActual, warning := GetLeftFootFlare(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetLeftFootFlare(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if internalActual >= normalActual {
		t.Errorf("GetLeftFootFlare(%+v, %+v) internal result %f is supposed to be less than the normal result %f", keypoints, calibrationInfoFaceOn, internalActual, normalActual)
	}
	if math.Abs(internalActual-internalExpected) > 0.01 {
		t.Errorf("GetLeftFootFlare(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, internalActual, internalExpected)
	}
}

func TestGetRightFootFlare(t *testing.T) {
	// normal flare
	keypoints := &skp.Body25PoseKeypoints{
		RBigToe: &skp.Keypoint{
			X:          415.270,
			Y:          1821.722,
			Confidence: 1.0,
		},
		RHeel: &skp.Keypoint{
			X:          415.120,
			Y:          1767.229,
			Confidence: 1.0,
		},
	}
	normalExpected := -0.115
	normalActual, warning := GetRightFootFlare(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetRightFootFlare(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(normalActual-normalExpected) > 0.01 {
		t.Errorf("GetRightFootFlare(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, normalActual, normalExpected)
	}
	// external -> shift right big toe out
	keypoints.RBigToe = &skp.Keypoint{
		X:          385.270,
		Y:          1821.722,
		Confidence: 1.0,
	}
	externalExpected := 28.755
	externalActual, warning := GetRightFootFlare(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetRightFootFlare(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if externalActual <= normalActual {
		t.Errorf("GetRightFootFlare(%+v, %+v) external result %f is supposed to be greater than the normal result %f", keypoints, calibrationInfoFaceOn, externalActual, normalActual)
	}
	if math.Abs(externalActual-externalExpected) > 0.01 {
		t.Errorf("GetRightFootFlare(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, externalActual, externalExpected)
	}
	// internal -> shift right big toe in
	keypoints.RBigToe = &skp.Keypoint{
		X:          435.270,
		Y:          1821.722,
		Confidence: 1.0,
	}
	internalExpected := -20.250
	internalActual, warning := GetRightFootFlare(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetRightFootFlare(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if internalActual >= normalActual {
		t.Errorf("GetRightFootFlare(%+v, %+v) internal result %f is supposed to be less than the normal result %f", keypoints, calibrationInfoFaceOn, internalActual, normalActual)
	}
	if math.Abs(internalActual-internalExpected) > 0.01 {
		t.Errorf("GetRightFootFlare(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, internalActual, internalExpected)
	}
}

func TestGetStanceWidth(t *testing.T) {
	// normal width
	keypoints := &skp.Body25PoseKeypoints{
		LHeel: &skp.Keypoint{
			X:          605.136,
			Y:          1760.744,
			Confidence: 1.0,
		},
		RHeel: &skp.Keypoint{
			X:          415.120,
			Y:          1767.229,
			Confidence: 1.0,
		},
		Midhip: &skp.Keypoint{
			X:          510.483,
			Y:          1441.562,
			Confidence: 1.0,
		},
		Neck: &skp.Keypoint{
			X:          510.222,
			Y:          1216.812,
			Confidence: 1.0,
		},
	}
	normalExpected := 0.845
	normalActual, warning := GetStanceWidth(keypoints)
	if warning != nil {
		t.Errorf("GetStanceWidth(%+v) has an unexpected warning: %v", keypoints, warning)
	}
	if math.Abs(normalActual-normalExpected) > 0.01 {
		t.Errorf("GetStanceWidth(%+v) = %f; expected %f", keypoints, normalActual, normalExpected)
	}
	// wide -> shift left heel farther
	keypoints.LHeel = &skp.Keypoint{
		X:          705.136,
		Y:          1760.744,
		Confidence: 1.0,
	}
	wideExpected := 1.290
	wideActual, warning := GetStanceWidth(keypoints)
	if warning != nil {
		t.Errorf("GetStanceWidth(%+v) has an unexpected warning: %v", keypoints, warning)
	}
	if wideActual <= normalActual {
		t.Errorf("GetStanceWidth(%+v) wide result %f is supposed to be greater than the normal result %f", keypoints, wideActual, normalActual)
	}
	if math.Abs(wideActual-wideExpected) > 0.01 {
		t.Errorf("GetStanceWidth(%+v) = %f; expected %f", keypoints, wideActual, wideExpected)
	}
	// narrow -> shift left heel closer
	keypoints.LHeel = &skp.Keypoint{
		X:          508.136,
		Y:          1760.744,
		Confidence: 1.0,
	}
	narrowExpected := 0.414
	narrowActual, warning := GetStanceWidth(keypoints)
	if warning != nil {
		t.Errorf("GetStanceWidth(%+v) has an unexpected warning: %v", keypoints, warning)
	}
	if narrowActual >= normalActual {
		t.Errorf("GetStanceWidth(%+v) narrow result %f is supposed to be less than the normal result %f", keypoints, narrowActual, normalActual)
	}
	if math.Abs(narrowActual-narrowExpected) > 0.01 {
		t.Errorf("GetStanceWidth(%+v) = %f; expected %f", keypoints, narrowActual, narrowExpected)
	}
}

func TestGetShoulderTilt(t *testing.T) {
	// neutral tilt
	keypoints := &skp.Body25PoseKeypoints{
		LShoulder: &skp.Keypoint{
			X:          578.036,
			Y:          1210.232,
			Confidence: 1.0,
		},
		RShoulder: &skp.Keypoint{
			X:          428.688,
			Y:          1217.123,
			Confidence: 1.0,
		},
	}
	neutralExpected := 2.643
	neutralActual, warning := GetShoulderTilt(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetShoulderTilt(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetShoulderTilt(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, neutralActual, neutralExpected)
	}
	// right tilt -> raise left shoulder
	keypoints.LShoulder = &skp.Keypoint{
		X:          578.036,
		Y:          1123.232,
		Confidence: 1.0,
	}
	rightTiltExpected := 32.158
	rightTiltActual, warning := GetShoulderTilt(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetShoulderTilt(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if rightTiltActual <= neutralActual {
		t.Errorf("GetShoulderTilt(%+v, %+v) right tilt result %f is supposed to be greater than the neutral result %f", keypoints, calibrationInfoFaceOn, rightTiltActual, neutralActual)
	}
	if math.Abs(rightTiltActual-rightTiltExpected) > 0.01 {
		t.Errorf("GetShoulderTilt(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, rightTiltActual, rightTiltExpected)
	}
	// left tilt -> lower left shoulder
	keypoints.LShoulder = &skp.Keypoint{
		X:          578.036,
		Y:          1324.232,
		Confidence: 1.0,
	}
	leftTiltExpected := -35.645
	leftTiltActual, warning := GetShoulderTilt(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetShoulderTilt(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if leftTiltActual >= neutralActual {
		t.Errorf("GetShoulderTilt(%+v, %+v) left tilt result %f is supposed to be less than the neutral result %f", keypoints, calibrationInfoFaceOn, leftTiltActual, neutralActual)
	}
	if math.Abs(leftTiltActual-leftTiltExpected) > 0.01 {
		t.Errorf("GetShoulderTilt(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, leftTiltActual, leftTiltExpected)
	}
}

func TestGetWaistTilt(t *testing.T) {
	// neutral tilt
	keypoints := &skp.Body25PoseKeypoints{
		LHip: &skp.Keypoint{
			X:          564.604,
			Y:          1441.614,
			Confidence: 1.0,
		},
		RHip: &skp.Keypoint{
			X:          462.534,
			Y:          1441.622,
			Confidence: 1.0,
		},
	}
	neutralExpected := 0.0065
	neutralActual, warning := GetWaistTilt(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetWaistTilt(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetWaistTilt(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, neutralActual, neutralExpected)
	}
	// right tilt -> raise left hip
	keypoints.LHip = &skp.Keypoint{
		X:          564.604,
		Y:          1341.614,
		Confidence: 1.0,
	}
	rightTiltExpected := 44.417
	rightTiltActual, warning := GetWaistTilt(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetWaistTilt(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if rightTiltActual <= neutralActual {
		t.Errorf("GetWaistTilt(%+v, %+v) right tilt result %f is supposed to be greater than the neutral result %f", keypoints, calibrationInfoFaceOn, rightTiltActual, neutralActual)
	}
	if math.Abs(rightTiltActual-rightTiltExpected) > 0.01 {
		t.Errorf("GetWaistTilt(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, rightTiltActual, rightTiltExpected)
	}
	// left tilt -> lower left hip
	keypoints.LHip = &skp.Keypoint{
		X:          564.604,
		Y:          1563.614,
		Confidence: 1.0,
	}
	leftTiltExpected := -50.078
	leftTiltActual, warning := GetWaistTilt(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetWaistTilt(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if leftTiltActual >= neutralActual {
		t.Errorf("GetWaistTilt(%+v, %+v) left tilt result %f is supposed to be less than the neutral result %f", keypoints, calibrationInfoFaceOn, leftTiltActual, neutralActual)
	}
	if math.Abs(leftTiltActual-leftTiltExpected) > 0.01 {
		t.Errorf("GetWaistTilt(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, leftTiltActual, leftTiltExpected)
	}
}

func TestGetShaftLean(t *testing.T) {
	// neutral shaft lean
	centralClubButt := calibrationInfoFaceOn.ClubButtPoint
	centralExpected := 2.532
	centralActual, warning := GetShaftLean(calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetShaftLean(%+v) has an unexpected warning: %v", calibrationInfoFaceOn, warning)
	}
	if math.Abs(centralActual-centralExpected) > 0.01 {
		t.Errorf("GetShaftLean(%+v) = %f; expected %f", calibrationInfoFaceOn, centralActual, centralExpected)
	}
	// forward -> shift club butt forward
	calibrationInfoFaceOn.ClubButtPoint = skp.Keypoint{
		X:          632,
		Y:          1480,
		Confidence: 1.0,
	}
	forwardExpected := 15.128
	forwardActual, warning := GetShaftLean(calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetShaftLean(%+v) has an unexpected warning: %v", calibrationInfoFaceOn, warning)
	}
	if forwardActual <= centralActual {
		t.Errorf("GetShaftLean(%+v) forward result %f is supposed to be greater than the central result %f", calibrationInfoFaceOn, forwardActual, centralActual)
	}
	if math.Abs(forwardActual-forwardExpected) > 0.01 {
		t.Errorf("GetShaftLean(%+v) = %f; expected %f", calibrationInfoFaceOn, forwardActual, forwardExpected)
	}
	// backward -> shift club butt backwards
	calibrationInfoFaceOn.ClubButtPoint = skp.Keypoint{
		X:          422,
		Y:          1480,
		Confidence: 1.0,
	}
	backwardExpected := -10.547
	backwardActual, warning := GetShaftLean(calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetShaftLean(%+v) has an unexpected warning: %v", calibrationInfoFaceOn, warning)
	}
	if backwardActual >= centralActual {
		t.Errorf("GetShaftLean(%+v) backward result %f is supposed to be less than the central result %f", calibrationInfoFaceOn, backwardActual, centralActual)
	}
	if math.Abs(backwardActual-backwardExpected) > 0.01 {
		t.Errorf("GetShaftLean(%+v) = %f; expected %f", calibrationInfoFaceOn, backwardActual, backwardExpected)
	}
	// reset club head point
	calibrationInfoFaceOn.ClubButtPoint = centralClubButt
}

func TestGetBallPosition(t *testing.T) {
	// central ball position
	centralGolfBallPoint := calibrationInfoFaceOn.GolfBallPoint
	keypoints := &skp.Body25PoseKeypoints{
		LHeel: &skp.Keypoint{
			X:          605.136,
			Y:          1760.744,
			Confidence: 1.0,
		},
		RHeel: &skp.Keypoint{
			X:          415.120,
			Y:          1767.229,
			Confidence: 1.0,
		},
	}
	centralExpected := 3.715
	centralActual, warning := GetBallPosition(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetBallPosition(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(centralActual-centralExpected) > 0.01 {
		t.Errorf("GetBallPosition(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, centralActual, centralExpected)
	}
	// forward -> shift golf ball forward
	calibrationInfoFaceOn.GolfBallPoint = skp.Keypoint{
		X:          632,
		Y:          1944,
		Confidence: 1.0,
	}
	forwardExpected := 32.143
	forwardActual, warning := GetBallPosition(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetBallPosition(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if forwardActual <= centralActual {
		t.Errorf("GetBallPosition(%+v, %+v) forward result %f is supposed to be greater than the central result %f", keypoints, calibrationInfoFaceOn, forwardActual, centralActual)
	}
	if math.Abs(forwardActual-forwardExpected) > 0.01 {
		t.Errorf("GetBallPosition(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, forwardActual, forwardExpected)
	}
	// back -> shift golf ball back
	calibrationInfoFaceOn.GolfBallPoint = skp.Keypoint{
		X:          445,
		Y:          1944,
		Confidence: 1.0,
	}
	backExpected := -21.844
	backActual, warning := GetBallPosition(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetBallPosition(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if backActual >= centralActual {
		t.Errorf("GetBallPosition(%+v, %+v) back result %f is supposed to be less than the central result %f", keypoints, calibrationInfoFaceOn, backActual, centralActual)
	}
	if math.Abs(backActual-backExpected) > 0.01 {
		t.Errorf("GetBallPosition(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, backActual, backExpected)
	}
	// reset golf ball point
	calibrationInfoFaceOn.GolfBallPoint = centralGolfBallPoint
}

func TestGetHeadPosition(t *testing.T) {
	// central head position
	keypoints := &skp.Body25PoseKeypoints{
		Nose: &skp.Keypoint{
			X:          510.237,
			Y:          1155.832,
			Confidence: 1.0,
		},
		LHeel: &skp.Keypoint{
			X:          605.136,
			Y:          1760.744,
			Confidence: 1.0,
		},
		RHeel: &skp.Keypoint{
			X:          415.120,
			Y:          1767.229,
			Confidence: 1.0,
		},
	}
	centralExpected := 1.964
	centralActual, warning := GetHeadPosition(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetHeadPosition(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(centralActual-centralExpected) > 0.01 {
		t.Errorf("GetHeadPosition(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, centralActual, centralExpected)
	}
	// forward -> shift nose forward
	keypoints.Nose = &skp.Keypoint{
		X:          603.237,
		Y:          1155.832,
		Confidence: 1.0,
	}
	forwardExpected := 10.659
	forwardActual, warning := GetHeadPosition(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetHeadPosition(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if forwardActual <= centralActual {
		t.Errorf("GetHeadPosition(%+v, %+v) forward result %f is supposed to be greater than the central result %f", keypoints, calibrationInfoFaceOn, forwardActual, centralActual)
	}
	if math.Abs(forwardActual-forwardExpected) > 0.01 {
		t.Errorf("GetHeadPosition(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, forwardActual, forwardExpected)
	}
	// back -> shift nose backward
	keypoints.Nose = &skp.Keypoint{
		X:          412.237,
		Y:          1155.832,
		Confidence: 1.0,
	}
	backExpected := -7.189
	backActual, warning := GetHeadPosition(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetHeadPosition(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if backActual >= centralActual {
		t.Errorf("GetHeadPosition(%+v, %+v) back result %f is supposed to be less than the central result %f", keypoints, calibrationInfoFaceOn, backActual, centralActual)
	}
	if math.Abs(backActual-backExpected) > 0.01 {
		t.Errorf("GetHeadPosition(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, backActual, backExpected)
	}
}

func TestGetChestPosition(t *testing.T) {
	// central chest position
	keypoints := &skp.Body25PoseKeypoints{
		Neck: &skp.Keypoint{
			X:          510.222,
			Y:          1216.812,
			Confidence: 1.0,
		},
		LHeel: &skp.Keypoint{
			X:          605.136,
			Y:          1760.744,
			Confidence: 1.0,
		},
		RHeel: &skp.Keypoint{
			X:          415.120,
			Y:          1767.229,
			Confidence: 1.0,
		},
	}
	centralExpected := 1.964
	centralActual, warning := GetChestPosition(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetChestPosition(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(centralActual-centralExpected) > 0.01 {
		t.Errorf("GetChestPosition(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, centralActual, centralExpected)
	}
	// forward -> shift neck forward
	keypoints.Neck = &skp.Keypoint{
		X:          634.222,
		Y:          1216.812,
		Confidence: 1.0,
	}
	forwardExpected := 14.732
	forwardActual, warning := GetChestPosition(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetChestPosition(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if forwardActual <= centralActual {
		t.Errorf("GetChestPosition(%+v, %+v) forward result %f is supposed to be greater than the central result %f", keypoints, calibrationInfoFaceOn, forwardActual, centralActual)
	}
	if math.Abs(forwardActual-forwardExpected) > 0.01 {
		t.Errorf("GetChestPosition(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, forwardActual, forwardExpected)
	}
	// back -> shift neck backward
	keypoints.Neck = &skp.Keypoint{
		X:          422.222,
		Y:          1216.812,
		Confidence: 1.0,
	}
	backExpected := -7.172
	backActual, warning := GetChestPosition(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetChestPosition(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if backActual >= centralActual {
		t.Errorf("GetChestPosition(%+v, %+v) back result %f is supposed to be less than the central result %f", keypoints, calibrationInfoFaceOn, backActual, centralActual)
	}
	if math.Abs(backActual-backExpected) > 0.01 {
		t.Errorf("GetChestPosition(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, backActual, backExpected)
	}
}

func TestGetMidhipPosition(t *testing.T) {
	// central midhip position
	keypoints := &skp.Body25PoseKeypoints{
		Midhip: &skp.Keypoint{
			X:          510.483,
			Y:          1441.562,
			Confidence: 1.0,
		},
		LHeel: &skp.Keypoint{
			X:          605.136,
			Y:          1760.744,
			Confidence: 1.0,
		},
		RHeel: &skp.Keypoint{
			X:          415.120,
			Y:          1767.229,
			Confidence: 1.0,
		},
	}
	centralExpected := 2.017
	centralActual, warning := GetMidhipPosition(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetMidhipPosition(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if math.Abs(centralActual-centralExpected) > 0.01 {
		t.Errorf("GetMidhipPosition(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, centralActual, centralExpected)
	}
	// forward -> shift midhip forward
	keypoints.Midhip = &skp.Keypoint{
		X:          632.483,
		Y:          1441.562,
		Confidence: 1.0,
	}
	forwardExpected := 22.735
	forwardActual, warning := GetMidhipPosition(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetMidhipPosition(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if forwardActual <= centralActual {
		t.Errorf("GetMidhipPosition(%+v, %+v) forward result %f is supposed to be greater than the central result %f", keypoints, calibrationInfoFaceOn, forwardActual, centralActual)
	}
	if math.Abs(forwardActual-forwardExpected) > 0.01 {
		t.Errorf("GetMidhipPosition(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, forwardActual, forwardExpected)
	}
	// back -> shift midhip backward
	keypoints.Midhip = &skp.Keypoint{
		X:          413.483,
		Y:          1441.562,
		Confidence: 1.0,
	}
	backExpected := -14.731
	backActual, warning := GetMidhipPosition(keypoints, calibrationInfoFaceOn)
	if warning != nil {
		t.Errorf("GetMidhipPosition(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfoFaceOn, warning)
	}
	if backActual >= centralActual {
		t.Errorf("GetMidhipPosition(%+v, %+v) back result %f is supposed to be less than the central result %f", keypoints, calibrationInfoFaceOn, backActual, centralActual)
	}
	if math.Abs(backActual-backExpected) > 0.01 {
		t.Errorf("GetMidhipPosition(%+v, %+v) = %f; expected %f", keypoints, calibrationInfoFaceOn, backActual, backExpected)
	}
}
