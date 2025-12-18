package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirfrank96/go-server/controller"
)

func startServices(ctx context.Context, controller *controller.Controller) error {
	go func() {
		if err := controller.StartKeypointsServer(); err != nil {
			log.Fatalf("Could not start keypoints server %w", err)
		}
	}()
	log.Printf("Started Golf Keypoints Server")
	if err := controller.StartDatabaseClient(ctx); err != nil {
		return fmt.Errorf("could not start database: %w", err)
	}
	log.Printf("Started Database Client")
	if err := controller.StartCvClient(); err != nil {
		return fmt.Errorf("could not start cvclient %w", err)
	}
	log.Printf("Started CV client")
	return nil
}

func stopServices(ctx context.Context, controller *controller.Controller) error {
	if err := controller.StopKeypointsServer(); err != nil {
		return fmt.Errorf("could not stop keypoints server %w", err)
	}
	log.Printf("Stopped Golf Keypoints Server")
	if err := controller.CloseDatabaseClient(ctx); err != nil {
		return fmt.Errorf("could not stop database client %w", err)
	}
	log.Printf("Stopped Database Client")
	if err := controller.CloseCvClient(); err != nil {
		return fmt.Errorf("could not close cvclient %w", err)
	}
	log.Printf("Closed Cv client")
	return nil
}

func main() {
	ctx := context.Background()
	controller := controller.NewController()
	log.Printf("Starting services")
	err := startServices(ctx, controller)
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
	err = stopServices(ctx, controller)
	if err != nil {
		log.Fatalf("Could not stop services: %w", err)
	}
}
