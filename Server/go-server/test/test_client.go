// Acts as a test client for entire go-server
// Simulates mobile or desktop client that sends raw images/videos to be processed by opencv/openpose
// Pull image/video from path and send to client_api_mgr. client_api_mgr will forward to server_mgr for some processing.
// server_mgr will send to cv_api_mgr to package and send to computervision python wrapper for opencv/openpose processing.

package main

import (
	"bytes"
	"context"
	"flag"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"time"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cvsportsserveraddr = flag.String("addr", "localhost:50052", "the address to connect to")
)

func grabImageDecodeAndEncodeAsJpg(path string) *bytes.Buffer {
	// Grab example image to process, decode image, then encode as jpg
	img, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open jpg file: %w", err)
	}
	defer img.Close()
	imgDecode, _, err := image.Decode(img)
	if err != nil {
		log.Fatalf("failed to decode original image: %w", err)
	}
	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, imgDecode, nil)
	if err != nil {
		log.Fatalf("Error encoding original image to jpeg: %w", err)
	}
	return buffer
}

func main() {
	log.Printf("Starting test_client")
	ctx := context.Background()
	flag.Parse()

	// Set up a connection to the cv_api server.
	conn, err := grpc.NewClient(*cvsportsserveraddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Init ComputerVisionGolf grpc client
	c := cv.NewGolfComputerVisionServiceClient(conn)
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	faceOnPath := `C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\static\faceon.jpg`
	buffer := grabImageDecodeAndEncodeAsJpg(faceOnPath)

	stream, err := c.ShowDTLPoseImage(ctx)
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
	if err := stream.Send(&cv.Image{Name: "Image from go to python", Bytes: buffer.Bytes()}); err != nil {
		log.Fatalf("client.GetOpenPoseFaceOnImage: stream.Send() failed: %v", err)
	}
	stream.CloseSend()
	log.Printf("Sent data in test_client")
	// Once receive stream is done, goroutine finishes
	<-waitc
	log.Printf("Received data in test_client")

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
	jpegFile, err := os.Create(`C:\Users\Franklin\Desktop\Computer Vision Sports\Server\go-server\test\test.jpg`)
	if err != nil {
		log.Fatalf("Failed to create test.jpg: %v", err)
	}
	defer jpegFile.Close()
	_, err = jpegFile.Write(jpegBytes)
	if err != nil {
		log.Fatalf("Failed to write JPeG file: %v", err)
	}

	log.Printf("Ending go test_client")
}
