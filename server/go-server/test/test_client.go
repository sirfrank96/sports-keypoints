// Acts as a test client for entire go-server
// Simulates mobile or desktop client that sends raw images/videos to be processed by computervision
// Pull image/video from path and send to client_api_mgr. client_api_mgr will forward to server_mgr for some processing.
// server_mgr will send to cv_api_mgr to package and send to computervision python wrapper for computervision processing.

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	gclient "github.com/sirfrank96/go-server/test/golf-keypoints-client"
	uclient "github.com/sirfrank96/go-server/test/user-client"
)

var (
	cvsportsserveraddr = flag.String("addr", "localhost:50052", "the address to connect to")
)

var currentFileDirectory string

func testMainCodeFlowDtl(ctx context.Context, uClient skp.UserServiceClient, gClient skp.GolfKeypointsServiceClient) {
	createUserResponse, err := uclient.CreateUser(ctx, uClient, "user_1", "password123", "user_1@gmail.com")
	if err != nil {
		log.Fatalf("Failed to create user: %s", err.Error())
	}
	log.Printf("CreateUserResponse: %+v", createUserResponse)

	registerUserResponse, err := uclient.RegisterUser(ctx, uClient, "user_1", "password123")
	if err != nil {
		log.Fatalf("Failed to register user: %s", err.Error())
	}
	log.Printf("RegisterUserResponse: %+v", registerUserResponse)

	uploadInputImageResponse, err := gclient.UploadInputImage(ctx, gClient, registerUserResponse.SessionToken, filepath.Join(currentFileDirectory, "static", "dtl-feetalign-neutral.jpg"), skp.ImageType_DTL)
	if err != nil {
		log.Fatalf("Failed to upload input image: %s", err.Error())
	}
	log.Printf("UploadInputImageResponse: %+v", uploadInputImageResponse)

	calibrationImgAxesPath := filepath.Join(currentFileDirectory, "static", "dtl-feetalign-axescalibration.jpg")
	calibrationImgVanishingPath := filepath.Join(currentFileDirectory, "static", "dtl-feetalign-vanishingpointcalibration.jpg")
	calibrateInputImageResponse, err := gclient.CalibrateInputImage(ctx, gClient, registerUserResponse.SessionToken, uploadInputImageResponse.InputImageId, skp.ImageType_DTL, skp.CalibrationType_FULL_CALIBRATION, skp.FeetLineMethod_USE_HEEL_LINE, calibrationImgAxesPath, calibrationImgVanishingPath)
	if err != nil {
		log.Fatalf("Failed to calibrate input image: %s", err.Error())
	}
	log.Printf("CalibrateInputImageResponse: %+v", calibrateInputImageResponse)

	calculateGolfKeypointsResponse, err := gclient.CalculateGolfKeypoints(ctx, gClient, registerUserResponse.SessionToken, uploadInputImageResponse.InputImageId, filepath.Join(currentFileDirectory, "dtl-feetalign-neutral-output.jpg"))
	if err != nil {
		log.Fatalf("Failed to calculate golf keypoints: %s", err.Error())
	}
	log.Printf("Calculate golf keypoints dtl: %+v", calculateGolfKeypointsResponse.GolfKeypoints.DtlGolfSetupPoints)
	newBodyKeypoints := &skp.Body25PoseKeypoints{
		LShoulder: &skp.Keypoint{
			X:          408,
			Y:          1176,
			Confidence: 0.99,
		},
	}
	updateBodyKeypointsResponse, err := gclient.UpdateBodyKeypoints(ctx, gClient, registerUserResponse.SessionToken, uploadInputImageResponse.InputImageId, newBodyKeypoints)
	if err != nil {
		log.Fatalf("Failed to update body keypoints: %s", err.Error())
	}
	log.Printf("Update body keypoints dtl: %+v", updateBodyKeypointsResponse.UpdatedGolfKeypoints.DtlGolfSetupPoints)
}

func testMainCodeFlowFaceOn(ctx context.Context, uClient skp.UserServiceClient, gClient skp.GolfKeypointsServiceClient) {
	createUserResponse, err := uclient.CreateUser(ctx, uClient, "user_1", "password123", "user_1@gmail.com")
	if err != nil {
		log.Fatalf("Failed to create user: %s", err.Error())
	}
	log.Printf("CreateUserResponse: %+v", createUserResponse)

	registerUserResponse, err := uclient.RegisterUser(ctx, uClient, "user_1", "password123")
	if err != nil {
		log.Fatalf("Failed to register user: %s", err.Error())
	}
	log.Printf("RegisterUserResponse: %+v", registerUserResponse)

	uploadInputImageResponse, err := gclient.UploadInputImage(ctx, gClient, registerUserResponse.SessionToken, filepath.Join(currentFileDirectory, "static", "faceon.jpg"), skp.ImageType_FACE_ON)
	if err != nil {
		log.Fatalf("Failed to upload input image: %s", err.Error())
	}
	log.Printf("UploadInputImageResponse: %+v", uploadInputImageResponse)

	calibrationImgPath := filepath.Join(currentFileDirectory, "static", "faceon-sidebend-goodcalibration.jpg")
	calibrateInputImageResponse, err := gclient.CalibrateInputImage(ctx, gClient, registerUserResponse.SessionToken, uploadInputImageResponse.InputImageId, skp.ImageType_FACE_ON, skp.CalibrationType_FULL_CALIBRATION, skp.FeetLineMethod_USE_TOE_LINE, calibrationImgPath, "")
	if err != nil {
		log.Fatalf("Failed to calibrate input image: %s", err.Error())
	}
	log.Printf("CalibrateInputImageResponse: %+v", calibrateInputImageResponse)

	calculateGolfKeypointsResponse, err := gclient.CalculateGolfKeypoints(ctx, gClient, registerUserResponse.SessionToken, uploadInputImageResponse.InputImageId, filepath.Join(currentFileDirectory, "dtl-spineangle-normal.jpg"))
	if err != nil {
		log.Fatalf("Failed to calculate golf keypoints: %s", err.Error())
	}
	log.Printf("Calculate golf keypoints face on: %+v", calculateGolfKeypointsResponse.GolfKeypoints.FaceonGolfSetupPoints)
}

func testNoCalibration(ctx context.Context, uClient skp.UserServiceClient, gClient skp.GolfKeypointsServiceClient) {
	createUserResponse, err := uclient.CreateUser(ctx, uClient, "user_1", "password123", "user_1@gmail.com")
	if err != nil {
		log.Fatalf("Failed to create user: %s", err.Error())
	}
	log.Printf("CreateUserResponse: %+v", createUserResponse)

	registerUserResponse, err := uclient.RegisterUser(ctx, uClient, "user_1", "password123")
	if err != nil {
		log.Fatalf("Failed to register user: %s", err.Error())
	}
	log.Printf("RegisterUserResponse: %+v", registerUserResponse)

	uploadInputImageResponse, err := gclient.UploadInputImage(ctx, gClient, registerUserResponse.SessionToken, filepath.Join(currentFileDirectory, "static", "dtl-spineangle-normal.jpg"), skp.ImageType_DTL)
	if err != nil {
		log.Fatalf("Failed to upload input image: %s", err.Error())
	}
	log.Printf("UploadInputImageResponse: %+v", uploadInputImageResponse)

	// no call to calibrateInputImage
	calculateGolfKeypointsResponse, err := gclient.CalculateGolfKeypoints(ctx, gClient, registerUserResponse.SessionToken, uploadInputImageResponse.InputImageId, filepath.Join(currentFileDirectory, "dtl-spineangle-normal.jpg"))
	if err != nil {
		log.Fatalf("Failed to calculate golf keypoints: %s", err.Error())
	}
	log.Printf("Calculate golf keypoints dtl: %+v", calculateGolfKeypointsResponse.GolfKeypoints.DtlGolfSetupPoints)

	// calibrateInputImage NO_CALIBRATION
	calibrateInputImageResponse, err := gclient.CalibrateInputImage(ctx, gClient, registerUserResponse.SessionToken, uploadInputImageResponse.InputImageId, skp.ImageType_DTL, skp.CalibrationType_NO_CALIBRATION, skp.FeetLineMethod_USE_TOE_LINE, "", "")
	if err != nil {
		log.Fatalf("Failed to calibrate input image: %s", err.Error())
	}
	log.Printf("CalibrateInputImageResponse: %+v", calibrateInputImageResponse)

	calculateGolfKeypointsResponse, err = gclient.CalculateGolfKeypoints(ctx, gClient, registerUserResponse.SessionToken, uploadInputImageResponse.InputImageId, filepath.Join(currentFileDirectory, "dtl-spineangle-normal.jpg"))
	if err != nil {
		log.Fatalf("Failed to calculate golf keypoints: %s", err.Error())
	}
	log.Printf("Calculate golf keypoints dtl: %+v", calculateGolfKeypointsResponse.GolfKeypoints.DtlGolfSetupPoints)

	// calibrateInputImage AXES_CALIBRATION_ONLY
	calibrationImgAxesPath := filepath.Join(currentFileDirectory, "static", "dtl-feetalign-axescalibration.jpg")
	calibrateInputImageResponse, err = gclient.CalibrateInputImage(ctx, gClient, registerUserResponse.SessionToken, uploadInputImageResponse.InputImageId, skp.ImageType_DTL, skp.CalibrationType_AXES_CALIBRATION_ONLY, skp.FeetLineMethod_USE_TOE_LINE, calibrationImgAxesPath, "")
	if err != nil {
		log.Fatalf("Failed to calibrate input image: %s", err.Error())
	}
	log.Printf("CalibrateInputImageResponse: %+v", calibrateInputImageResponse)

	calculateGolfKeypointsResponse, err = gclient.CalculateGolfKeypoints(ctx, gClient, registerUserResponse.SessionToken, uploadInputImageResponse.InputImageId, filepath.Join(currentFileDirectory, "dtl-spineangle-normal.jpg"))
	if err != nil {
		log.Fatalf("Failed to calculate golf keypoints: %s", err.Error())
	}
	log.Printf("Calculate golf keypoints dtl: %+v", calculateGolfKeypointsResponse.GolfKeypoints.DtlGolfSetupPoints)
}

func main() {
	log.Printf("Starting test_client")
	ctx := context.Background()
	flag.Parse()

	// get current executables path
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("Could not get current executable info")
	}
	currentFileDirectory = path.Dir(executable)

	// init user grpc client
	uClient, closeUserConn, err := uclient.InitUserServiceGrpcClient(*cvsportsserveraddr)
	if err != nil {
		log.Fatalf("Failed to connect to user: %v", err)
	}
	defer closeUserConn()
	// init golf keypoints grpc client
	gClient, closeGolfConn, err := gclient.InitGolfKeypointsServiceGrpcClient(*cvsportsserveraddr)
	if err != nil {
		log.Fatalf("Failed to connect to golf keypoints: %v", err)
	}
	defer closeGolfConn()

	testMainCodeFlowDtl(ctx, uClient, gClient)
	testMainCodeFlowFaceOn(ctx, uClient, gClient)
	testNoCalibration(ctx, uClient, gClient)

	log.Printf("Ending go test_client")
}
