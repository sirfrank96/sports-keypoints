package opencvclient

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	openposeaddr = flag.String("openposeaddr", "localhost:50051", "the address to connect to")
)

// TODO: Return errors for all functions instead of log.fatalf
// TODO: Pass around objects instead of slices???

type OpenCvClientManager struct {
	conn *grpc.ClientConn
	c    skp.OpenCVAndPoseServiceClient
}

func NewOpenCvClientManager() *OpenCvClientManager {
	o := &OpenCvClientManager{}
	log.Printf("New OpenCV Client Mgr")
	return o
}

func (o *OpenCvClientManager) StartOpenCvClient() error {
	log.Printf("Starting OpenCvClient")
	flag.Parse()
	openposeURI := os.Getenv("OPENPOSE_URI")
	if openposeURI == "" {
		openposeURI = *openposeaddr
	}
	// Set up a connection to the opencvandpose server.
	var err error
	o.conn, err = grpc.NewClient(openposeURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("could not connect grpc client: %w", err)
	}
	// Init OpenCvAndPoseService grpc client
	o.c = skp.NewOpenCVAndPoseServiceClient(o.conn)
	return nil
}

func (o *OpenCvClientManager) CloseOpenCvClient() error {
	log.Printf("Closing OpenCvClient")
	o.conn.Close()
	return nil
}

func (o *OpenCvClientManager) GetOpenPoseImage(img []byte) (*skp.GetOpenPoseImageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	getOpenPoseImageRequest := &skp.GetOpenPoseImageRequest{Image: img}
	getOpenPoseImageResponse, err := o.c.GetOpenPoseImage(ctx, getOpenPoseImageRequest)
	if err != nil {
		return nil, fmt.Errorf("opencv/openpose client GetOpenPoseImage failed: %w", err)
	}
	return getOpenPoseImageResponse, nil
}

func (o *OpenCvClientManager) GetOpenPoseData(img []byte) (*skp.GetOpenPoseDataResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	getOpenPoseDataRequest := &skp.GetOpenPoseDataRequest{Image: img}
	getOpenPoseDataResponse, err := o.c.GetOpenPoseData(ctx, getOpenPoseDataRequest)
	if err != nil {
		return nil, fmt.Errorf("opencv/openpose client GetOpenPoseData failed: %w", err)
	}
	return getOpenPoseDataResponse, nil
}

func (o *OpenCvClientManager) GetOpenPoseAll(img []byte) (*skp.GetOpenPoseAllResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	getOpenPoseAllRequest := &skp.GetOpenPoseAllRequest{Image: img}
	getOpenPoseAllResponse, err := o.c.GetOpenPoseAll(ctx, getOpenPoseAllRequest)
	if err != nil {
		return nil, fmt.Errorf("opencv/openpose client GetOpenPoseAll failed: %w", err)
	}
	return getOpenPoseAllResponse, nil
}

func (o *OpenCvClientManager) GetOpenPoseImagesFromFromVideo(images [][]byte) ([]*skp.GetOpenPoseImageResponse, error) {
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
	responses := []*skp.GetOpenPoseImageResponse{}
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
		if err := stream.Send(&skp.GetOpenPoseImageRequest{Image: img}); err != nil {
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
