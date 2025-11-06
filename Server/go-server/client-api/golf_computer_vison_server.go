package clientapi

import (
	"context"
	"io"
	//"log"
	//"sync"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	processor "github.com/sirfrank96/go-server/processor"
)

type computerVisionGolfServer struct {
	cv.UnimplementedComputerVisionGolfServiceServer
	ctx context.Context
	//mutex sync.Mutex
	prcsr *processor.Processor
}

func createNewComputerVisionGolfServer(ctx context.Context, processor *processor.Processor) *computerVisionGolfServer {
	c := &computerVisionGolfServer{}
	c.ctx = ctx
	c.prcsr = processor
	return c
}

func (c *computerVisionGolfServer) ShowDTLPoseImage(ctx context.Context, request *cv.ShowDTLPoseImageRequest) (*cv.ShowDTLPoseImageResponse, error) {
	return c.prcsr.ShowDTLPoseImage(request)
}

func (c *computerVisionGolfServer) ShowFaceOnPoseImage(ctx context.Context, request *cv.ShowFaceOnPoseImageRequest) (*cv.ShowFaceOnPoseImageResponse, error) {
	return c.prcsr.ShowFaceOnPoseImage(request)
}

func (c *computerVisionGolfServer) GetDTLPoseSetupPoints(ctx context.Context, request *cv.GetDTLPoseSetupPointsRequest) (*cv.GetDTLPoseSetupPointsResponse, error) {
	return c.prcsr.GetDTLPoseSetupPoints(request)
}

func (c *computerVisionGolfServer) GetFaceOnPoseSetupPoints(ctx context.Context, request *cv.GetFaceOnPoseSetupPointsRequest) (*cv.GetFaceOnPoseSetupPointsResponse, error) {
	return c.prcsr.GetFaceOnPoseSetupPoints(request)
}

func (c *computerVisionGolfServer) IdentifyDTLPoseDifferences(ctx context.Context, request *cv.IdentifyDTLPoseDifferencesRequest) (*cv.IdentifyDTLPoseDifferencesResponse, error) {
	return c.prcsr.IdentifyDTLPoseDifferences(request)
}

func (c *computerVisionGolfServer) IdentifyFaceOnPoseDifferences(ctx context.Context, request *cv.IdentifyFaceOnPoseDifferencesRequest) (*cv.IdentifyFaceOnPoseDifferencesResponse, error) {
	return c.prcsr.IdentifyFaceOnPoseDifferences(request)
}

func (c *computerVisionGolfServer) ShowDTLPoseImagesFromVideo(stream cv.ComputerVisionGolfService_ShowDTLPoseImagesFromVideoServer) error {
	requests := []*cv.ShowDTLPoseImageRequest{}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		requests = append(requests, in)
	}
	responses, err := c.prcsr.ShowDTLPoseImagesFromVideo(requests)
	if err != nil {
		return err
	}
	for _, response := range responses {
		if err := stream.Send(response); err != nil {
			return err
		}
	}
	return nil
}

func (c *computerVisionGolfServer) ShowFaceOnPoseImagesFromVideo(stream cv.ComputerVisionGolfService_ShowFaceOnPoseImagesFromVideoServer) error {
	requests := []*cv.ShowFaceOnPoseImageRequest{}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		requests = append(requests, in)
	}
	responses, err := c.prcsr.ShowFaceOnPoseImagesFromVideo(requests)
	if err != nil {
		return err
	}
	for _, response := range responses {
		if err := stream.Send(response); err != nil {
			return err
		}
	}
	return nil
}
