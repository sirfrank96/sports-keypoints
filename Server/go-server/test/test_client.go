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

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	cvclient "github.com/sirfrank96/go-server/test/computer-vision-client"
)

var (
	cvsportsserveraddr = flag.String("addr", "localhost:50052", "the address to connect to")
)

var currentFileDirectory string

func testShowDTLPoseImage(ctx context.Context, c cv.ComputerVisionGolfServiceClient) {
	log.Printf("Starting testShowDTLPoseImage...\n")
	err := cvclient.ShowDTLPoseImage(ctx, c, filepath.Join(currentFileDirectory, "static", "dtl.jpg"), filepath.Join(currentFileDirectory, "dtl.jpg"))
	if err != nil {
		log.Fatalf("Error in ShowDTLPoseImage %w", err)
	}
	log.Printf("Finished testShowDTLPoseImage...\n")
}

func testShowFaceOnPoseImage(ctx context.Context, c cv.ComputerVisionGolfServiceClient) {
	log.Printf("Starting testShowFaceOnPoseImage...\n")
	err := cvclient.ShowFaceOnPoseImage(ctx, c, filepath.Join(currentFileDirectory, "static", "faceon.jpg"), filepath.Join(currentFileDirectory, "faceon.jpg"))
	if err != nil {
		log.Fatalf("Error in ShowFaceOnPoseImage %w", err)
	}
	log.Printf("Finished testShowFaceOnPoseImage...\n")
}

func testShowDTLPoseImagesFromVideo(ctx context.Context, c cv.ComputerVisionGolfServiceClient) {
	log.Printf("Starting testShowDTLPoseImagesFromVideo...\n")
	err := cvclient.ShowDTLPoseImagesFromVideo(ctx, c, filepath.Join(currentFileDirectory, "static", "DTLVid.mp4"), filepath.Join(currentFileDirectory, "dltvid"))
	if err != nil {
		log.Fatalf("Error in ShowDTLPoseImagesFromVideo %w", err)
	}
	log.Printf("Finished testShowDTLPoseImagesFromVideo...\n")
}

func testGetFaceOnPoseSetupPoints(ctx context.Context, c cv.ComputerVisionGolfServiceClient) {
	log.Printf("Starting testGetFaceOnPoseSetupPoints...\n")

	// Side bend imgs
	log.Printf("Side bend tests...\n")
	calibrationImgPath := filepath.Join(currentFileDirectory, "static", "faceon-sidebend-goodcalibration.jpg")
	faceOnPoseSetupPointsResponse, err := cvclient.GetFaceOnPoseSetupPoints(ctx, c, calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints neutral image %w", err)
	}
	log.Printf("Neutral side bend is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.SideBend.Data)
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-sidebend-left.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints left side bend image %w", err)
	}
	log.Printf("Left side bend is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.SideBend.Data)
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-sidebend-right.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints right side bend image %w", err)
	}
	log.Printf("Right side bend is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.SideBend.Data)

	// Bad calibration img
	log.Printf("Bad calibration tests...\n")
	calibrationImgPath = filepath.Join(currentFileDirectory, "static", "faceon-sidebend-badcalibration.jpg")
	_, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon.jpg"))
	if err != nil {
		log.Printf("GetFaceOnPoseSetupPoints with bad calibration failed successfully: %v", err)
	} else {
		log.Fatalf("Supposed to get error with bad calibration image")
	}

	// Tilted calibration imgs
	log.Printf("Tilted calibration tests...\n")
	calibrationImgPath = filepath.Join(currentFileDirectory, "static", "faceon-sidebend-tilted-calibration.jpg")
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-sidebend-tilted.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints tilted neutral side bend image %w", err)
	}
	log.Printf("Tilted neutral side bend is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.SideBend.Data)
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-sidebend-tilted-left.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints tilted left side bend image %w", err)
	}
	log.Printf("Tilted left side bend is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.SideBend.Data)
	faceOnPoseSetupPointsResponse, err = cvclient.GetFaceOnPoseSetupPoints(ctx, c, calibrationImgPath, filepath.Join(currentFileDirectory, "static", "faceon-sidebend-tilted-right.jpg"))
	if err != nil {
		log.Fatalf("Error in GetFaceOnPoseSetupPoints tilted right side bend image %w", err)
	}
	log.Printf("Tilted right side bend is %f\n", faceOnPoseSetupPointsResponse.SetupPoints.SideBend.Data)

	log.Printf("Finished testGetFaceOnPoseSetupPoints...\n")
}

func testGetDTLPoseSetupPoints(ctx context.Context, c cv.ComputerVisionGolfServiceClient) {
	log.Printf("Starting testGetDTLPoseSetupPoints...\n")

	// Spine angle imgs
	log.Printf("Spine angle tests...\n")
	calibrationImgAxesPath := filepath.Join(currentFileDirectory, "static", "dtl-spineangle-axescalibration.jpg")
	calibrationImgVanishingPath := filepath.Join(currentFileDirectory, "static", "dtl-spineangle-vanishingpointcalibration.jpg")
	dTLPoseSetupPointsResponse, err := cvclient.GetDTLPoseSetupPoints(ctx, c, calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, filepath.Join(currentFileDirectory, "static", "dtl-spineangle-normal.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints normal spine angle %w", err)
	}
	log.Printf("Normal spine angle is %f", dTLPoseSetupPointsResponse.SetupPoints.SpineAngle.Data)
	dTLPoseSetupPointsResponse, err = cvclient.GetDTLPoseSetupPoints(ctx, c, calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, filepath.Join(currentFileDirectory, "static", "dtl-spineangle-big.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints big spine angle %w", err)
	}
	log.Printf("Big spine angle is %f", dTLPoseSetupPointsResponse.SetupPoints.SpineAngle.Data)
	dTLPoseSetupPointsResponse, err = cvclient.GetDTLPoseSetupPoints(ctx, c, calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, filepath.Join(currentFileDirectory, "static", "dtl-spineangle-small.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints small spine angle %w", err)
	}
	log.Printf("Small spine angle is %f", dTLPoseSetupPointsResponse.SetupPoints.SpineAngle.Data)

	// Feet align imgs
	log.Printf("Feet align tests...\n")
	calibrationImgAxesPath = `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-feetalign-axescalibration.jpg`
	calibrationImgVanishingPath = `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-feetalign-vanishingpointcalibration.jpg`
	dTLPoseSetupPointsResponse, err = cvclient.GetDTLPoseSetupPoints(ctx, c, calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_HEEL_LINE, filepath.Join(currentFileDirectory, "static", "dtl-feetalign-neutral.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints neutral feet align on heel line %w", err)
	}
	log.Printf("Neutral feet alignment on heel line (bad heel line) is %f, associated warning msg is %s", dTLPoseSetupPointsResponse.SetupPoints.FeetAlignment.Data, dTLPoseSetupPointsResponse.SetupPoints.FeetAlignment.Warning)
	dTLPoseSetupPointsResponse, err = cvclient.GetDTLPoseSetupPoints(ctx, c, calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, filepath.Join(currentFileDirectory, "static", "dtl-feetalign-neutral.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints neutral feet align on toe line %w", err)
	}
	log.Printf("Neutral feet alignment on toe line is %f, associated warning msg is %s", dTLPoseSetupPointsResponse.SetupPoints.FeetAlignment.Data, dTLPoseSetupPointsResponse.SetupPoints.FeetAlignment.Warning)
	dTLPoseSetupPointsResponse, err = cvclient.GetDTLPoseSetupPoints(ctx, c, calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, filepath.Join(currentFileDirectory, "static", "dtl-feetalign-open.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints open feet align on toe line %w", err)
	}
	log.Printf("Open feet alignment on toe line is %f, associated warning msg is %s", dTLPoseSetupPointsResponse.SetupPoints.FeetAlignment.Data, dTLPoseSetupPointsResponse.SetupPoints.FeetAlignment.Warning)
	dTLPoseSetupPointsResponse, err = cvclient.GetDTLPoseSetupPoints(ctx, c, calibrationImgAxesPath, calibrationImgVanishingPath, cv.FeetLineMethod_USE_TOE_LINE, filepath.Join(currentFileDirectory, "static", "dtl-feetalign-closed.jpg"))
	if err != nil {
		log.Fatalf("Error in GetDTLPoseSetupPoints closed feet align on toe line %w", err)
	}
	log.Printf("Closed feet alignment one toe line is %f", dTLPoseSetupPointsResponse.SetupPoints.FeetAlignment.Data)

	log.Printf("Finished testGetDTLPoseSetupPoints...\n")
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

	testShowDTLPoseImage(ctx, c)
	testShowFaceOnPoseImage(ctx, c)
	testShowDTLPoseImagesFromVideo(ctx, c)
	testGetFaceOnPoseSetupPoints(ctx, c)
	testGetDTLPoseSetupPoints(ctx, c)

	log.Printf("Ending go test_client")
}
