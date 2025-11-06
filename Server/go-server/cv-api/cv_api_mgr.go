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

func (o *OpenCvApiManager) StartOpenCvApiClient() error {
	log.Printf("Starting OpenCvApiClient")
	flag.Parse()
	// Set up a connection to the opencvandpose server.
	var err error
	o.conn, err = grpc.NewClient(*opencvaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("could not connect grpc client: %w", err)
	}
	// Init OpenCvAndPoseService grpc client
	o.c = cv.NewOpenCVAndPoseServiceClient(o.conn)
	return nil
}

func (o *OpenCvApiManager) CloseOpenCvApiClient() error {
	log.Printf("Closing OpenCvApiClient")
	o.conn.Close()
	return nil
}

func (o *OpenCvApiManager) GetOpenPoseImage(img []byte) (*cv.GetOpenPoseImageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	getOpenPoseImageRequest := &cv.GetOpenPoseImageRequest{Image: &cv.Image{Name: "", Bytes: img}}
	getOpenPoseImageResponse, err := o.c.GetOpenPoseImage(ctx, getOpenPoseImageRequest)
	if err != nil {
		return nil, fmt.Errorf("opencv/openpose client GetOpenPoseImage failed: %w", err)
	}
	return getOpenPoseImageResponse, nil
}

func (o *OpenCvApiManager) GetOpenPoseData(img []byte) (*cv.GetOpenPoseDataResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	getOpenPoseDataRequest := &cv.GetOpenPoseDataRequest{Image: &cv.Image{Name: "", Bytes: img}}
	getOpenPoseDataResponse, err := o.c.GetOpenPoseData(ctx, getOpenPoseDataRequest)
	if err != nil {
		return nil, fmt.Errorf("opencv/openpose client GetOpenPoseData failed: %w", err)
	}
	return getOpenPoseDataResponse, nil
}

func (o *OpenCvApiManager) GetOpenPoseImagesFromFromVideo(images [][]byte) ([]*cv.GetOpenPoseImageResponse, error) {
	// TODO: Get rid of timeout? // or configure stream vs nonstream timeouts
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	stream, err := o.c.GetOpenPoseImagesFromVideo(ctx)
	if err != nil {
		return nil, fmt.Errorf("opencv/open pose client GetOpenPoseImagesFromVideo get stream failed: %w", err)
	}

	// Start goroutine that waits for return images from stream
	waitc := make(chan struct{})
	errChan := make(chan error)
	responses := []*cv.GetOpenPoseImageResponse{}
	go func() {
		responseIdx := 0
		for {
			responseIdx += 1
			response, err := stream.Recv()
			if err == io.EOF { // read done
				close(waitc)
				close(errChan)
				return
			}
			if err != nil {
				errChan <- fmt.Errorf("failed to receive stream for response #%d: %w", responseIdx, err)
				return
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
	if <-errChan != nil {
		return nil, fmt.Errorf("could not get all responses: %w", err)
	}
	log.Printf("Received all images in cv_api_mgr")
	return responses, nil
}
