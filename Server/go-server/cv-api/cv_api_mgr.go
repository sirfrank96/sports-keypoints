package cvapi

import (
	"bytes"
	"context"
	"flag"
	//"image"
	"image/jpeg"
	"io"
	"log"
	//"os"
	"time"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	opencvaddr = flag.String("addr", "localhost:50051", "the address to connect to")
)

type OpenCvApiManager struct {
	conn *grpc.ClientConn
	c    cv.OpenCVAndPoseServiceClient
}

func NewOpenCvApiManager() *OpenCvApiManager {
	o := &OpenCvApiManager{}
	log.Printf("New OpenCV Api Mgr")
	return o
}

func (o *OpenCvApiManager) StartOpenCVApiClient() {
	log.Printf("Starting cv_api_mgr client")
	flag.Parse()
}

func (o *OpenCvApiManager) GetOpenPoseImage(img []byte) []byte {
	// Set up a connection to the opencvandpose server.
	var err error
	o.conn, err = grpc.NewClient(*opencvaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer o.conn.Close()
	// Init OpenCvAndPoseService grpc client
	o.c = cv.NewOpenCVAndPoseServiceClient(o.conn)
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Create rpc stream for GetOpenPoseImage
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := o.c.GetOpenPoseImage(ctx)
	if err != nil {
		log.Fatalf("c.GetOpenPoseImage failed: %v", err)
	}
	// Start goroutine that waits for return data from stream, concatenates bytes for images that are chunked
	waitc := make(chan struct{})
	imgSliceBytes := []byte{}
	go func() {
		for {
			imgReturn, err := stream.Recv()
			if err == io.EOF { // read done
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive stream: %v", err)
			}
			log.Printf("Received from stream")
			imgSliceBytes = append(imgSliceBytes, imgReturn.GetBytes()...)
		}
	}()
	// Send image via stream
	if err := stream.Send(&cv.Image{Name: "Image from go to python", Bytes: img}); err != nil {
		log.Fatalf("client.GetOpenPoseFaceOnImage: stream.Send() failed: %v", err)
	}
	stream.CloseSend()
	log.Printf("Sent data from cv_api_mgr")
	// Once receive stream is done, goroutine finishes
	<-waitc
	log.Printf("Received data in cv_api_mgr")
	// Convert bytes received to a jpg and write to a file in cwd
	imgReturnDecode, err := jpeg.Decode(bytes.NewReader(imgSliceBytes))
	if err != nil {
		log.Fatalf("failed to decode return image: %w", err)
	}
	buf := new(bytes.Buffer) //var opts jpeg.Options // opts.Quality = 80
	err = jpeg.Encode(buf, imgReturnDecode, nil)
	if err != nil {
		log.Fatalf("Failed to encode return image to JPEG: %v", err)
	}
	jpegBytes := buf.Bytes()
	return jpegBytes
}
