package gclient

import (
	"context"
	"fmt"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	testutil "github.com/sirfrank96/go-server/test/test-util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//TODO: RETURN ERRORS INSTEAD OF LOG.FATALF

// Middle arg is a close function, should be called by calling function
func InitGolfKeypointsServiceGrpcClient(serveraddr string) (cv.GolfKeypointsServiceClient, func() error, error) {
	// Set up a connection to the cv_api server.
	conn, err := grpc.NewClient(serveraddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, conn.Close, err
	}
	// Init ComputerVisionGolf grpc client
	c := cv.NewGolfKeypointsServiceClient(conn)
	return c, conn.Close, nil
}

func UploadInputImage(ctx context.Context, gclient cv.GolfKeypointsServiceClient, sessionToken string, inputImgPath string, imageType cv.ImageType) (*cv.UploadInputImageResponse, error) {
	file, closeFile, err := testutil.GetFileFromPath(inputImgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to getFileFromPath: %w", err)
	}
	defer closeFile()
	bytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for original image: %w", err)
	}
	request := &cv.UploadInputImageRequest{
		SessionToken: sessionToken,
		ImageType:    imageType,
		Image:        bytesEncodedAsJpg,
	}
	return gclient.UploadInputImage(ctx, request)
}

func CalibrateInputImage(ctx context.Context, gclient cv.GolfKeypointsServiceClient, sessionToken string, inputImgId string, imageType cv.ImageType, calibrationType cv.CalibrationType, feetLineMethod cv.FeetLineMethod, calibrationImgAxesPath string, calibrationImgVanishingPointPath string) (*cv.CalibrateInputImageResponse, error) {
	request := &cv.CalibrateInputImageRequest{
		SessionToken:    sessionToken,
		InputImageId:    inputImgId,
		CalibrationType: calibrationType,
		FeetLineMethod:  feetLineMethod,
	}
	if calibrationType != cv.CalibrationType_NO_CALIBRATION {
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

		if imageType == cv.ImageType_DTL && calibrationType != cv.CalibrationType_AXES_CALIBRATION_ONLY {
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

func CalculateGolfKeypoints(ctx context.Context, gclient cv.GolfKeypointsServiceClient, sessionToken string, inputImgId string, outputImgPath string) (*cv.CalculateGolfKeypointsResponse, error) {
	request := &cv.CalculateGolfKeypointsRequest{
		SessionToken: sessionToken,
		InputImageId: inputImgId,
	}

	response, err := gclient.CalculateGolfKeypoints(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("could not calculate golf keypoints: %w", err)
	}

	imgSliceBytes := response.GolfKeypoints.OutputImg
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
