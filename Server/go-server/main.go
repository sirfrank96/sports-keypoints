package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	clientapimgr "github.com/sirfrank96/go-server/client-api"
)

func startServices(clientApiMgr *clientapimgr.ClientApiManager) error {
	go clientApiMgr.StartGolfComputerVisionServer()
	log.Printf("Started Golf Computer Vision Server")
	go clientApiMgr.StartOpenCVApiClient()
	log.Printf("Started OpenCV Api client")
	return nil
}

func stopServices(clientApiMgr *clientapimgr.ClientApiManager) error {
	clientApiMgr.StopGolfComputerVisionServer()
	return nil
}

func main() {
	ctx := context.Background()

	clientApiMgr := clientapimgr.NewClientApiManager(ctx)
	log.Printf("Starting services")
	startServices(clientApiMgr)

	// Set up a channel to listen for OS signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Waiting for sigint")
	<-stopChan

	log.Printf("Stopping services")
	stopServices(clientApiMgr)
}
