package cvclient

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
	cvaddr = flag.String("cvaddr", "localhost:50051", "computervision address to connect to")
)

type CvClientManager struct {
	conn   *grpc.ClientConn
	client skp.ComputerVisionServiceClient
}

func NewCvClientManager() *CvClientManager {
	c := &CvClientManager{}
	log.Printf("New Cv Client Mgr")
	return c
}

func (c *CvClientManager) StartCvClient() error {
	log.Printf("Starting CvClient")
	flag.Parse()
	cvURI := os.Getenv("COMPUTER_VISION_URI")
	if cvURI == "" {
		cvURI = *cvaddr
	}
	// Set up a connection to the computervision server.
	var err error
	c.conn, err = grpc.NewClient(cvURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("could not connect grpc client: %w", err)
	}
	// Init ComputerVisionService grpc client
	c.client = skp.NewComputerVisionServiceClient(c.conn)
	return nil
}

func (c *CvClientManager) CloseCvClient() error {
	log.Printf("Closing CvClient")
	c.conn.Close()
	return nil
}

func (c *CvClientManager) GetPoseImage(img []byte) (*skp.GetPoseImageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	getPoseImageRequest := &skp.GetPoseImageRequest{Image: img}
	getPoseImageResponse, err := c.client.GetPoseImage(ctx, getPoseImageRequest)
	if err != nil {
		return nil, fmt.Errorf("computervision client GetPoseImage failed: %w", err)
	}
	return getPoseImageResponse, nil
}

func (c *CvClientManager) GetPoseData(img []byte) (*skp.GetPoseDataResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	getPoseDataRequest := &skp.GetPoseDataRequest{Image: img}
	getPoseDataResponse, err := c.client.GetPoseData(ctx, getPoseDataRequest)
	if err != nil {
		return nil, fmt.Errorf("computervision client GetPoseData failed: %w", err)
	}
	return getPoseDataResponse, nil
}

func (c *CvClientManager) GetPoseAll(img []byte) (*skp.GetPoseAllResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	getPoseAllRequest := &skp.GetPoseAllRequest{Image: img}
	getPoseAllResponse, err := c.client.GetPoseAll(ctx, getPoseAllRequest)
	if err != nil {
		return nil, fmt.Errorf("computervision client GetPoseAll failed: %w", err)
	}
	return getPoseAllResponse, nil
}

func (c *CvClientManager) GetPoseImagesFromFromVideo(images [][]byte) ([]*skp.GetPoseImageResponse, error) {
	// TODO: Get rid of timeout? // or configure stream vs nonstream timeouts
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	stream, err := c.client.GetPoseImagesFromVideo(ctx)
	if err != nil {
		return nil, fmt.Errorf("computervision pose client GetPoseImagesFromVideo get stream failed: %w", err)
	}

	// Start goroutine that waits for return images from stream
	waitc := make(chan struct{})
	errChan := make(chan error)
	responses := []*skp.GetPoseImageResponse{}
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
		if err := stream.Send(&skp.GetPoseImageRequest{Image: img}); err != nil {
			log.Fatalf("getPoseImagesFromVideo for img #%d stream.Send() failed: %v", idx, err)
		}
	}
	stream.CloseSend()
	log.Printf("Sent all images from cv_api_mgr")

	// Once receive stream is done, goroutine finishes
	<-waitc
	if <-errChan != nil {
		return nil, fmt.Errorf("could not get all responses: %w", err)
	}
	log.Printf("Received all images in cv_client_mgr")
	return responses, nil
}
