package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	clientapimgr "github.com/sirfrank96/go-server/client-api"
)

func startServices(clientApiMgr *clientapimgr.ClientApiManager) error {
	go func() {
		err := clientApiMgr.StartComputerVisionGolfServer()
		if err != nil {
			log.Fatalf("Could not start computervisiongolfserver %w", err)
		}
	}()
	log.Printf("Started Golf Computer Vision Server")
	err := clientApiMgr.StartOpenCvApiClient()
	if err != nil {
		return fmt.Errorf("could not start opencvapiclient %w", err)
	}
	log.Printf("Started OpenCV Api client")
	return nil
}

func stopServices(clientApiMgr *clientapimgr.ClientApiManager) error {
	err := clientApiMgr.StopGolfComputerVisionServer()
	if err != nil {
		return fmt.Errorf("could not stop computervisiongolfserver %w", err)
	}
	log.Printf("Stopped Golf Computer Vision Server")
	err = clientApiMgr.CloseOpenCvApiClient()
	if err != nil {
		return fmt.Errorf("could not close opencvapiclient %w", err)
	}
	log.Printf("Closed OpenCv Api client")
	return nil
}

func main() {
	ctx := context.Background()
	clientApiMgr := clientapimgr.NewClientApiManager(ctx)
	log.Printf("Starting services")
	err := startServices(clientApiMgr)
	if err != nil {
		log.Fatalf("Could not start services: %w", err)
	}
	// Set up a channel to listen for OS signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Waiting for computer vision golf requests at port 50052")
	log.Printf("Waiting for sigint to stop services...")
	<-stopChan
	log.Printf("Stopping services")
	err = stopServices(clientApiMgr)
	if err != nil {
		log.Fatalf("Could not stop services: %w", err)
	}
}
