// Acts as a test client for entire go-server
// Simulates mobile or desktop client that sends raw images/videos to be processed by opencv/openpose
// Pull image/video from path and send to client_api_mgr. client_api_mgr will forward to server_mgr for some processing.
// server_mgr will send to cv_api_mgr to package and send to computervision python wrapper for opencv/openpose processing.

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

/*
func testGetDTLPoseImage(ctx context.Context, c cv.ComputerVisionGolfServiceClient, userId string) {
	log.Printf("Starting testGetDTLPoseImage...\n")
	_, err := cvclient.GetDTLPoseImage(ctx, c, userId, "", filepath.Join(currentFileDirectory, "static", "dtl.jpg"), filepath.Join(currentFileDirectory, "dtl.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseImage %w", err)
	}
	log.Printf("Finished testGetDTLPoseImage...\n")
}

func testGetFaceOnPoseImage(ctx context.Context, c cv.ComputerVisionGolfServiceClient, userId string) {
	log.Printf("Starting testGetFaceOnPoseImage...\n")
	_, err := cvclient.GetFaceOnPoseImage(ctx, c, userId, "", filepath.Join(currentFileDirectory, "static", "faceon.jpg"), filepath.Join(currentFileDirectory, "faceon.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseImage %w", err)
	}
	log.Printf("Finished testGetFaceOnPoseImage...\n")
}

func testGetDTLPoseImagesFromVideo(ctx context.Context, c cv.ComputerVisionGolfServiceClient, userId string) {
	log.Printf("Starting testGetDTLPoseImagesFromVideo...\n")
	err := cvclient.GetDTLPoseImagesFromVideo(ctx, c, userId, filepath.Join(currentFileDirectory, "static", "DTLVid.mp4"), filepath.Join(currentFileDirectory, "dltvid"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseImagesFromVideo %w", err)
	}
	log.Printf("Finished testGetDTLPoseImagesFromVideo...\n")
}

func testGetFaceOnPoseSetupPoints(ctx context.Context, c cv.ComputerVisionGolfServiceClient, userId string) {
	log.Printf("Starting testGetFaceOnPoseSetupPoints...\n")

	// Side bend imgs
	log.Printf("Side bend tests...\n")
	calibrationImgPath := filepath.Join(currentFileDirectory, "static", "faceon-sidebend-goodcalibration.jpg")
	faceOnPoseSetupPointsResponse, err := cvclient.GetFaceOnPoseSetupPoints(ctx, c, userId, "", calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints neutral side bend image %w", err)
	}
	log.Printf("Neutral side bend is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.SideBend.Data)
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, userId, "", calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-sidebend-left.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints left side bend image %w", err)
	}
	log.Printf("Left side bend is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.SideBend.Data)
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, userId, "", calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-sidebend-right.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints right side bend image %w", err)
	}
	log.Printf("Right side bend is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.SideBend.Data)

	// Bad calibration img
	log.Printf("Bad calibration tests...\n")
	calibrationImgPath = filepath.Join(currentFileDirectory, "static", "faceon-sidebend-badcalibration.jpg")
	_, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, userId, "", calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon.jpg"))
	if err != nil {
		log.Printf("GetFaceOnPoseSetupPoints with bad calibration failed successfully: %v", err)
	} else {
		log.Fatalf("Supposed to get error with bad calibration image")
	}

	// Tilted calibration imgs
	log.Printf("Tilted calibration tests...\n")
	calibrationImgPath = filepath.Join(currentFileDirectory, "static", "faceon-sidebend-tilted-calibration.jpg")
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, userId, "", calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-sidebend-tilted.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints tilted neutral side bend image %w", err)
	}
	log.Printf("Tilted neutral side bend is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.SideBend.Data)
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, userId, "", calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-sidebend-tilted-left.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints tilted left side bend image %w", err)
	}
	log.Printf("Tilted left side bend is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.SideBend.Data)
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, userId, "", calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-sidebend-tilted-right.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints tilted right side bend image %w", err)
	}
	log.Printf("Tilted right side bend is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.SideBend.Data)

	// Feet flare imgs
	log.Printf("Feet flare tests...\n")
	calibrationImgPath = filepath.Join(currentFileDirectory, "static", "faceon-feetflare-calibration.jpg")
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, userId, "", calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-feetflare-neutral.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints neutral feet flare image %w", err)
	}
	log.Printf("Neutral left foot flare is %f, right foot flare is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.LFootFlare.Data, faceOnPoseSetupPointsResponse.SetupPoints.RFootFlare.Data)
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, userId, "", calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-feetflare-external.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints external feet flare image %w", err)
	}
	log.Printf("External left foot flare is %f, right foot flare is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.LFootFlare.Data, faceOnPoseSetupPointsResponse.SetupPoints.RFootFlare.Data)
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, userId, "", calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-feetflare-internal.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints internal feet flare image %w", err)
	}
	log.Printf("Internal left foot flare is %f, right foot flare is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.LFootFlare.Data, faceOnPoseSetupPointsResponse.SetupPoints.RFootFlare.Data)

	log.Printf("Finished testGetFaceOnPoseSetupPoints...\n")
}

func testGetDTLPoseSetupPoints(ctx context.Context, c cv.ComputerVisionGolfServiceClient, userId string) {
	log.Printf("Starting testGetDTLPoseSetupPoints...\n")

	// Spine angle imgs
	log.Printf("Spine angle tests...\n")
	calibrationImgAxesPath := filepath.Join(currentFileDirectory, "static", "dtl-spineangle-axescalibration.jpg")
	calibrationImgVanishingPath := filepath.Join(currentFileDirectory, "static", "dtl-spineangle-vanishingpointcalibration.jpg")
	dTLPoseSetupPointsResponse, err := cvclient.GetDTLPoseSetupPoints(ctx, c, userId, "", calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, filepath.Join(currentFileDirectory, "static", "dtl-spineangle-normal.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints normal spine angle %w", err)
	}
	log.Printf("Normal spine angle is %f", dTLPoseSetupPointsResponse.SetupPoints.SpineAngle.Data)
	dTLPoseSetupPointsResponse, err = cvclient.GetDTLPoseSetupPoints(ctx, c, userId, "", calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, filepath.Join(currentFileDirectory, "static", "dtl-spineangle-big.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints big spine angle %w", err)
	}
	log.Printf("Big spine angle is %f", dTLPoseSetupPointsResponse.SetupPoints.SpineAngle.Data)
	dTLPoseSetupPointsResponse, err = cvclient.GetDTLPoseSetupPoints(ctx, c, userId, "", calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, filepath.Join(currentFileDirectory, "static", "dtl-spineangle-small.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints small spine angle %w", err)
	}
	log.Printf("Small spine angle is %f", dTLPoseSetupPointsResponse.SetupPoints.SpineAngle.Data)

	// Feet align imgs
	log.Printf("Feet align tests...\n")
	calibrationImgAxesPath = `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-feetalign-axescalibration.jpg`
	calibrationImgVanishingPath = `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-feetalign-vanishingpointcalibration.jpg`
	_, err = cvclient.GetDTLPoseSetupPoints(ctx, c, userId, "", calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_HEEL_LINE, filepath.Join(currentFileDirectory, "static", "dtl-feetalign-neutral.jpg"))
	if err != nil {
		log.Printf("Bad heel line error is %w", err)
	} else {
		log.Fatalf("Supposed to get error with heel line that has low confidence in calibrated image")
	}
	dTLPoseSetupPointsResponse, err = cvclient.GetDTLPoseSetupPoints(ctx, c, userId, "", calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, filepath.Join(currentFileDirectory, "static", "dtl-feetalign-neutral.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints neutral feet align on toe line %w", err)
	}
	log.Printf("Neutral feet alignment on toe line is %f, associated warning msg is %s", dTLPoseSetupPointsResponse.SetupPoints.FeetAlignment.Data, dTLPoseSetupPointsResponse.SetupPoints.FeetAlignment.Warning)
	dTLPoseSetupPointsResponse, err = cvclient.GetDTLPoseSetupPoints(ctx, c, userId, "", calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, filepath.Join(currentFileDirectory, "static", "dtl-feetalign-open.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints open feet align on toe line %w", err)
	}
	log.Printf("Open feet alignment on toe line is %f, associated warning msg is %s", dTLPoseSetupPointsResponse.SetupPoints.FeetAlignment.Data, dTLPoseSetupPointsResponse.SetupPoints.FeetAlignment.Warning)
	dTLPoseSetupPointsResponse, err = cvclient.GetDTLPoseSetupPoints(ctx, c, userId, "", calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, filepath.Join(currentFileDirectory, "static", "dtl-feetalign-closed.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints closed feet align on toe line %w", err)
	}
	log.Printf("Closed feet alignment one toe line is %f", dTLPoseSetupPointsResponse.SetupPoints.FeetAlignment.Data)

	log.Printf("Finished testGetDTLPoseSetupPoints...\n")
}

func testDb(ctx context.Context, c cv.ComputerVisionGolfServiceClient) {
	user, err := c.CreateUser(ctx, &cv.CreateUserRequest{UserName: "user_1", Email: "user_1@gmail.com"})
	if err != nil {
		log.Fatalf("Error creating user %w", err)
	}
	log.Printf("CreateUser user: %+v", user)
	id := user.Id
	user, err = c.ReadUser(ctx, &cv.ReadUserRequest{Id: id})
	if err != nil {
		log.Fatalf("Error reading user %w", err)
	}
	log.Printf("ReadUser user: %+v", user)
	user, err = c.UpdateUser(ctx, &cv.UpdateUserRequest{Id: id, UserName: "user_1_update", Email: "user_1_update@gmail.com"})
	if err != nil {
		log.Fatalf("Error updating user %w", err)
	}
	log.Printf("Updated user: %+v", user)

	userId := user.Id
	// getdtlimage
	response, err := cvclient.GetDTLPoseImage(ctx, c, userId, "", filepath.Join(currentFileDirectory, "static", "dtl-spineangle-normal.jpg"), filepath.Join(currentFileDirectory, "dtl-spineangle-normal.jpg"))
	if err != nil {
		log.Fatalf("Error getting dtl image %w", err)
	}
	//readdata
	imageInfo, err := cvclient.ReadImageInfo(ctx, c, userId, response.InputImgId)
	if err != nil {
		log.Fatalf("Error reading image info %w", err)
	}
	log.Printf("Image info: userid: %s, imgid: %s, imagetype %s\n", imageInfo.UserId, imageInfo.InputImgId, imageInfo.ImageType)
	if imageInfo.OutputImg == nil {
		log.Fatalf("No output image")
	}
	if imageInfo.CalibrationImgAxes != nil {
		log.Fatalf("Calibration axes should not exist")
	}
	if imageInfo.CalibrationImgVanishingPoint != nil {
		log.Fatalf("Calibration vanishing point should not exist")
	}
	if imageInfo.OutputKeypoints != nil {
		log.Fatalf("Output keypoints should not exist")
	}
	if imageInfo.DtlGolfSetupPoints != nil {
		log.Fatalf("Dtl golf setup ponitns should not exist")
	}
	//getdtldata
	calibrationImgAxesPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-feetalign-axescalibration.jpg`
	calibrationImgVanishingPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-feetalign-vanishingpointcalibration.jpg`
	_, err = cvclient.GetDTLPoseSetupPoints(ctx, c, userId, response.InputImgId, calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, "")
	if err != nil {
		log.Fatalf("Error getting dtl setup points %w", err)
	}
	//readdata
	imageInfo, err = cvclient.ReadImageInfo(ctx, c, userId, response.InputImgId)
	if err != nil {
		log.Fatalf("Error reading image info %w", err)
	}
	log.Printf("Image info: userid: %s, imgid: %s, imagetype %s, outputkeypoitns: %+v, dtlsetuppoints: %+v\n", imageInfo.UserId, imageInfo.InputImgId, imageInfo.ImageType, imageInfo.OutputKeypoints, imageInfo.DtlGolfSetupPoints)
	if imageInfo.CalibrationImgAxes == nil {
		log.Fatalf("No calibration axes")
	}
	if imageInfo.CalibrationImgVanishingPoint == nil {
		log.Fatalf("No calibration vanishing point")
	}
	if imageInfo.OutputKeypoints == nil {
		log.Fatalf("No outputkeypoints")
	}
	if imageInfo.DtlGolfSetupPoints == nil {
		log.Fatalf("No Dtl golf setup ponitns")
	}

	deleteUserResponse, err := c.DeleteUser(ctx, &cv.DeleteUserRequest{Id: id})
	if err != nil {
		log.Fatalf("Error deleting user %w", err)
	}
	log.Printf("Deleted user: %t", deleteUserResponse.Success)
	_, err = c.ReadUser(ctx, &cv.ReadUserRequest{Id: id})
	if err != nil {
		log.Printf("Successfully got error for reading user that doesn't exist %w", err)
	}
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

	// init computer vision grpc client
	c, closeConn, err := cvclient.InitComputerVisionGolfServiceGrpcClient(*cvsportsserveraddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer closeConn()

	//createUser, then pass userId
	//testShowDTLPoseImage(ctx, c, userId)
	//testShowFaceOnPoseImage(ctx, c, userId)
	//testShowDTLPoseImagesFromVideo(ctx, c, userId)
	//testGetFaceOnPoseSetupPoints(ctx, c, userId)
	//testGetDTLPoseSetupPoints(ctx, c, userId)
	testDb(ctx, c)

	log.Printf("Ending go test_client")
}*/

func testMainCodeFlow(ctx context.Context, uClient skp.UserServiceClient, gClient skp.GolfKeypointsServiceClient) {
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

	calibrationImgAxesPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-feetalign-axescalibration.jpg`
	calibrationImgVanishingPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-feetalign-vanishingpointcalibration.jpg`
	calibrateInputImageResponse, err := gclient.CalibrateInputImage(ctx, gClient, registerUserResponse.SessionToken, uploadInputImageResponse.InputImageId, skp.ImageType_DTL, skp.CalibrationType_FULL_CALIBRATION, skp.FeetLineMethod_USE_TOE_LINE, calibrationImgAxesPath, calibrationImgVanishingPath)
	if err != nil {
		log.Fatalf("Failed to calibrate input image: %s", err.Error())
	}
	log.Printf("CalibrateInputImageResponse: %+v", calibrateInputImageResponse)

	calculateGolfKeypointsResponse, err := gclient.CalculateGolfKeypoints(ctx, gClient, registerUserResponse.SessionToken, uploadInputImageResponse.InputImageId, filepath.Join(currentFileDirectory, "dtl-spineangle-normal.jpg"))
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

	testMainCodeFlow(ctx, uClient, gClient)

	log.Printf("Ending go test_client")
}
