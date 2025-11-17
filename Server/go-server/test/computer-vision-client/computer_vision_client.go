package cvclient

/*import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	testutil "github.com/sirfrank96/go-server/test/test-util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//TODO: RETURN ERRORS INSTEAD OF LOG.FATALF

// Middle arg is a close function, should be called by calling function
func InitComputerVisionGolfServiceGrpcClient(serveraddr string) (cv.ComputerVisionGolfServiceClient, func() error, error) {
	// Set up a connection to the cv_api server.
	conn, err := grpc.NewClient(serveraddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, conn.Close, err
	}
	// Init ComputerVisionGolf grpc client
	c := cv.NewComputerVisionGolfServiceClient(conn)
	return c, conn.Close, nil
}

func ReadImageInfo(ctx context.Context, c cv.ComputerVisionGolfServiceClient, userId string, imgId string) (*cv.ImageInfo, error) {
	return c.ReadImageInfo(ctx, &cv.ReadImageInfoRequest{
		UserId: userId,
		ImgId:  imgId,
	})
}

func GetDTLPoseImage(ctx context.Context, c cv.ComputerVisionGolfServiceClient, userId string, imgId string, inputImgPath string, outputImgPath string) (*cv.GetDTLPoseImageResponse, error) {
	var request cv.GetDTLPoseImageRequest
	if imgId != "" {
		request = cv.GetDTLPoseImageRequest{
			UserId: userId,
			ImgId:  imgId,
		}
	} else {
		// Send 1 image to ShowDTLPoseImage
		file, closeFile, err := testutil.GetFileFromPath(inputImgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to getFileFromPath: %w", err)
		}
		defer closeFile()
		bytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for original image: %w", err)
		}
		request = cv.GetDTLPoseImageRequest{
			UserId: userId,
			Image: &cv.Image{
				Bytes: bytesEncodedAsJpg,
			},
		}
	}
	response, err := c.GetDTLPoseImage(ctx, &request)
	if err != nil {
		return nil, fmt.Errorf("c.ShowDTLPoseImage failed: %v", err)
	}
	imgSliceBytes := response.Image.Bytes
	jpegBytes, err := testutil.DecodeAndEncodeBytesAsJpg(imgSliceBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decodeAndEncodeBytesAsJpg for return image: %w", err)
	}
	close, err := testutil.WriteBytesToJpgFile(jpegBytes, outputImgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to writeBytesToJpgFile: %w", err)
	}
	defer close()
	return response, nil
}

func GetFaceOnPoseImage(ctx context.Context, c cv.ComputerVisionGolfServiceClient, userId string, imgId string, inputImgPath string, outputImgPath string) (*cv.GetFaceOnPoseImageResponse, error) {
	var request cv.GetFaceOnPoseImageRequest
	if imgId != "" {
		request = cv.GetFaceOnPoseImageRequest{
			UserId: userId,
			ImgId:  imgId,
		}
	} else {
		// Send 1 image to ShowFaceOnPoseImage
		file, closeFile, err := testutil.GetFileFromPath(inputImgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to getFileFromPath: %w", err)
		}
		defer closeFile()
		bytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for original image: %w", err)
		}
		request = cv.GetFaceOnPoseImageRequest{
			UserId: userId,
			Image: &cv.Image{
				Bytes: bytesEncodedAsJpg,
			},
		}
	}
	response, err := c.GetFaceOnPoseImage(ctx, &request)
	if err != nil {
		return nil, fmt.Errorf("c.ShowFaceOnPoseImage failed: %v", err)
	}
	imgSliceBytes := response.Image.Bytes
	jpegBytes, err := testutil.DecodeAndEncodeBytesAsJpg(imgSliceBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decodeAndEncodeBytesAsJpg for return image: %w", err)
	}
	close, err := testutil.WriteBytesToJpgFile(jpegBytes, outputImgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to writeBytesToJpgFile: %w", err)
	}
	defer close()
	return response, nil
}

func GetDTLPoseImagesFromVideo(ctx context.Context, c cv.ComputerVisionGolfServiceClient, userId string, inputImgPath string, outputImgPathBase string) error {
	// Start goroutine that waits for return data from stream, concatenates bytes for images that are chunked
	stream, err := c.GetDTLPoseImagesFromVideo(ctx)
	if err != nil {
		return fmt.Errorf("c.ShowDTLPoseImagesFromVideo failed: %v", err)
	}
	waitc := make(chan struct{})
	returnImages := [][]byte{}
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF { // read done
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive stream: %v", err)
			}
			log.Printf("Received from stream")
			returnImages = append(returnImages, response.Image.Bytes)
		}
	}()

	// Get video and break up into jpgs to send via stream
	cmd := exec.Command("ffmpeg", "-i", inputImgPath, "-f", "image2pipe", "-c:v", "mjpeg", "-r", "5", "pipe:1") //5 fps
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run FFmpeg: %w", err)
	}
	xFFWasPrevByte := false
	var currImgBytes []byte
	numImgs := 0
	for {
		b, err := out.ReadByte()
		if err != nil {
			if err == io.EOF {
				log.Printf("Finished reading file")
				break
			} else {
				return fmt.Errorf("could not read byte: %w", err)
			}
		}
		currImgBytes = append(currImgBytes, b)
		if b == byte(0xff) {
			xFFWasPrevByte = true
		} else if b == byte(0xd8) && xFFWasPrevByte { // start of an img
			xFFWasPrevByte = false
			currImgBytes = []byte{byte(0xff), byte(0xd8)}
			numImgs += 1
		} else if b == byte(0xd9) && xFFWasPrevByte { // end of an img
			xFFWasPrevByte = false
			bytesEncodedAsJpg, err := testutil.DecodeAndEncodeBytesAsJpg(currImgBytes)
			if err != nil {
				return fmt.Errorf("failed to decodeAndEncodeFileAsJpg for original image: %w", err)
			}
			err = stream.Send(
				&cv.GetDTLPoseImageRequest{
					UserId: userId,
					Image: &cv.Image{
						Bytes: bytesEncodedAsJpg,
					},
				})
			if err != nil {
				return fmt.Errorf("cv.GetDTLPoseImagesFromVideo: stream.Send() failed: %v", err)
			}
			log.Printf("Sent img # %d, size of img is %d", numImgs, len(currImgBytes))
		} else {
			xFFWasPrevByte = false
		}
	}

	stream.CloseSend()
	log.Printf("Sent all data in showDTLPoseImagesFromVideo\n")
	<-waitc
	log.Printf("Received all data in showDTLPoseImagesFromVideo\n")

	// iterate over all processed images and output to jpgs
	for idx, imgSliceBytes := range returnImages {
		jpegBytes, err := testutil.DecodeAndEncodeBytesAsJpg(imgSliceBytes)
		if err != nil {
			return fmt.Errorf("failed to decodeAndEncodeBytesAsJpg for return image: %w", err)
		}
		currOutputImgPath := outputImgPathBase + strconv.Itoa(idx) + ".jpg"
		close, err := testutil.WriteBytesToJpgFile(jpegBytes, currOutputImgPath)
		if err != nil {
			return fmt.Errorf("failed to writeBytesToJpgFile: %w", err)
		}
		defer close()
	}
	return nil
}

func GetFaceOnPoseSetupPoints(ctx context.Context, c cv.ComputerVisionGolfServiceClient, userId string, imgId string, calibrationImgPath string, inputImgPath string) (*cv.GetFaceOnPoseSetupPointsResponse, error) {
	request := &cv.GetFaceOnPoseSetupPointsRequest{
		UserId:          userId,
		CalibratedImage: &cv.CalibratedFaceOnImage{},
	}
	if imgId != "" { // Get input image from db
		request.ImgId = imgId
	} else { // Get input image from path
		file, closeFile, err := testutil.GetFileFromPath(inputImgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to getFileFromPath file: %w", err)
		}
		defer closeFile()
		bytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for image: %w", err)
		}
		request.CalibratedImage.Image = &cv.Image{
			Bytes: bytesEncodedAsJpg,
		}
	}
	// Get Calibration Image
	if calibrationImgPath != "" {
		calibrationFileAxes, closeFile, err := testutil.GetFileFromPath(calibrationImgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to getFileFromPath calibration file axes %w", err)
		}
		defer closeFile()
		calibrationAxesBytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(calibrationFileAxes)
		if err != nil {
			return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for calibration image axes: %w", err)
		}
		request.CalibratedImage.CalibrationImageAxes = &cv.Image{
			Bytes: calibrationAxesBytesEncodedAsJpg,
		}
	}

	getFaceOnPoseSetupPointsResponse, err := c.GetFaceOnPoseSetupPoints(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("c.GetFaceOnPoseSetupPoints failed for image: %v", err)
	}
	return getFaceOnPoseSetupPointsResponse, nil
}

func GetDTLPoseSetupPoints(ctx context.Context, c cv.ComputerVisionGolfServiceClient, userId string, imgId string, calibrationImgAxesPath string, calibrationImgVanishingPointPath string, feetLineMethod cv.FeetLineMethod, inputImgPath string) (*cv.GetDTLPoseSetupPointsResponse, error) {
	request := &cv.GetDTLPoseSetupPointsRequest{
		UserId: userId,
		CalibratedImage: &cv.CalibratedDTLImage{
			FeetLineMethod: feetLineMethod,
		},
	}
	if imgId != "" { // Get input image from db
		request.ImgId = imgId
	} else { // Get input image from path
		file, closeFile, err := testutil.GetFileFromPath(inputImgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to getFileFromPath file: %w", err)
		}
		defer closeFile()
		bytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for image: %w", err)
		}
		request.CalibratedImage.Image = &cv.Image{
			Bytes: bytesEncodedAsJpg,
		}
	}
	// Get Calibration Image Axes
	if calibrationImgAxesPath != "" {
		calibrationFileAxes, closeFile, err := testutil.GetFileFromPath(calibrationImgAxesPath)
		if err != nil {
			return nil, fmt.Errorf("failed to getFileFromPath calibration file axes %w", err)
		}
		defer closeFile()
		calibrationAxesBytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(calibrationFileAxes)
		if err != nil {
			return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for calibration image axes: %w", err)
		}
		request.CalibratedImage.CalibrationImageAxes = &cv.Image{
			Bytes: calibrationAxesBytesEncodedAsJpg,
		}
	}
	// Get Vanishing Point Calibration Image
	if calibrationImgVanishingPointPath != "" {
		calibrationFileVanishingPoint, closeFile, err := testutil.GetFileFromPath(calibrationImgVanishingPointPath)
		if err != nil {
			return nil, fmt.Errorf("failed to getFileFromPath calibration file vanishing point %w", err)
		}
		defer closeFile()
		calibrationVanishingBytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(calibrationFileVanishingPoint)
		if err != nil {
			return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for calibration image vanishing point: %w", err)
		}
		request.CalibratedImage.CalibrationImageVanishingPoint = &cv.Image{
			Bytes: calibrationVanishingBytesEncodedAsJpg,
		}
	}

	getDTLPoseSetupPointsResponse, err := c.GetDTLPoseSetupPoints(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("c.GetDTLPoseSetupPoints failed for image: %v", err)
	}
	return getDTLPoseSetupPointsResponse, nil
}*/
