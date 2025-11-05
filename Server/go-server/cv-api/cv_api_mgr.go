package cvapi

import (
	//"bytes"
	"context"
	"flag"
	//"image"
	//"image/jpeg"
	"io"
	"log"
	//"os"
	"fmt"
	"time"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	opencvaddr = flag.String("addr", "localhost:50051", "the address to connect to")
)

// TODO: Return errors for all functions instead of log.fatalf
// TODO: Pass around objects instead of slices???

type OpenCvApiManager struct {
	conn *grpc.ClientConn
	c    cv.OpenCVAndPoseServiceClient
}

func NewOpenCvApiManager() *OpenCvApiManager {
	o := &OpenCvApiManager{}
	log.Printf("New OpenCV Api Mgr")
	return o
}

func (o *OpenCvApiManager) StartOpenCvApiClient() {
	log.Printf("Starting OpenCvApiClient")
	flag.Parse()
	// Set up a connection to the opencvandpose server.
	var err error
	o.conn, err = grpc.NewClient(*opencvaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// Init OpenCvAndPoseService grpc client
	o.c = cv.NewOpenCVAndPoseServiceClient(o.conn)
}

func (o *OpenCvApiManager) CloseOpenCvApiClient() {
	log.Printf("Closing OpenCvApiClient")
	o.conn.Close()
}

func (o *OpenCvApiManager) GetOpenPoseImage(img []byte) *cv.GetOpenPoseImageResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	getOpenPoseImageRequest := &cv.GetOpenPoseImageRequest{Image: &cv.Image{Name: "", Bytes: img}}
	getOpenPoseImageResponse, err := o.c.GetOpenPoseImage(ctx, getOpenPoseImageRequest)
	if err != nil {
		log.Fatalf("c.GetOpenPoseImage failed: %v", err)
	}
	return getOpenPoseImageResponse
}

func (o *OpenCvApiManager) GetOpenPoseImagesFromFromVideo(images [][]byte) []*cv.GetOpenPoseImageResponse {
	// Create rpc stream for GetOpenPoseImage
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	stream, err := o.c.GetOpenPoseImagesFromVideo(ctx)
	if err != nil {
		log.Fatalf("c.GetOpenPoseImagesFromVideo initial call failed: %v", err)
	}
	// Start goroutine that waits for return images from stream
	waitc := make(chan struct{})
	responses := []*cv.GetOpenPoseImageResponse{}
	go func() {
		responseIdx := 0
		for {
			responseIdx += 1
			response, err := stream.Recv()
			if err == io.EOF { // read done
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive stream for response #%d: %v", responseIdx, err)
			}
			log.Printf("Received from stream, response #%d", responseIdx)
			responses = append(responses, response)

		}
	}()

	// Send all images via stream
	for idx, img := range images {
		if err := stream.Send(&cv.GetOpenPoseImageRequest{Image: &cv.Image{Name: fmt.Sprintf("Img #%d", idx), Bytes: img}}); err != nil {
			log.Fatalf("getOpenPoseImagesFromVideo for img #%d stream.Send() failed: %v", idx, err)
		}
	}

	stream.CloseSend()
	log.Printf("Sent all images from cv_api_mgr")
	// Once receive stream is done, goroutine finishes
	<-waitc
	log.Printf("Received all images in cv_api_mgr")

	// TODO: Probably dont have to decode and encode (python wrapper will have done this)
	// Decode and encode images as jpg
	/*processedImages := [][]byte{}
	for idx, imgReturn := range imagesReturned {
		imgReturnDecode, err := jpeg.Decode(bytes.NewReader(imgReturn))
		if err != nil {
			log.Fatalf("failed to decode return image #%d: %w", idx, err)
		}
		buf := new(bytes.Buffer) //var opts jpeg.Options // opts.Quality = 80
		err = jpeg.Encode(buf, imgReturnDecode, nil)
		if err != nil {
			log.Fatalf("Failed to encode return image #%d to jpg: %v", idx, err)
		}
		jpegBytes := buf.Bytes()
		processedImages = append(processedImages, jpegBytes)
	}*/
	return responses
}
