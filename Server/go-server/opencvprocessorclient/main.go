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

	cv "github.com/sirfrank96/go-server/computervision"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	log.Printf("Starting go client")
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := cv.NewComputerVisionClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	img, err := os.Open(`C:\Users\Franklin\Desktop\Computer Vision Golf\Server\go-server\static\faceon.jpg`)
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

	//imgReturn, err := c.GetOpenPoseFaceOnImage(ctx, &cv.FaceOnImage{Name: "Faceon from go to python", Image: buffer.Bytes()})
	//if err != nil {
	//	log.Fatalf("grpc error: %v", err)
	//}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := c.GetOpenPoseFaceOnImage(ctx)
	if err != nil {
		log.Fatalf("c.GetOpenPoseFaceOnImage failed: %v", err)
	}

	// wait for processed image
	waitc := make(chan struct{})
	imgSliceBytes := []byte{}
	go func() {
		for {
			imgReturn, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive stream: %v", err)
			}
			log.Printf("Received from stream")
			imgSliceBytes = append(imgSliceBytes, imgReturn.GetImage()...)
		}
	}()

	// Send image
	if err := stream.Send(&cv.FaceOnImage{Name: "Faceon from go to python", Image: buffer.Bytes()}); err != nil {
		log.Fatalf("client.GetOpenPoseFaceOnImage: stream.Send() failed: %v", err)
	}
	stream.CloseSend()
	log.Printf("Sent data")

	<-waitc
	log.Printf("Received data")

	log.Printf("starting to decode bytes and encode to jpg")
	// converting bytes to a jpg
	log.Printf("byte slice size is %d", len(imgSliceBytes))
	imgReturnDecode, err := jpeg.Decode(bytes.NewReader(imgSliceBytes))
	if err != nil {
		log.Fatalf("failed to decode return image: %w", err)
	}

	//var buf bytes.Buffer
	//var opts jpeg.Options
	//opts.Quality = 80 // Set quality to 80
	//err = jpeg.Encode(&buf, imgReturnDecode, &opts)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, imgReturnDecode, nil)
	if err != nil {
		log.Fatalf("Failed to encode return image to JPEG: %v", err)
	}

	jpegBytes := buf.Bytes()

	jpegFile, err := os.Create(`C:\Users\Franklin\Desktop\Computer Vision Golf\Server\go-server\opencvprocessorclient\test.jpg`)
	if err != nil {
		log.Fatalf("Failed to create test.jpg: %v", err)
	}
	defer jpegFile.Close()

	_, err = jpegFile.Write(jpegBytes)
	if err != nil {
		log.Fatalf("Failed to write JPeG file: %v", err)
	}

	log.Printf("Ending go client")

}
