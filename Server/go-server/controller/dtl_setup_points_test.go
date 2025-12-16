package controller

import (
	"math"
	"testing"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

var calibrationInfo *util.CalibrationInfo = &util.CalibrationInfo{
	CalibrationType: skp.CalibrationType_FULL_CALIBRATION,
	FeetLineMethod:  skp.FeetLineMethod_USE_HEEL_LINE,
	HorAxisLine: util.Line{
		Slope:      0.0324,
		YIntercept: 1679.177,
		PointOnLine: util.Point{
			XPos: 625.811,
			YPos: 1699.506,
		},
	},
	VertAxisLine: util.Line{
		Slope:      -15699.017,
		YIntercept: 8123993.528,
		PointOnLine: util.Point{
			XPos: 517.395,
			YPos: 1400.628,
		},
	},
	VanishingPoint: util.Point{
		XPos: 517.404,
		YPos: 1264.629,
	},
	GolfBallPoint: skp.Keypoint{
		X:          648,
		Y:          1736,
		Confidence: 1.0,
	},
	ClubButtPoint: skp.Keypoint{
		X:          412,
		Y:          1408,
		Confidence: 1.0,
	},
	ClubHeadPoint: skp.Keypoint{
		X:          628,
		Y:          1732,
		Confidence: 1.0,
	},
	ShoulderTilt: skp.Double{
		Data:    10.0,
		Warning: "",
	},
}

func TestGetSpineAngle(t *testing.T) {
	// neutral spine
	keypoints := &skp.Body25PoseKeypoints{
		Midhip: &skp.Keypoint{
			X:          299.682,
			Y:          1400.564,
			Confidence: 1.0,
		},
		Neck: &skp.Keypoint{
			X:          401.453,
			Y:          1196.713,
			Confidence: 1.0,
		},
	}
	neutralExpected := 26.526
	neutralActual, warning := GetSpineAngle(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetSpineAngle(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetSpineAngle(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, neutralActual, neutralExpected)
	}
	// bent over
	keypoints.Neck = &skp.Keypoint{
		X:          420.345,
		Y:          1295.637,
		Confidence: 1.0,
	}
	bentOverExpected := 48.986
	bentOverActual, warning := GetSpineAngle(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetSpineAngle(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if bentOverActual <= neutralActual {
		t.Errorf("GetSpineAngle(%+v, %+v) bent over result %f is supposed to be greater than the neutral result %f", keypoints, calibrationInfo, bentOverActual, neutralActual)
	}
	if math.Abs(bentOverActual-bentOverExpected) > 0.01 {
		t.Errorf("GetSpineAngle(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, bentOverActual, bentOverExpected)
	}
	// upright
	keypoints.Neck = &skp.Keypoint{
		X:          320.455,
		Y:          1096.342,
		Confidence: 1.0,
	}
	uprightExpected := 3.902
	uprightActual, warning := GetSpineAngle(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetSpineAngle(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if uprightActual >= neutralActual {
		t.Errorf("GetSpineAngle(%+v, %+v) upright result %f is supposed to be less than the neutral result %f", keypoints, calibrationInfo, uprightActual, neutralActual)
	}
	if math.Abs(uprightActual-uprightExpected) > 0.01 {
		t.Errorf("GetSpineAngle(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, uprightActual, uprightExpected)
	}
}

func TestFeetAlignment(t *testing.T) {
	keypoints := &skp.Body25PoseKeypoints{
		LHeel: &skp.Keypoint{
			X:          326.782,
			Y:          1706.318,
			Confidence: 1.0,
		},
		RHeel: &skp.Keypoint{
			X:          272.357,
			Y:          1815.146,
			Confidence: 1.0,
		},
		LBigToe: &skp.Keypoint{
			X:          401.705,
			Y:          1699.557,
			Confidence: 1.0,
		},
		RBigToe: &skp.Keypoint{
			X:          381.126,
			Y:          1814.789,
			Confidence: 1.0,
		},
	}
	// heel feet line method
	heelExpected := 2.574
	heelActual, warning := GetFeetAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetFeetAlignment(%+v, %+v) heel has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if math.Abs(heelActual-heelExpected) > 0.01 {
		t.Errorf("GetFeetAlignment(%+v, %+v) = %f; heel expected %f", keypoints, calibrationInfo, heelActual, heelExpected)
	}
	// toe feet line method
	calibrationInfo.FeetLineMethod = skp.FeetLineMethod_USE_TOE_LINE
	toeExpected := -3.787
	toeActual, warning := GetFeetAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetFeetAlignment(%+v, %+v) toe has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if math.Abs(toeActual-toeExpected) > 0.01 {
		t.Errorf("GetFeetAlignment(%+v, %+v) = %f; toe expected %f", keypoints, calibrationInfo, toeActual, toeExpected)
	}
}

func TestHeelAlignment(t *testing.T) {
	// neutral alignment
	keypoints := &skp.Body25PoseKeypoints{
		LHeel: &skp.Keypoint{
			X:          326.782,
			Y:          1706.318,
			Confidence: 1.0,
		},
		RHeel: &skp.Keypoint{
			X:          272.357,
			Y:          1815.146,
			Confidence: 1.0,
		},
	}
	neutralExpected := 2.574
	neutralActual, warning := GetHeelAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetHeelAlignment(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetHeelAlignment(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, neutralActual, neutralExpected)
	}
	// open alignment
	keypoints.LHeel = &skp.Keypoint{
		X:          282.782,
		Y:          1706.318,
		Confidence: 1.0,
	}
	openExpected := -18.523
	openActual, warning := GetHeelAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetHeelAlignment(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if openActual >= neutralActual {
		t.Errorf("GetHeelAlignment(%+v, %+v) open result %f is supposed to be less than the neutral result %f", keypoints, calibrationInfo, openActual, neutralActual)
	}
	if openActual > 0 {
		t.Errorf("GetHeelAlignment(%+v, %+v) open result %f is supposed to be less than 0", keypoints, calibrationInfo, openActual)
	}
	if math.Abs(openActual-openExpected) > 0.01 {
		t.Errorf("GetHeelAlignment(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, openActual, openExpected)
	}
	// closed alignment
	keypoints.LHeel = &skp.Keypoint{
		X:          374.782,
		Y:          1706.318,
		Confidence: 1.0,
	}
	closedExpected := 19.269
	closedActual, warning := GetHeelAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetHeelAlignment(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if closedActual <= neutralActual {
		t.Errorf("GetHeelAlignment(%+v, %+v) closed result %f is supposed to be greater than the neutral result %f", keypoints, calibrationInfo, closedActual, neutralActual)
	}
	if closedActual < 0 {
		t.Errorf("GetHeelAlignment(%+v, %+v) open result %f is supposed to be greater than 0", keypoints, calibrationInfo, closedActual)
	}
	if math.Abs(closedActual-closedExpected) > 0.01 {
		t.Errorf("GetHeelAlignment(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, closedActual, closedExpected)
	}
}

func TestToeAlignment(t *testing.T) {
	// neutral alignment
	keypoints := &skp.Body25PoseKeypoints{
		LBigToe: &skp.Keypoint{
			X:          401.705,
			Y:          1699.557,
			Confidence: 1.0,
		},
		RBigToe: &skp.Keypoint{
			X:          381.126,
			Y:          1814.789,
			Confidence: 1.0,
		},
	}
	neutralExpected := -3.787
	neutralActual, warning := GetToeAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetToeAlignment(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetToeAlignment(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, neutralActual, neutralExpected)
	}
	// open alignment
	keypoints.LBigToe = &skp.Keypoint{
		X:          302.705,
		Y:          1699.557,
		Confidence: 1.0,
	}
	openExpected := -48.149
	openActual, warning := GetToeAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetToeAlignment(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if openActual >= neutralActual {
		t.Errorf("GetToeAlignment(%+v, %+v) open result %f is supposed to be less than the neutral result %f", keypoints, calibrationInfo, openActual, neutralActual)
	}
	if openActual > 0 {
		t.Errorf("GetToeAlignment(%+v, %+v) open result %f is supposed to be less than 0", keypoints, calibrationInfo, openActual)
	}
	if math.Abs(openActual-openExpected) > 0.01 {
		t.Errorf("GetToeAlignment(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, openActual, openExpected)
	}
	// closed alignment
	keypoints.LBigToe = &skp.Keypoint{
		X:          456.705,
		Y:          1699.557,
		Confidence: 1.0,
	}
	closedExpected := 19.347
	closedActual, warning := GetToeAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetToeAlignment(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if closedActual <= neutralActual {
		t.Errorf("GetToeAlignment(%+v, %+v) closed result %f is supposed to be greater than the neutral result %f", keypoints, calibrationInfo, closedActual, neutralActual)
	}
	if closedActual < 0 {
		t.Errorf("GetToeAlignment(%+v, %+v) open result %f is supposed to be greater than 0", keypoints, calibrationInfo, closedActual)
	}
	if math.Abs(closedActual-closedExpected) > 0.01 {
		t.Errorf("GetToeAlignment(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, closedActual, closedExpected)
	}
}

func TestShoulderAlignment(t *testing.T) {
	// neutral alignment
	keypoints := &skp.Body25PoseKeypoints{
		LShoulder: &skp.Keypoint{
			X:          440.589,
			Y:          1211.155,
			Confidence: 1.0,
		},
		RShoulder: &skp.Keypoint{
			X:          401.740,
			Y:          1217.086,
			Confidence: 1.0,
		},
	}
	neutralExpected := 5.520
	neutralActual, warning := GetShoulderAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetShoulderAlignment(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetShoulderAlignment(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, neutralActual, neutralExpected)
	}
	// open alignment
	keypoints.LShoulder = &skp.Keypoint{
		X:          395.456,
		Y:          1205.086,
		Confidence: 1.0,
	}
	openExpected := -103.438
	openActual, warning := GetShoulderAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetShoulderAlignment(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if openActual >= neutralActual {
		t.Errorf("GetShoulderAlignment(%+v, %+v) open result %f is supposed to be less than the neutral result %f", keypoints, calibrationInfo, openActual, neutralActual)
	}
	if openActual > 0 {
		t.Errorf("GetShoulderAlignment(%+v, %+v) open result %f is supposed to be less than 0", keypoints, calibrationInfo, openActual)
	}
	if math.Abs(openActual-openExpected) > 0.01 {
		t.Errorf("GetShoulderAlignment(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, openActual, openExpected)
	}
	// closed alignment
	keypoints.LShoulder = &skp.Keypoint{
		X:          440.589,
		Y:          1240.562,
		Confidence: 1.0,
	}
	closedExpected := 45.345
	closedActual, warning := GetShoulderAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetShoulderAlignment(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if closedActual <= neutralActual {
		t.Errorf("GetShoulderAlignment(%+v, %+v) closed result %f is supposed to be greater than the neutral result %f", keypoints, calibrationInfo, closedActual, neutralActual)
	}
	if closedActual < 0 {
		t.Errorf("GetShoulderAlignment(%+v, %+v) open result %f is supposed to be greater than 0", keypoints, calibrationInfo, closedActual)
	}
	if math.Abs(closedActual-closedExpected) > 0.01 {
		t.Errorf("GetShoulderAlignment(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, closedActual, closedExpected)
	}
}

func TestWaistAlignment(t *testing.T) {
	// neutral alignment
	keypoints := &skp.Body25PoseKeypoints{
		LHip: &skp.Keypoint{
			X:          306.232,
			Y:          1393.867,
			Confidence: 1.0,
		},
		RHip: &skp.Keypoint{
			X:          275.587,
			Y:          1414.083,
			Confidence: 1.0,
		},
	}
	neutralExpected := -1.694
	neutralActual, warning := GetWaistAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetWaistAlignment(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetWaistAlignment(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, neutralActual, neutralExpected)
	}
	// open alignment
	keypoints.LHip = &skp.Keypoint{
		X:          276.232,
		Y:          1393.867,
		Confidence: 1.0,
	}
	openExpected := -56.454
	openActual, warning := GetWaistAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetWaistAlignment(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if openActual >= neutralActual {
		t.Errorf("GetWaistAlignment(%+v, %+v) open result %f is supposed to be less than the neutral result %f", keypoints, calibrationInfo, openActual, neutralActual)
	}
	if openActual > 0 {
		t.Errorf("GetWaistAlignment(%+v, %+v) open result %f is supposed to be less than 0", keypoints, calibrationInfo, openActual)
	}
	if math.Abs(openActual-openExpected) > 0.01 {
		t.Errorf("GetWaistAlignment(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, openActual, openExpected)
	}
	// closed alignment
	keypoints.LHip = &skp.Keypoint{
		X:          346.232,
		Y:          1393.867,
		Confidence: 1.0,
	}
	closedExpected := 15.748
	closedActual, warning := GetWaistAlignment(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetWaistAlignment(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if closedActual <= neutralActual {
		t.Errorf("GetWaistAlignment(%+v, %+v) closed result %f is supposed to be greater than the neutral result %f", keypoints, calibrationInfo, closedActual, neutralActual)
	}
	if closedActual < 0 {
		t.Errorf("GetWaistAlignment(%+v, %+v) open result %f is supposed to be greater than 0", keypoints, calibrationInfo, closedActual)
	}
	if math.Abs(closedActual-closedExpected) > 0.01 {
		t.Errorf("GetWaistAlignment(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, closedActual, closedExpected)
	}
}

func TestKneeBend(t *testing.T) {
	// standard knee bend
	keypoints := &skp.Body25PoseKeypoints{
		RHip: &skp.Keypoint{
			X:          299.587,
			Y:          1414.083,
			Confidence: 1.0,
		},
		RKnee: &skp.Keypoint{
			X:          306.190,
			Y:          1597.461,
			Confidence: 1.0,
		},
		RAnkle: &skp.Keypoint{
			X:          292.723,
			Y:          1794.595,
			Confidence: 1.0,
		},
	}
	standardExpected := 5.969
	standardActual, warning := GetKneeBend(keypoints)
	if warning != nil {
		t.Errorf("GetKneeBend(%+v) has an unexpected warning: %v", keypoints, warning)
	}
	if math.Abs(standardActual-standardExpected) > 0.01 {
		t.Errorf("GetKneeBend(%+v) = %f; expected %f", keypoints, standardActual, standardExpected)
	}
	// big knee bend
	keypoints.RKnee = &skp.Keypoint{
		X:          394.567,
		Y:          1597.461,
		Confidence: 1.0,
	}
	bigBendExpected := 54.703
	bigBendActual, warning := GetKneeBend(keypoints)
	if warning != nil {
		t.Errorf("GetKneeBend(%+v) has an unexpected warning: %v", keypoints, warning)
	}
	if bigBendActual <= standardActual {
		t.Errorf("GetKneeBend(%+v) big bend result %f is supposed to be greater than the standard result %f", keypoints, bigBendActual, standardActual)
	}
	if math.Abs(bigBendActual-bigBendExpected) > 0.01 {
		t.Errorf("GetKneeBend(%+v) = %f; expected %f", keypoints, bigBendActual, bigBendExpected)
	}
	// small knee bend
	keypoints.RKnee = &skp.Keypoint{
		X:          299.987,
		Y:          1597.461,
		Confidence: 1.0,
	}
	smallBendExpected := 2.235
	smallBendActual, warning := GetKneeBend(keypoints)
	if warning != nil {
		t.Errorf("GetKneeBend(%+v) has an unexpected warning: %v", keypoints, warning)
	}
	if bigBendActual <= standardActual {
		t.Errorf("GetKneeBend(%+v) small bend result %f is supposed to be less than the standard result %f", keypoints, smallBendActual, standardActual)
	}
	if math.Abs(smallBendActual-smallBendExpected) > 0.01 {
		t.Errorf("GetKneeBend(%+v) = %f; expected %f", keypoints, smallBendActual, smallBendExpected)
	}
}

func TestDistanceFromBall(t *testing.T) {
	// standard distance
	keypoints := &skp.Body25PoseKeypoints{
		Midhip: &skp.Keypoint{
			X:          299.682,
			Y:          1400.564,
			Confidence: 1.0,
		},
		Neck: &skp.Keypoint{
			X:          401.453,
			Y:          1196.713,
			Confidence: 1.0,
		},
		LBigToe: &skp.Keypoint{
			X:          401.705,
			Y:          1699.557,
			Confidence: 1.0,
		},
		RBigToe: &skp.Keypoint{
			X:          381.126,
			Y:          1814.789,
			Confidence: 1.0,
		},
	}
	standardExpected := 1.092
	standardActual, warning := GetDistanceFromBall(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetDistanceFromBall(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if math.Abs(standardActual-standardExpected) > 0.01 {
		t.Errorf("GetDistanceFromBall(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, standardActual, standardExpected)
	}
	// far from ball
	keypoints.LBigToe = &skp.Keypoint{
		X:          301.705,
		Y:          1699.557,
		Confidence: 1.0,
	}
	keypoints.RBigToe = &skp.Keypoint{
		X:          281.126,
		Y:          1814.789,
		Confidence: 1.0,
	}
	farExpected := 1.524
	farActual, warning := GetDistanceFromBall(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetDistanceFromBall(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if farActual <= standardActual {
		t.Errorf("GetDistanceFromBall(%+v, %+v) far result %f is supposed to be greater than the standard result %f", keypoints, calibrationInfo, farActual, standardActual)
	}
	if math.Abs(farActual-farExpected) > 0.01 {
		t.Errorf("GetDistanceFromBall(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, farActual, farExpected)
	}
	// close to ball
	keypoints.LBigToe = &skp.Keypoint{
		X:          501.705,
		Y:          1699.557,
		Confidence: 1.0,
	}
	keypoints.RBigToe = &skp.Keypoint{
		X:          481.126,
		Y:          1814.789,
		Confidence: 1.0,
	}
	closeExpected := 0.660
	closeActual, warning := GetDistanceFromBall(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetDistanceFromBall(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if closeActual >= standardActual {
		t.Errorf("GetDistanceFromBall(%+v, %+v) close result %f is supposed to be less than the standard result %f", keypoints, calibrationInfo, closeActual, standardActual)
	}
	if math.Abs(closeActual-closeExpected) > 0.01 {
		t.Errorf("GetDistanceFromBall(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, closeActual, closeExpected)
	}
}

func TestGetUlnarDeviation(t *testing.T) {
	// standard ulnar deviation
	keypoints := &skp.Body25PoseKeypoints{
		RElbow: &skp.Keypoint{
			X:          415.229,
			Y:          1346.263,
			Confidence: 1.0,
		},
		RWrist: &skp.Keypoint{
			X:          442.338,
			Y:          1434.486,
			Confidence: 1.0,
		},
	}
	standardExpected := 165.115
	standardActual, warning := GetUlnarDeviation(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetUlnarDeviation(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if math.Abs(standardActual-standardExpected) > 0.01 {
		t.Errorf("GetUlnarDeviation(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, standardActual, standardExpected)
	}
	// high hands
	keypoints.RWrist = &skp.Keypoint{
		X:          442.338,
		Y:          1420.001,
		Confidence: 1.0,
	}
	highExpected := 169.429
	highActual, warning := GetUlnarDeviation(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetUlnarDeviation(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if highActual <= standardActual {
		t.Errorf("GetUlnarDeviation(%+v, %+v) high hands result %f is supposed to be greater than the standard result %f", keypoints, calibrationInfo, highActual, standardActual)
	}
	if math.Abs(highActual-highExpected) > 0.01 {
		t.Errorf("GetUlnarDeviation(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, highActual, highExpected)
	}
	// low hands
	keypoints.RWrist = &skp.Keypoint{
		X:          442.338,
		Y:          1603.782,
		Confidence: 1.0,
	}
	lowExpected := 130.638
	lowActual, warning := GetUlnarDeviation(keypoints, calibrationInfo)
	if warning != nil {
		t.Errorf("GetUlnarDeviation(%+v, %+v) has an unexpected warning: %v", keypoints, calibrationInfo, warning)
	}
	if highActual <= standardActual {
		t.Errorf("GetUlnarDeviation(%+v, %+v) low hands result %f is supposed to be less than the standard result %f", keypoints, calibrationInfo, lowActual, standardActual)
	}
	if math.Abs(lowActual-lowExpected) > 0.01 {
		t.Errorf("GetUlnarDeviation(%+v, %+v) = %f; expected %f", keypoints, calibrationInfo, lowActual, lowExpected)
	}
}
