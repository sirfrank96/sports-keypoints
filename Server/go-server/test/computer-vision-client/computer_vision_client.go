package cvclient

import (
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

func ShowDTLPoseImage(ctx context.Context, c cv.ComputerVisionGolfServiceClient, inputImgPath string, outputImgPath string) error {
	// Send 1 image to ShowDTLPoseImage
	file, closeFile, err := testutil.GetFileFromPath(inputImgPath)
	if err != nil {
		return fmt.Errorf("failed to getFileFromPath: %w", err)
	}
	defer closeFile()
	bytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(file)
	if err != nil {
		return fmt.Errorf("failed to decodeAndEncodeFileAsJpg for original image: %w", err)
	}
	response, err := c.ShowDTLPoseImage(ctx,
		&cv.ShowDTLPoseImageRequest{
			Image: &cv.Image{
				Name:  "DTL img",
				Bytes: bytesEncodedAsJpg,
			},
		})
	if err != nil {
		return fmt.Errorf("c.ShowDTLPoseImage failed: %v", err)
	}
	imgSliceBytes := response.Image.Bytes
	jpegBytes, err := testutil.DecodeAndEncodeBytesAsJpg(imgSliceBytes)
	if err != nil {
		return fmt.Errorf("failed to decodeAndEncodeBytesAsJpg for return image: %w", err)
	}
	close, err := testutil.WriteBytesToJpgFile(jpegBytes, outputImgPath)
	if err != nil {
		return fmt.Errorf("failed to writeBytesToJpgFile: %w", err)
	}
	defer close()
	return nil
}

func ShowFaceOnPoseImage(ctx context.Context, c cv.ComputerVisionGolfServiceClient, inputImgPath string, outputImgPath string) error {
	// Send 1 image to ShowFaceOnPoseImage
	file, closeFile, err := testutil.GetFileFromPath(inputImgPath)
	if err != nil {
		return fmt.Errorf("failed to getFileFromPath: %w", err)
	}
	defer closeFile()
	bytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(file)
	if err != nil {
		return fmt.Errorf("failed to decodeAndEncodeFileAsJpg for original image: %w", err)
	}
	response, err := c.ShowFaceOnPoseImage(ctx,
		&cv.ShowFaceOnPoseImageRequest{
			Image: &cv.Image{
				Name:  "Image from go to python",
				Bytes: bytesEncodedAsJpg,
			},
		})
	if err != nil {
		return fmt.Errorf("c.ShowFaceOnPoseImage failed: %v", err)
	}
	imgSliceBytes := response.Image.Bytes
	jpegBytes, err := testutil.DecodeAndEncodeBytesAsJpg(imgSliceBytes)
	if err != nil {
		return fmt.Errorf("failed to decodeAndEncodeBytesAsJpg for return image: %w", err)
	}
	close, err := testutil.WriteBytesToJpgFile(jpegBytes, outputImgPath)
	if err != nil {
		return fmt.Errorf("failed to writeBytesToJpgFile: %w", err)
	}
	defer close()
	return nil
}

func ShowDTLPoseImagesFromVideo(ctx context.Context, c cv.ComputerVisionGolfServiceClient, inputImgPath string, outputImgPathBase string) error {
	// Start goroutine that waits for return data from stream, concatenates bytes for images that are chunked
	stream, err := c.ShowDTLPoseImagesFromVideo(ctx)
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
				&cv.ShowDTLPoseImageRequest{
					Image: &cv.Image{
						Name:  fmt.Sprintf("Image %d from go to python", numImgs),
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

func GetFaceOnPoseSetupPoints(ctx context.Context, c cv.ComputerVisionGolfServiceClient, calibrationImgPath string, inputImgPath string) (*cv.GetFaceOnPoseSetupPointsResponse, error) {
	// Get Calibration Image
	calibrationFile, closeFile, err := testutil.GetFileFromPath(calibrationImgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to getFileFromPath calibration file %w", err)
	}
	defer closeFile()
	calibrationBytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(calibrationFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for calibration image: %w", err)
	}

	file, closeFile, err := testutil.GetFileFromPath(inputImgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to getFileFromPath file: %w", err)
	}
	defer closeFile()
	bytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for image: %w", err)
	}
	getFaceOnPoseSetupPointsResponse, err := c.GetFaceOnPoseSetupPoints(ctx,
		&cv.GetFaceOnPoseSetupPointsRequest{
			CalibratedImage: &cv.CalibratedFaceOnImage{
				CalibrationImageAxes: &cv.Image{
					Name:  "Calibration img",
					Bytes: calibrationBytesEncodedAsJpg,
				},
				Image: &cv.Image{
					Name:  "Img",
					Bytes: bytesEncodedAsJpg,
				},
			},
		})
	if err != nil {
		return nil, fmt.Errorf("c.GetFaceOnPoseSetupPoints failed for image: %v", err)
	}
	return getFaceOnPoseSetupPointsResponse, nil
}

func GetDTLPoseSetupPoints(ctx context.Context, c cv.ComputerVisionGolfServiceClient, calibrationImgAxesPath string, calibrationImgVanishingPointPath string, feetLineMethod cv.FeetLineMethod, inputImgPath string) (*cv.GetDTLPoseSetupPointsResponse, error) {
	// Get Calibration Image Axes
	calibrationFileAxes, closeFile, err := testutil.GetFileFromPath(calibrationImgAxesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to getFileFromPath calibration file axes %w", err)
	}
	defer closeFile()
	calibrationAxesBytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(calibrationFileAxes)
	if err != nil {
		return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for calibration image axes: %w", err)
	}

	// Get Vanishing Point Calibration Image
	calibrationFileVanishingPoint, closeFile, err := testutil.GetFileFromPath(calibrationImgVanishingPointPath)
	if err != nil {
		return nil, fmt.Errorf("failed to getFileFromPath calibration file vanishing point %w", err)
	}
	defer closeFile()
	calibrationVanishingBytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(calibrationFileVanishingPoint)
	if err != nil {
		return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for calibration image vanishing point: %w", err)
	}

	file, closeFile, err := testutil.GetFileFromPath(inputImgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to getFileFromPath file: %w", err)
	}
	defer closeFile()
	bytesEncodedAsJpg, err := testutil.DecodeAndEncodeFileAsJpg(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decodeAndEncodeFileAsJpg for image: %w", err)
	}
	getDTLPoseSetupPointsResponse, err := c.GetDTLPoseSetupPoints(ctx,
		&cv.GetDTLPoseSetupPointsRequest{
			CalibratedImage: &cv.CalibratedDTLImage{
				FeetLineMethod: feetLineMethod,
				CalibrationImageAxes: &cv.Image{
					Name:  "Calibration img axes",
					Bytes: calibrationAxesBytesEncodedAsJpg,
				},
				CalibrationImageVanishingPoint: &cv.Image{
					Name:  "Calibration img vanishing point",
					Bytes: calibrationVanishingBytesEncodedAsJpg,
				},
				Image: &cv.Image{
					Name:  "Img",
					Bytes: bytesEncodedAsJpg,
				},
			},
		})
	if err != nil {
		return nil, fmt.Errorf("c.GetDTLPoseSetupPoints failed for image: %v", err)
	}
	return getDTLPoseSetupPointsResponse, nil
}
