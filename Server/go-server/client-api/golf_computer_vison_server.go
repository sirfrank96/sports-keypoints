package clientapi

import (
	"context"
	"io"
	//"log"
	//"sync"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	processor "github.com/sirfrank96/go-server/processor"
)

type golfComputerVisionServer struct {
	cv.UnimplementedGolfComputerVisionServiceServer
	ctx context.Context
	//mutex sync.Mutex
	prcsr *processor.Processor
}

func createNewGolfComputerVisionServer(ctx context.Context, processor *processor.Processor) *golfComputerVisionServer {
	g := &golfComputerVisionServer{}
	g.ctx = ctx
	g.prcsr = processor
	return g
}

func (g *golfComputerVisionServer) ShowDTLPoseImage(stream cv.GolfComputerVisionService_ShowDTLPoseImageServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		processedImg := g.prcsr.ShowDTLPoseImage(in.Bytes)
		if err = stream.Send(&cv.Image{Name: "Processed image", Bytes: processedImg}); err != nil {
			return err
		}
	}
	return nil
}

func (g *golfComputerVisionServer) IdentifyDTLPoseDifferences(stream cv.GolfComputerVisionService_IdentifyDTLPoseDifferencesServer) error {
	return nil
}

func (g *golfComputerVisionServer) ShowFaceOnPoseImage(stream cv.GolfComputerVisionService_ShowFaceOnPoseImageServer) error {
	return nil
}

func (g *golfComputerVisionServer) IdentifyFaceOnPoseDifferences(stream cv.GolfComputerVisionService_IdentifyFaceOnPoseDifferencesServer) error {
	return nil
}
