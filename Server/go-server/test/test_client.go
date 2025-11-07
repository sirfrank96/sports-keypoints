// Acts as a test client for entire go-server
// Simulates mobile or desktop client that sends raw images/videos to be processed by opencv/openpose
// Pull image/video from path and send to client_api_mgr. client_api_mgr will forward to server_mgr for some processing.
// server_mgr will send to cv_api_mgr to package and send to computervision python wrapper for opencv/openpose processing.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cvsportsserveraddr = flag.String("addr", "localhost:50052", "the address to connect to")
)

// TODO: 1 conn and client
// Middle arg is a close function, should be called by calling function
func initComputerVisionGolfServiceGrpcClient(serveraddr string) (cv.ComputerVisionGolfServiceClient, func() error, error) {
	// Set up a connection to the cv_api server.
	conn, err := grpc.NewClient(*cvsportsserveraddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, conn.Close, err
	}
	// Init ComputerVisionGolf grpc client
	c := cv.NewComputerVisionGolfServiceClient(conn)
	return c, conn.Close, nil
}

// Middle arg is a close function, should be called by calling function
func getFileFromPath(path string) (*os.File, func() error, error) {
	// Grab example image to process, decode image, then encode as jpg
	file, err := os.Open(path)
	if err != nil {
		return nil, file.Close, fmt.Errorf("failed to open file: %w", err)
	}
	return file, file.Close, nil
}

func decodeAndEncodeFileAsJpg(file *os.File) ([]byte, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file: %w", err)
	}
	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode img to jpg: %w", err)
	}
	return buffer.Bytes(), nil
}

func decodeAndEncodeBytesAsJpg(byteSlice []byte) ([]byte, error) {
	// Convert bytes received to a jpg and write to a file in cwd
	imgReturnDecode, err := jpeg.Decode(bytes.NewReader(byteSlice))
	if err != nil {
		return nil, fmt.Errorf("failed to decode return image: %w", err)
	}
	buf := new(bytes.Buffer) //var opts jpeg.Options // opts.Quality = 80
	err = jpeg.Encode(buf, imgReturnDecode, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode return image to jpg: %w", err)
	}
	return buf.Bytes(), nil
}

// first arg is a close function, should be called by calling function
func writeBytesToJpgFile(byteSlice []byte, path string) (func() error, error) {
	jpegFile, err := os.Create(path)
	if err != nil {
		return jpegFile.Close, fmt.Errorf("failed to create test.jpg: %v", err)
	}
	_, err = jpegFile.Write(byteSlice)
	if err != nil {
		return jpegFile.Close, fmt.Errorf("failed to write jpg file: %v", err)
	}
	return jpegFile.Close, nil
}

func testShowDTLPoseImage(ctx context.Context) {
	log.Printf("Starting testShowDTLPoseImage...")
	c, closeConn, err := initComputerVisionGolfServiceGrpcClient(*cvsportsserveraddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer closeConn()
	// Send 1 image to ShowDTLPoseImage
	dtlPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl.jpg`
	file, closeFile, err := getFileFromPath(dtlPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath: %w", err)
	}
	defer closeFile()
	bytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(file)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for original image: %w", err)
	}
	response, err := c.ShowDTLPoseImage(ctx, &cv.ShowDTLPoseImageRequest{Image: &cv.Image{Name: "Image from go to python", Bytes: bytesEncodedAsJpg}})
	if err != nil {
		log.Fatalf("c.GetOpenPoseImage failed: %v", err)
	}
	log.Printf("Sent and received data in testShowDTLPoseImage")
	imgSliceBytes := response.Image.Bytes
	jpegBytes, err := decodeAndEncodeBytesAsJpg(imgSliceBytes)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeBytesAsJpg for return image: %w", err)
	}
	close, err := writeBytesToJpgFile(jpegBytes, `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\dtl.jpg`)
	if err != nil {
		log.Fatalf("Failed to writeBytesToJpgFile: %w", err)
	}
	defer close()
}

func testShowFaceOnPoseImage(ctx context.Context) {
	log.Printf("Starting testShowFaceOnPoseImage...")
	c, closeConn, err := initComputerVisionGolfServiceGrpcClient(*cvsportsserveraddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer closeConn()
	// Send 1 image to ShowDTLPoseImage
	faceOnPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\faceon.jpg`
	file, closeFile, err := getFileFromPath(faceOnPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath: %w", err)
	}
	defer closeFile()
	bytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(file)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for original image: %w", err)
	}
	response, err := c.ShowFaceOnPoseImage(ctx, &cv.ShowFaceOnPoseImageRequest{Image: &cv.Image{Name: "Image from go to python", Bytes: bytesEncodedAsJpg}})
	if err != nil {
		log.Fatalf("c.GetOpenPoseImage failed: %v", err)
	}
	log.Printf("Sent and received data in testShowDTLPoseImage")
	imgSliceBytes := response.Image.Bytes
	jpegBytes, err := decodeAndEncodeBytesAsJpg(imgSliceBytes)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeBytesAsJpg for return image: %w", err)
	}
	close, err := writeBytesToJpgFile(jpegBytes, `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\faceon.jpg`)
	if err != nil {
		log.Fatalf("Failed to writeBytesToJpgFile: %w", err)
	}
	defer close()
}

func testShowDTLPoseImagesFromVideo(ctx context.Context) {
	log.Printf("Starting testShowDTLPoseImagesFromVideo...")
	c, closeConn, err := initComputerVisionGolfServiceGrpcClient(*cvsportsserveraddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer closeConn()
	// Start goroutine that waits for return data from stream, concatenates bytes for images that are chunked
	stream, err := c.ShowDTLPoseImagesFromVideo(ctx)
	if err != nil {
		log.Fatalf("c.GetOpenPoseImage failed: %v", err)
	}
	waitc := make(chan struct{})
	returnImages := [][]byte{}
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF { // read done
				log.Printf("Read EOF")
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
	dtlVideoPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\DTLVid.mp4`
	cmd := exec.Command("ffmpeg", "-i", dtlVideoPath, "-f", "image2pipe", "-c:v", "mjpeg", "-r", "5", "pipe:1") //5 fps
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run FFmpeg: %w", err)
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
				log.Fatalf("Could not read byte: %w", err)
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
			bytesEncodedAsJpg, err := decodeAndEncodeBytesAsJpg(currImgBytes)
			if err != nil {
				log.Fatalf("Failed to decodeAndEncodeFileAsJpg for original image: %w", err)
			}
			if err := stream.Send(&cv.ShowDTLPoseImageRequest{Image: &cv.Image{Name: fmt.Sprintf("Image %d from go to python", numImgs), Bytes: bytesEncodedAsJpg}}); err != nil {
				log.Fatalf("client.GetOpenPoseFaceOnImage: stream.Send() failed: %v", err)
			}
			log.Printf("Sent img # %d, size of img is %d", numImgs, len(currImgBytes))
		} else {
			xFFWasPrevByte = false
		}
	}

	stream.CloseSend()
	log.Printf("Sent all data in testShowDTLPoseImagesFromVideo")
	<-waitc
	log.Printf("Received all data in testShowDTLPoseImagesFromVideo")

	// iterate over all processed images and output to jpgs
	for idx, imgSliceBytes := range returnImages {
		jpegBytes, err := decodeAndEncodeBytesAsJpg(imgSliceBytes)
		if err != nil {
			log.Fatalf("Failed to decodeAndEncodeBytesAsJpg for return image: %w", err)
		}
		close, err := writeBytesToJpgFile(jpegBytes, `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\dltvid`+strconv.Itoa(idx)+`.jpg`)
		if err != nil {
			log.Fatalf("Failed to writeBytesToJpgFile: %w", err)
		}
		defer close()
	}
}

func testGetFaceOnPoseSetupPoints(ctx context.Context) {
	log.Printf("Starting testGetFaceOnPoseSetupPoints...")
	c, closeConn, err := initComputerVisionGolfServiceGrpcClient(*cvsportsserveraddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer closeConn()

	// Get Calibration Image
	calibrationImgPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\faceon-goodcalibration.jpg`
	calibrationfile, closeFile, err := getFileFromPath(calibrationImgPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath calibration file %w", err)
	}
	defer closeFile()
	calibrationBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(calibrationfile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for calibration image: %w", err)
	}

	// Get Neutral Side Bend Image
	neutralImgPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\faceon.jpg`
	neutralFile, closeFile, err := getFileFromPath(neutralImgPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath neutralfile: %w", err)
	}
	defer closeFile()
	neutralBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(neutralFile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for neutral image: %w", err)
	}
	getFaceOnPoseSetupPointsNeutralResponse, err := c.GetFaceOnPoseSetupPoints(ctx, &cv.GetFaceOnPoseSetupPointsRequest{CalibratedImage: &cv.CalibratedFaceOnImage{CalibrationImageAxes: &cv.Image{Name: "Calibration img", Bytes: calibrationBytesEncodedAsJpg}, Image: &cv.Image{Name: "Neutral side bend img", Bytes: neutralBytesEncodedAsJpg}}})
	if err != nil {
		log.Fatalf("c.GetFaceOnPoseSetupPoints failed for neutral image: %v", err)
	}
	log.Printf("Sent and received data for neutral side bend GetFaceonPoseSetupPoints")
	log.Printf("Neutral side bend is %f", getFaceOnPoseSetupPointsNeutralResponse.SetupPoints.SideBend)

	// Get Left Side Bend Image
	leftSideBendImgPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\faceon-leftsidebend.jpg`
	leftSideBendFile, closeFile, err := getFileFromPath(leftSideBendImgPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath leftSideBendFile: %w", err)
	}
	defer closeFile()
	leftSideBendBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(leftSideBendFile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for left side bend image: %w", err)
	}
	getFaceOnPoseSetupPointsLeftSideBendResponse, err := c.GetFaceOnPoseSetupPoints(ctx, &cv.GetFaceOnPoseSetupPointsRequest{CalibratedImage: &cv.CalibratedFaceOnImage{CalibrationImageAxes: &cv.Image{Name: "Calibration img", Bytes: calibrationBytesEncodedAsJpg}, Image: &cv.Image{Name: "Left side bend img", Bytes: leftSideBendBytesEncodedAsJpg}}})
	if err != nil {
		log.Fatalf("c.GetFaceOnPoseSetupPoints failed for left side bend image: %v", err)
	}
	log.Printf("Sent and received data for left side bend testGetFaceOnPoseSetupPoints")
	log.Printf("Left side bend is %f", getFaceOnPoseSetupPointsLeftSideBendResponse.SetupPoints.SideBend)

	// Get Right Side Bend Image
	rightSideBendImgPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\faceon-rightsidebend.jpg`
	rightSideBendFile, closeFile, err := getFileFromPath(rightSideBendImgPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath rightSideBendFile: %w", err)
	}
	defer closeFile()
	rightSideBendBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(rightSideBendFile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for right side bend image: %w", err)
	}
	getFaceOnPoseSetupPointsRightSideBendResponse, err := c.GetFaceOnPoseSetupPoints(ctx, &cv.GetFaceOnPoseSetupPointsRequest{CalibratedImage: &cv.CalibratedFaceOnImage{CalibrationImageAxes: &cv.Image{Name: "Calibration img", Bytes: calibrationBytesEncodedAsJpg}, Image: &cv.Image{Name: "Right side bend img", Bytes: rightSideBendBytesEncodedAsJpg}}})
	if err != nil {
		log.Fatalf("c.GetFaceOnPoseSetupPoints failed for right side bend image: %v", err)
	}
	log.Printf("Sent and received data for right side bend testGetFaceOnPoseSetupPoints")
	log.Printf("Right side bend is %f", getFaceOnPoseSetupPointsRightSideBendResponse.SetupPoints.SideBend)
}

func testGetFaceOnPoseSetupPointsBadCalibration(ctx context.Context) {
	log.Printf("Starting testGetFaceOnPoseSetupPointsBadCalibration...")
	c, closeConn, err := initComputerVisionGolfServiceGrpcClient(*cvsportsserveraddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer closeConn()

	// Get Bad Calibration Image
	badCalibrationImgPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\faceon-badcalibration.jpg`
	calibrationfile, closeFile, err := getFileFromPath(badCalibrationImgPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath bad calibration file %w", err)
	}
	defer closeFile()
	calibrationBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(calibrationfile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for bad calibration image: %w", err)
	}

	_, err = c.GetFaceOnPoseSetupPoints(ctx, &cv.GetFaceOnPoseSetupPointsRequest{CalibratedImage: &cv.CalibratedFaceOnImage{CalibrationImageAxes: &cv.Image{Name: "Calibration img", Bytes: calibrationBytesEncodedAsJpg}, Image: nil}})
	if err != nil {
		log.Printf("c.GetFaceOnPoseSetupPointsBadCalibration failed successfully: %v", err)
	} else {
		log.Fatalf("Supposed to get error with bad calibration image")
	}
}

func testGetFaceOnPoseSetupPointsTiltedCalibration(ctx context.Context) {
	log.Printf("Starting testGetFaceOnPoseSetupPointsTiltedCalibration...")
	c, closeConn, err := initComputerVisionGolfServiceGrpcClient(*cvsportsserveraddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer closeConn()

	// Get Tilted Calibration Image
	calibrationImgPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\faceon-tiltedcalibration.jpg`
	calibrationfile, closeFile, err := getFileFromPath(calibrationImgPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath tilted calibration file %w", err)
	}
	defer closeFile()
	calibrationBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(calibrationfile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for tilted calibration image: %w", err)
	}

	// Get Tilted Neutral Side Bend Image
	neutralImgPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\faceon-tilted.jpg`
	neutralFile, closeFile, err := getFileFromPath(neutralImgPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath tilted neutralfile: %w", err)
	}
	defer closeFile()
	neutralBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(neutralFile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for tilted neutral image: %w", err)
	}
	getFaceOnPoseSetupPointsNeutralResponse, err := c.GetFaceOnPoseSetupPoints(ctx, &cv.GetFaceOnPoseSetupPointsRequest{CalibratedImage: &cv.CalibratedFaceOnImage{CalibrationImageAxes: &cv.Image{Name: "Calibration img", Bytes: calibrationBytesEncodedAsJpg}, Image: &cv.Image{Name: "Neutral side bend img", Bytes: neutralBytesEncodedAsJpg}}})
	if err != nil {
		log.Fatalf("c.GetFaceOnPoseSetupPoints failed for tilted neutral image: %v", err)
	}
	log.Printf("Sent and received data for neutral side bend GetFaceonPoseSetupPointsTiltedCalibration")
	log.Printf("Tilted neutral side bend is %f", getFaceOnPoseSetupPointsNeutralResponse.SetupPoints.SideBend)

	// Get Tilted Left Side Bend Image
	leftSideBendImgPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\faceon-tiltedleftsidebend.jpg`
	leftSideBendFile, closeFile, err := getFileFromPath(leftSideBendImgPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath tilted leftSideBendFile: %w", err)
	}
	defer closeFile()
	leftSideBendBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(leftSideBendFile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for tilted left side bend image: %w", err)
	}
	getFaceOnPoseSetupPointsLeftSideBendResponse, err := c.GetFaceOnPoseSetupPoints(ctx, &cv.GetFaceOnPoseSetupPointsRequest{CalibratedImage: &cv.CalibratedFaceOnImage{CalibrationImageAxes: &cv.Image{Name: "Calibration img", Bytes: calibrationBytesEncodedAsJpg}, Image: &cv.Image{Name: "Left side bend img", Bytes: leftSideBendBytesEncodedAsJpg}}})
	if err != nil {
		log.Fatalf("c.GetFaceOnPoseSetupPoints failed for tilted left side bend image: %v", err)
	}
	log.Printf("Sent and received data for left side bend testGetFaceOnPoseSetupPointsTiltedCalibration")
	log.Printf("Tilted left side bend is %f", getFaceOnPoseSetupPointsLeftSideBendResponse.SetupPoints.SideBend)

	// Get Tilted Right Side Bend Image
	rightSideBendImgPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\faceon-tiltedrightsidebend.jpg`
	rightSideBendFile, closeFile, err := getFileFromPath(rightSideBendImgPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath tilted rightSideBendFile: %w", err)
	}
	defer closeFile()
	rightSideBendBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(rightSideBendFile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for tilted right side bend image: %w", err)
	}
	getFaceOnPoseSetupPointsRightSideBendResponse, err := c.GetFaceOnPoseSetupPoints(ctx, &cv.GetFaceOnPoseSetupPointsRequest{CalibratedImage: &cv.CalibratedFaceOnImage{CalibrationImageAxes: &cv.Image{Name: "Calibration img", Bytes: calibrationBytesEncodedAsJpg}, Image: &cv.Image{Name: "Right side bend img", Bytes: rightSideBendBytesEncodedAsJpg}}})
	if err != nil {
		log.Fatalf("c.GetFaceOnPoseSetupPoints failed for tilted right side bend image: %v", err)
	}
	log.Printf("Sent and received data for right side bend testGetFaceOnPoseSetupPointsTiltedCalibration")
	log.Printf("Tilted right side bend is %f", getFaceOnPoseSetupPointsRightSideBendResponse.SetupPoints.SideBend)
}

func testGetDTLPoseSetupPoints(ctx context.Context) {
	log.Printf("Starting testGetDTLPoseSetupPoints...")
	c, closeConn, err := initComputerVisionGolfServiceGrpcClient(*cvsportsserveraddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer closeConn()

	// Get Calibration Axes Image
	calibrationImgAxesPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-goodcalibrationaxes.jpg`
	calibrationAxesFile, closeFile, err := getFileFromPath(calibrationImgAxesPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath calibration axes file %w", err)
	}
	defer closeFile()
	calibrationAxesBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(calibrationAxesFile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for calibration axes image: %w", err)
	}

	// Get Calibration Vanishing Point Image
	calibrationImgVanishingPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-goodcalibrationvanishingpoint.jpg`
	calibrationVanishingFile, closeFile, err := getFileFromPath(calibrationImgVanishingPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath calibration vanishing file %w", err)
	}
	defer closeFile()
	calibrationVanishingBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(calibrationVanishingFile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for calibration vanishing image: %w", err)
	}

	// Get Normal spine angle
	normalImgPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-normal.jpg`
	normalFile, closeFile, err := getFileFromPath(normalImgPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath normal file: %w", err)
	}
	defer closeFile()
	normalBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(normalFile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for normal image: %w", err)
	}
	normalRequest := &cv.GetDTLPoseSetupPointsRequest{
		CalibratedImage: &cv.CalibratedDTLImage{
			CalibrationImageAxes: &cv.Image{
				Name:  "Calibration axes img",
				Bytes: calibrationAxesBytesEncodedAsJpg,
			},
			CalibrationImageVanishingPoint: &cv.Image{
				Name:  "Calibration vanishing point img",
				Bytes: calibrationVanishingBytesEncodedAsJpg,
			},
			Image: &cv.Image{
				Name:  "Normal side bend img",
				Bytes: normalBytesEncodedAsJpg,
			},
		},
	}
	getDTLPoseSetupPointsNormalResponse, err := c.GetDTLPoseSetupPoints(ctx, normalRequest)
	if err != nil {
		log.Fatalf("c.GetDTLPoseSetupPoints failed for normal image: %v", err)
	}
	log.Printf("Sent and received data for normal spine angle GetDTLPoseSetupPoints")
	log.Printf("Normal spine angle is %f", getDTLPoseSetupPointsNormalResponse.SetupPoints.SpineAngle)

	// Get Big spine angle
	bigSpineImgPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-bigspineangle.jpg`
	bigSpineFile, closeFile, err := getFileFromPath(bigSpineImgPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath big spine file: %w", err)
	}
	defer closeFile()
	bigSpineBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(bigSpineFile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for big spine image: %w", err)
	}
	bigSpineRequest := &cv.GetDTLPoseSetupPointsRequest{
		CalibratedImage: &cv.CalibratedDTLImage{
			CalibrationImageAxes: &cv.Image{
				Name:  "Calibration axes img",
				Bytes: calibrationAxesBytesEncodedAsJpg,
			},
			CalibrationImageVanishingPoint: &cv.Image{
				Name:  "Calibration vanishing point img",
				Bytes: calibrationVanishingBytesEncodedAsJpg,
			},
			Image: &cv.Image{
				Name:  "Big side bend img",
				Bytes: bigSpineBytesEncodedAsJpg,
			},
		},
	}
	getDTLPoseSetupPointsBigSpineResponse, err := c.GetDTLPoseSetupPoints(ctx, bigSpineRequest)
	if err != nil {
		log.Fatalf("c.GetDTLPoseSetupPoints failed for big spine image: %v", err)
	}
	log.Printf("Sent and received data for big spine angle GetDTLPoseSetupPoints")
	log.Printf("Big spine angle is %f", getDTLPoseSetupPointsBigSpineResponse.SetupPoints.SpineAngle)

	// Get Small spine angle
	smallSpineImgPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\static\dtl-smallspineangle.jpg`
	smallSpineFile, closeFile, err := getFileFromPath(smallSpineImgPath)
	if err != nil {
		log.Fatalf("Failed to getFileFromPath small spine file: %w", err)
	}
	defer closeFile()
	smallSpineBytesEncodedAsJpg, err := decodeAndEncodeFileAsJpg(smallSpineFile)
	if err != nil {
		log.Fatalf("Failed to decodeAndEncodeFileAsJpg for small spine image: %w", err)
	}
	smallSpineRequest := &cv.GetDTLPoseSetupPointsRequest{
		CalibratedImage: &cv.CalibratedDTLImage{
			CalibrationImageAxes: &cv.Image{
				Name:  "Calibration axes img",
				Bytes: calibrationAxesBytesEncodedAsJpg,
			},
			CalibrationImageVanishingPoint: &cv.Image{
				Name:  "Calibration vanishing point img",
				Bytes: calibrationVanishingBytesEncodedAsJpg,
			},
			Image: &cv.Image{
				Name:  "Small side bend img",
				Bytes: smallSpineBytesEncodedAsJpg,
			},
		},
	}
	getDTLPoseSetupPointsSmallSpineResponse, err := c.GetDTLPoseSetupPoints(ctx, smallSpineRequest)
	if err != nil {
		log.Fatalf("c.GetDTLPoseSetupPoints failed for small spine image: %v", err)
	}
	log.Printf("Sent and received data for small spine angle GetDTLPoseSetupPoints")
	log.Printf("Small spine angle is %f", getDTLPoseSetupPointsSmallSpineResponse.SetupPoints.SpineAngle)
}

func main() {
	log.Printf("Starting test_client")
	ctx := context.Background()
	flag.Parse()

	//testShowDTLPoseImage(ctx)

	//testShowFaceOnPoseImage(ctx)

	//testShowDTLPoseImagesFromVideo(ctx)

	//testGetFaceOnPoseSetupPoints(ctx)

	//testGetFaceOnPoseSetupPointsBadCalibration(ctx)

	//testGetFaceOnPoseSetupPointsTiltedCalibration(ctx)

	testGetDTLPoseSetupPoints(ctx)

	log.Printf("Ending go test_client")
}
