package gclient

import (
	"context"
	"fmt"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	testutil "github.com/sirfrank96/go-server/test/test-util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Middle arg is a close function, should be called by calling function
func InitGolfKeypointsServiceGrpcClient(serveraddr string) (skp.GolfKeypointsServiceClient, func() error, error) {
	// Set up a connection to the cv_api server.
	conn, err := grpc.NewClient(serveraddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, conn.Close, err
	}
	// Init ComputerVisionGolf grpc client
	c := skp.NewGolfKeypointsServiceClient(conn)
	return c, conn.Close, nil
}

func UploadInputImage(ctx context.Context, gclient skp.GolfKeypointsServiceClient, sessionToken string, inputImgPath string, imageType skp.ImageType) (*skp.UploadInputImageResponse, error) {
	file, closeFile, err := testutil.GetFileFromPath(inputImgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to getFileFromPath: %w", err)
	}
	defer closeFile()
	bytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for original image: %w", err)
	}
	request := &skp.UploadInputImageRequest{
		SessionToken: sessionToken,
		ImageType:    imageType,
		Image:        bytesEncodedAsJpg,
	}
	return gclient.UploadInputImage(ctx, request)
}

func CalibrateInputImage(ctx context.Context, gclient skp.GolfKeypointsServiceClient, sessionToken string, inputImgId string, imageType skp.ImageType, calibrationType skp.CalibrationType, feetLineMethod skp.FeetLineMethod, calibrationImgAxesPath string, calibrationImgVanishingPointPath string) (*skp.CalibrateInputImageResponse, error) {
	request := &skp.CalibrateInputImageRequest{
		SessionToken:    sessionToken,
		InputImageId:    inputImgId,
		CalibrationType: calibrationType,
		FeetLineMethod:  feetLineMethod,
	}
	if calibrationType != skp.CalibrationType_NO_CALIBRATION {
		calibrationFileAxes, closeFile, err := testutil.GetFileFromPath(calibrationImgAxesPath)
		if err != nil {
			return nil, fmt.Errorf("failed to getFileFromPath calibration file axes %w", err)
		}
		defer closeFile()
		calibrationAxesBytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(calibrationFileAxes)
		if err != nil {
			return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for calibration image axes: %w", err)
		}
		request.CalibrationImageAxes = calibrationAxesBytesEncodedAsJpg

		if imageType == skp.ImageType_DTL && calibrationType != skp.CalibrationType_AXES_CALIBRATION_ONLY {
			calibrationFileVanishingPoint, closeFile, err := testutil.GetFileFromPath(calibrationImgVanishingPointPath)
			if err != nil {
				return nil, fmt.Errorf("failed to getFileFromPath calibration file vanishing point %w", err)
			}
			defer closeFile()
			calibrationVanishingBytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(calibrationFileVanishingPoint)
			if err != nil {
				return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for calibration image vanishing point: %w", err)
			}
			request.CalibrationImageVanishingPoint = calibrationVanishingBytesEncodedAsJpg
		}
	}
	return gclient.CalibrateInputImage(ctx, request)
}

func CalculateGolfKeypoints(ctx context.Context, gclient skp.GolfKeypointsServiceClient, sessionToken string, inputImgId string, outputImgPath string) (*skp.CalculateGolfKeypointsResponse, error) {
	request := &skp.CalculateGolfKeypointsRequest{
		SessionToken: sessionToken,
		InputImageId: inputImgId,
	}

	response, err := gclient.CalculateGolfKeypoints(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("could not calculate golf keypoints: %w", err)
	}

	imgSliceBytes := response.OutputImage
	jpegBytes, err := testutil.DecodeAndEncodeBytesAsJpg(imgSliceBytes)
	if err != nil {
		return response, fmt.Errorf("failed to decodeAndEncodeBytesAsJpg for return image: %w", err)
	}
	close, err := testutil.WriteBytesToJpgFile(jpegBytes, outputImgPath)
	if err != nil {
		return response, fmt.Errorf("failed to writeBytesToJpgFile: %w", err)
	}
	defer close()

	return response, nil
}

func UpdateBodyKeypoints(ctx context.Context, gclient skp.GolfKeypointsServiceClient, sessionToken string, inputImgId string, newBodyKeypoints *skp.Body25PoseKeypoints) (*skp.UpdateBodyKeypointsResponse, error) {
	request := &skp.UpdateBodyKeypointsRequest{
		SessionToken:         sessionToken,
		InputImageId:         inputImgId,
		UpdatedBodyKeypoints: newBodyKeypoints,
	}
	return gclient.UpdateBodyKeypoints(ctx, request)
}
