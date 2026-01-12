package controller

import (
	"math"
	"testing"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

var golfSpecificDatapointsDTL *skp.GolfSpecificDatapoints = &skp.GolfSpecificDatapoints{
	GolfBall: &skp.Datapoint{
		X:          648,
		Y:          1736,
		Confidence: 1.0,
	},
	ClubButt: &skp.Datapoint{
		X:          412,
		Y:          1408,
		Confidence: 1.0,
	},
	ClubHead: &skp.Datapoint{
		X:          628,
		Y:          1732,
		Confidence: 1.0,
	},
	ShoulderTilt: &skp.Double{
		Data:    10.0,
		Warning: "",
	},
}

var calibrationInfoDTL *util.CalibrationInfo = &util.CalibrationInfo{
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
}

func TestGetSpineAngle(t *testing.T) {
	// neutral spine
	bodyDatapoints := &skp.Body25PoseDatapoints{
		Midhip: &skp.Datapoint{
			X:          299.682,
			Y:          1400.564,
			Confidence: 1.0,
		},
		Neck: &skp.Datapoint{
			X:          401.453,
			Y:          1196.713,
			Confidence: 1.0,
		},
	}
	neutralExpected := 26.526
	neutralActual, warning := GetSpineAngle(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetSpineAngle(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetSpineAngle(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoDTL, neutralActual, neutralExpected)
	}
	// bent over -> shift neck forward
	bodyDatapoints.Neck = &skp.Datapoint{
		X:          420.345,
		Y:          1295.637,
		Confidence: 1.0,
	}
	bentOverExpected := 48.986
	bentOverActual, warning := GetSpineAngle(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetSpineAngle(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if bentOverActual <= neutralActual {
		t.Errorf("GetSpineAngle(%+v, %+v) bent over result %f is supposed to be greater than the neutral result %f", bodyDatapoints, calibrationInfoDTL, bentOverActual, neutralActual)
	}
	if math.Abs(bentOverActual-bentOverExpected) > 0.01 {
		t.Errorf("GetSpineAngle(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoDTL, bentOverActual, bentOverExpected)
	}
	// upright -> shift neck more in line with midhip
	bodyDatapoints.Neck = &skp.Datapoint{
		X:          320.455,
		Y:          1096.342,
		Confidence: 1.0,
	}
	uprightExpected := 3.902
	uprightActual, warning := GetSpineAngle(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetSpineAngle(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if uprightActual >= neutralActual {
		t.Errorf("GetSpineAngle(%+v, %+v) upright result %f is supposed to be less than the neutral result %f", bodyDatapoints, calibrationInfoDTL, uprightActual, neutralActual)
	}
	if math.Abs(uprightActual-uprightExpected) > 0.01 {
		t.Errorf("GetSpineAngle(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoDTL, uprightActual, uprightExpected)
	}
}

func TestFeetAlignment(t *testing.T) {
	bodyDatapoints := &skp.Body25PoseDatapoints{
		LHeel: &skp.Datapoint{
			X:          326.782,
			Y:          1706.318,
			Confidence: 1.0,
		},
		RHeel: &skp.Datapoint{
			X:          272.357,
			Y:          1815.146,
			Confidence: 1.0,
		},
		LBigToe: &skp.Datapoint{
			X:          401.705,
			Y:          1699.557,
			Confidence: 1.0,
		},
		RBigToe: &skp.Datapoint{
			X:          381.126,
			Y:          1814.789,
			Confidence: 1.0,
		},
	}
	// heel feet line method
	heelExpected := 2.574
	heelActual, warning := GetFeetAlignment(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetFeetAlignment(%+v, %+v) heel has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if math.Abs(heelActual-heelExpected) > 0.01 {
		t.Errorf("GetFeetAlignment(%+v, %+v) = %f; heel expected %f", bodyDatapoints, calibrationInfoDTL, heelActual, heelExpected)
	}
	// toe feet line method
	calibrationInfoDTL.FeetLineMethod = skp.FeetLineMethod_USE_TOE_LINE
	toeExpected := -3.787
	toeActual, warning := GetFeetAlignment(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetFeetAlignment(%+v, %+v) toe has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if math.Abs(toeActual-toeExpected) > 0.01 {
		t.Errorf("GetFeetAlignment(%+v, %+v) = %f; toe expected %f", bodyDatapoints, calibrationInfoDTL, toeActual, toeExpected)
	}
}

func TestHeelAlignment(t *testing.T) {
	// neutral alignment
	bodyDatapoints := &skp.Body25PoseDatapoints{
		LHeel: &skp.Datapoint{
			X:          326.782,
			Y:          1706.318,
			Confidence: 1.0,
		},
		RHeel: &skp.Datapoint{
			X:          272.357,
			Y:          1815.146,
			Confidence: 1.0,
		},
	}
	neutralExpected := 2.574
	neutralActual, warning := GetHeelAlignment(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetHeelAlignment(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetHeelAlignment(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoDTL, neutralActual, neutralExpected)
	}
	// open alignment -> shift left heel open
	bodyDatapoints.LHeel = &skp.Datapoint{
		X:          282.782,
		Y:          1706.318,
		Confidence: 1.0,
	}
	openExpected := -18.523
	openActual, warning := GetHeelAlignment(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetHeelAlignment(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if openActual >= neutralActual {
		t.Errorf("GetHeelAlignment(%+v, %+v) open result %f is supposed to be less than the neutral result %f", bodyDatapoints, calibrationInfoDTL, openActual, neutralActual)
	}
	if openActual > 0 {
		t.Errorf("GetHeelAlignment(%+v, %+v) open result %f is supposed to be less than 0", bodyDatapoints, calibrationInfoDTL, openActual)
	}
	if math.Abs(openActual-openExpected) > 0.01 {
		t.Errorf("GetHeelAlignment(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoDTL, openActual, openExpected)
	}
	// closed alignment -> shift left heel closed
	bodyDatapoints.LHeel = &skp.Datapoint{
		X:          374.782,
		Y:          1706.318,
		Confidence: 1.0,
	}
	closedExpected := 19.269
	closedActual, warning := GetHeelAlignment(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetHeelAlignment(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if closedActual <= neutralActual {
		t.Errorf("GetHeelAlignment(%+v, %+v) closed result %f is supposed to be greater than the neutral result %f", bodyDatapoints, calibrationInfoDTL, closedActual, neutralActual)
	}
	if closedActual < 0 {
		t.Errorf("GetHeelAlignment(%+v, %+v) open result %f is supposed to be greater than 0", bodyDatapoints, calibrationInfoDTL, closedActual)
	}
	if math.Abs(closedActual-closedExpected) > 0.01 {
		t.Errorf("GetHeelAlignment(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoDTL, closedActual, closedExpected)
	}
}

func TestToeAlignment(t *testing.T) {
	// neutral alignment
	bodyDatapoints := &skp.Body25PoseDatapoints{
		LBigToe: &skp.Datapoint{
			X:          401.705,
			Y:          1699.557,
			Confidence: 1.0,
		},
		RBigToe: &skp.Datapoint{
			X:          381.126,
			Y:          1814.789,
			Confidence: 1.0,
		},
	}
	neutralExpected := -3.787
	neutralActual, warning := GetToeAlignment(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetToeAlignment(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetToeAlignment(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoDTL, neutralActual, neutralExpected)
	}
	// open alignment -> shift left big toe open
	bodyDatapoints.LBigToe = &skp.Datapoint{
		X:          302.705,
		Y:          1699.557,
		Confidence: 1.0,
	}
	openExpected := -48.149
	openActual, warning := GetToeAlignment(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetToeAlignment(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if openActual >= neutralActual {
		t.Errorf("GetToeAlignment(%+v, %+v) open result %f is supposed to be less than the neutral result %f", bodyDatapoints, calibrationInfoDTL, openActual, neutralActual)
	}
	if openActual > 0 {
		t.Errorf("GetToeAlignment(%+v, %+v) open result %f is supposed to be less than 0", bodyDatapoints, calibrationInfoDTL, openActual)
	}
	if math.Abs(openActual-openExpected) > 0.01 {
		t.Errorf("GetToeAlignment(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoDTL, openActual, openExpected)
	}
	// closed alignment -> shift left big toe closed
	bodyDatapoints.LBigToe = &skp.Datapoint{
		X:          456.705,
		Y:          1699.557,
		Confidence: 1.0,
	}
	closedExpected := 19.347
	closedActual, warning := GetToeAlignment(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetToeAlignment(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if closedActual <= neutralActual {
		t.Errorf("GetToeAlignment(%+v, %+v) closed result %f is supposed to be greater than the neutral result %f", bodyDatapoints, calibrationInfoDTL, closedActual, neutralActual)
	}
	if closedActual < 0 {
		t.Errorf("GetToeAlignment(%+v, %+v) open result %f is supposed to be greater than 0", bodyDatapoints, calibrationInfoDTL, closedActual)
	}
	if math.Abs(closedActual-closedExpected) > 0.01 {
		t.Errorf("GetToeAlignment(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoDTL, closedActual, closedExpected)
	}
}

func TestShoulderAlignment(t *testing.T) {
	// neutral alignment
	bodyDatapoints := &skp.Body25PoseDatapoints{
		LShoulder: &skp.Datapoint{
			X:          440.589,
			Y:          1211.155,
			Confidence: 1.0,
		},
		RShoulder: &skp.Datapoint{
			X:          401.740,
			Y:          1217.086,
			Confidence: 1.0,
		},
	}
	neutralExpected := 5.520
	neutralActual, warning := GetShoulderAlignment(bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetShoulderAlignment(%+v, %+v, %+v) has an unexpected warning: %v", bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetShoulderAlignment(%+v, %+v, %+v) = %f; expected %f", bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL, neutralActual, neutralExpected)
	}
	// open alignment -> shift left shoulder open
	bodyDatapoints.LShoulder = &skp.Datapoint{
		X:          395.456,
		Y:          1205.086,
		Confidence: 1.0,
	}
	openExpected := -103.438
	openActual, warning := GetShoulderAlignment(bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetShoulderAlignment(%+v, %+v, %+v) has an unexpected warning: %v", bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL, warning)
	}
	if openActual >= neutralActual {
		t.Errorf("GetShoulderAlignment(%+v, %+v, %+v) open result %f is supposed to be less than the neutral result %f", bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL, openActual, neutralActual)
	}
	if openActual > 0 {
		t.Errorf("GetShoulderAlignment(%+v, %+v, %+v) open result %f is supposed to be less than 0", bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL, openActual)
	}
	if math.Abs(openActual-openExpected) > 0.01 {
		t.Errorf("GetShoulderAlignment(%+v, %+v, %+v) = %f; expected %f", bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL, openActual, openExpected)
	}
	// closed alignment -> shift left shoulder closed
	bodyDatapoints.LShoulder = &skp.Datapoint{
		X:          440.589,
		Y:          1240.562,
		Confidence: 1.0,
	}
	closedExpected := 45.345
	closedActual, warning := GetShoulderAlignment(bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetShoulderAlignment(%+v, %+v, %+v) has an unexpected warning: %v", bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL, warning)
	}
	if closedActual <= neutralActual {
		t.Errorf("GetShoulderAlignment(%+v, %+v, %+v) closed result %f is supposed to be greater than the neutral result %f", bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL, closedActual, neutralActual)
	}
	if closedActual < 0 {
		t.Errorf("GetShoulderAlignment(%+v, %+v, %+v) open result %f is supposed to be greater than 0", bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL, closedActual)
	}
	if math.Abs(closedActual-closedExpected) > 0.01 {
		t.Errorf("GetShoulderAlignment(%+v, %+v, %+v) = %f; expected %f", bodyDatapoints, golfSpecificDatapointsDTL, calibrationInfoDTL, closedActual, closedExpected)
	}
}

func TestWaistAlignment(t *testing.T) {
	// neutral alignment
	bodyDatapoints := &skp.Body25PoseDatapoints{
		LHip: &skp.Datapoint{
			X:          306.232,
			Y:          1393.867,
			Confidence: 1.0,
		},
		RHip: &skp.Datapoint{
			X:          275.587,
			Y:          1414.083,
			Confidence: 1.0,
		},
	}
	neutralExpected := -1.694
	neutralActual, warning := GetWaistAlignment(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetWaistAlignment(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if math.Abs(neutralActual-neutralExpected) > 0.01 {
		t.Errorf("GetWaistAlignment(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoDTL, neutralActual, neutralExpected)
	}
	// open alignment -> shift left hip open
	bodyDatapoints.LHip = &skp.Datapoint{
		X:          276.232,
		Y:          1393.867,
		Confidence: 1.0,
	}
	openExpected := -56.454
	openActual, warning := GetWaistAlignment(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetWaistAlignment(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if openActual >= neutralActual {
		t.Errorf("GetWaistAlignment(%+v, %+v) open result %f is supposed to be less than the neutral result %f", bodyDatapoints, calibrationInfoDTL, openActual, neutralActual)
	}
	if openActual > 0 {
		t.Errorf("GetWaistAlignment(%+v, %+v) open result %f is supposed to be less than 0", bodyDatapoints, calibrationInfoDTL, openActual)
	}
	if math.Abs(openActual-openExpected) > 0.01 {
		t.Errorf("GetWaistAlignment(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoDTL, openActual, openExpected)
	}
	// closed alignment -> shift left hip closed
	bodyDatapoints.LHip = &skp.Datapoint{
		X:          346.232,
		Y:          1393.867,
		Confidence: 1.0,
	}
	closedExpected := 15.748
	closedActual, warning := GetWaistAlignment(bodyDatapoints, calibrationInfoDTL)
	if warning != nil {
		t.Errorf("GetWaistAlignment(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, calibrationInfoDTL, warning)
	}
	if closedActual <= neutralActual {
		t.Errorf("GetWaistAlignment(%+v, %+v) closed result %f is supposed to be greater than the neutral result %f", bodyDatapoints, calibrationInfoDTL, closedActual, neutralActual)
	}
	if closedActual < 0 {
		t.Errorf("GetWaistAlignment(%+v, %+v) open result %f is supposed to be greater than 0", bodyDatapoints, calibrationInfoDTL, closedActual)
	}
	if math.Abs(closedActual-closedExpected) > 0.01 {
		t.Errorf("GetWaistAlignment(%+v, %+v) = %f; expected %f", bodyDatapoints, calibrationInfoDTL, closedActual, closedExpected)
	}
}

func TestKneeBend(t *testing.T) {
	// standard knee bend
	bodyDatapoints := &skp.Body25PoseDatapoints{
		RHip: &skp.Datapoint{
			X:          299.587,
			Y:          1414.083,
			Confidence: 1.0,
		},
		RKnee: &skp.Datapoint{
			X:          306.190,
			Y:          1597.461,
			Confidence: 1.0,
		},
		RAnkle: &skp.Datapoint{
			X:          292.723,
			Y:          1794.595,
			Confidence: 1.0,
		},
	}
	standardExpected := 5.969
	standardActual, warning := GetKneeBend(bodyDatapoints)
	if warning != nil {
		t.Errorf("GetKneeBend(%+v) has an unexpected warning: %v", bodyDatapoints, warning)
	}
	if math.Abs(standardActual-standardExpected) > 0.01 {
		t.Errorf("GetKneeBend(%+v) = %f; expected %f", bodyDatapoints, standardActual, standardExpected)
	}
	// big knee bend -> shift knee forward
	bodyDatapoints.RKnee = &skp.Datapoint{
		X:          394.567,
		Y:          1597.461,
		Confidence: 1.0,
	}
	bigBendExpected := 54.703
	bigBendActual, warning := GetKneeBend(bodyDatapoints)
	if warning != nil {
		t.Errorf("GetKneeBend(%+v) has an unexpected warning: %v", bodyDatapoints, warning)
	}
	if bigBendActual <= standardActual {
		t.Errorf("GetKneeBend(%+v) big bend result %f is supposed to be greater than the standard result %f", bodyDatapoints, bigBendActual, standardActual)
	}
	if math.Abs(bigBendActual-bigBendExpected) > 0.01 {
		t.Errorf("GetKneeBend(%+v) = %f; expected %f", bodyDatapoints, bigBendActual, bigBendExpected)
	}
	// small knee bend -> shift knee more in line with hip and ankle
	bodyDatapoints.RKnee = &skp.Datapoint{
		X:          299.987,
		Y:          1597.461,
		Confidence: 1.0,
	}
	smallBendExpected := 2.235
	smallBendActual, warning := GetKneeBend(bodyDatapoints)
	if warning != nil {
		t.Errorf("GetKneeBend(%+v) has an unexpected warning: %v", bodyDatapoints, warning)
	}
	if bigBendActual <= standardActual {
		t.Errorf("GetKneeBend(%+v) small bend result %f is supposed to be less than the standard result %f", bodyDatapoints, smallBendActual, standardActual)
	}
	if math.Abs(smallBendActual-smallBendExpected) > 0.01 {
		t.Errorf("GetKneeBend(%+v) = %f; expected %f", bodyDatapoints, smallBendActual, smallBendExpected)
	}
}

func TestDistanceFromBall(t *testing.T) {
	// standard distance
	bodyDatapoints := &skp.Body25PoseDatapoints{
		Midhip: &skp.Datapoint{
			X:          299.682,
			Y:          1400.564,
			Confidence: 1.0,
		},
		Neck: &skp.Datapoint{
			X:          401.453,
			Y:          1196.713,
			Confidence: 1.0,
		},
		LBigToe: &skp.Datapoint{
			X:          401.705,
			Y:          1699.557,
			Confidence: 1.0,
		},
		RBigToe: &skp.Datapoint{
			X:          381.126,
			Y:          1814.789,
			Confidence: 1.0,
		},
	}
	standardExpected := 1.092
	standardActual, warning := GetDistanceFromBall(bodyDatapoints, golfSpecificDatapointsDTL)
	if warning != nil {
		t.Errorf("GetDistanceFromBall(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, golfSpecificDatapointsDTL, warning)
	}
	if math.Abs(standardActual-standardExpected) > 0.01 {
		t.Errorf("GetDistanceFromBall(%+v, %+v) = %f; expected %f", bodyDatapoints, golfSpecificDatapointsDTL, standardActual, standardExpected)
	}
	// far from ball -> shift feet farther from ball
	bodyDatapoints.LBigToe = &skp.Datapoint{
		X:          301.705,
		Y:          1699.557,
		Confidence: 1.0,
	}
	bodyDatapoints.RBigToe = &skp.Datapoint{
		X:          281.126,
		Y:          1814.789,
		Confidence: 1.0,
	}
	farExpected := 1.524
	farActual, warning := GetDistanceFromBall(bodyDatapoints, golfSpecificDatapointsDTL)
	if warning != nil {
		t.Errorf("GetDistanceFromBall(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, golfSpecificDatapointsDTL, warning)
	}
	if farActual <= standardActual {
		t.Errorf("GetDistanceFromBall(%+v, %+v) far result %f is supposed to be greater than the standard result %f", bodyDatapoints, golfSpecificDatapointsDTL, farActual, standardActual)
	}
	if math.Abs(farActual-farExpected) > 0.01 {
		t.Errorf("GetDistanceFromBall(%+v, %+v) = %f; expected %f", bodyDatapoints, golfSpecificDatapointsDTL, farActual, farExpected)
	}
	// close to ball -> shift feet closer to ball
	bodyDatapoints.LBigToe = &skp.Datapoint{
		X:          501.705,
		Y:          1699.557,
		Confidence: 1.0,
	}
	bodyDatapoints.RBigToe = &skp.Datapoint{
		X:          481.126,
		Y:          1814.789,
		Confidence: 1.0,
	}
	closeExpected := 0.660
	closeActual, warning := GetDistanceFromBall(bodyDatapoints, golfSpecificDatapointsDTL)
	if warning != nil {
		t.Errorf("GetDistanceFromBall(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, golfSpecificDatapointsDTL, warning)
	}
	if closeActual >= standardActual {
		t.Errorf("GetDistanceFromBall(%+v, %+v) close result %f is supposed to be less than the standard result %f", bodyDatapoints, golfSpecificDatapointsDTL, closeActual, standardActual)
	}
	if math.Abs(closeActual-closeExpected) > 0.01 {
		t.Errorf("GetDistanceFromBall(%+v, %+v) = %f; expected %f", bodyDatapoints, golfSpecificDatapointsDTL, closeActual, closeExpected)
	}
}

func TestGetUlnarDeviation(t *testing.T) {
	// standard ulnar deviation
	bodyDatapoints := &skp.Body25PoseDatapoints{
		RElbow: &skp.Datapoint{
			X:          415.229,
			Y:          1346.263,
			Confidence: 1.0,
		},
		RWrist: &skp.Datapoint{
			X:          442.338,
			Y:          1434.486,
			Confidence: 1.0,
		},
	}
	standardExpected := 165.115
	standardActual, warning := GetUlnarDeviation(bodyDatapoints, golfSpecificDatapointsDTL)
	if warning != nil {
		t.Errorf("GetUlnarDeviation(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, golfSpecificDatapointsDTL, warning)
	}
	if math.Abs(standardActual-standardExpected) > 0.01 {
		t.Errorf("GetUlnarDeviation(%+v, %+v) = %f; expected %f", bodyDatapoints, golfSpecificDatapointsDTL, standardActual, standardExpected)
	}
	// high hands -> shift wrist closer to elbow and clubhead line
	bodyDatapoints.RWrist = &skp.Datapoint{
		X:          442.338,
		Y:          1420.001,
		Confidence: 1.0,
	}
	highExpected := 169.429
	highActual, warning := GetUlnarDeviation(bodyDatapoints, golfSpecificDatapointsDTL)
	if warning != nil {
		t.Errorf("GetUlnarDeviation(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, golfSpecificDatapointsDTL, warning)
	}
	if highActual <= standardActual {
		t.Errorf("GetUlnarDeviation(%+v, %+v) high hands result %f is supposed to be greater than the standard result %f", bodyDatapoints, golfSpecificDatapointsDTL, highActual, standardActual)
	}
	if math.Abs(highActual-highExpected) > 0.01 {
		t.Errorf("GetUlnarDeviation(%+v, %+v) = %f; expected %f", bodyDatapoints, golfSpecificDatapointsDTL, highActual, highExpected)
	}
	// low hands -> shift wrist lower than elbow and clubhead line
	bodyDatapoints.RWrist = &skp.Datapoint{
		X:          442.338,
		Y:          1603.782,
		Confidence: 1.0,
	}
	lowExpected := 130.638
	lowActual, warning := GetUlnarDeviation(bodyDatapoints, golfSpecificDatapointsDTL)
	if warning != nil {
		t.Errorf("GetUlnarDeviation(%+v, %+v) has an unexpected warning: %v", bodyDatapoints, golfSpecificDatapointsDTL, warning)
	}
	if highActual <= standardActual {
		t.Errorf("GetUlnarDeviation(%+v, %+v) low hands result %f is supposed to be less than the standard result %f", bodyDatapoints, golfSpecificDatapointsDTL, lowActual, standardActual)
	}
	if math.Abs(lowActual-lowExpected) > 0.01 {
		t.Errorf("GetUlnarDeviation(%+v, %+v) = %f; expected %f", bodyDatapoints, golfSpecificDatapointsDTL, lowActual, lowExpected)
	}
}
