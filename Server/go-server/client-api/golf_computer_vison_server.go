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
	response := c.prcsr.ShowDTLPoseImage(request)
	return response, nil
}

func (c *computerVisionGolfServer) ShowFaceOnPoseImage(ctx context.Context, request *cv.ShowFaceOnPoseImageRequest) (*cv.ShowFaceOnPoseImageResponse, error) {
	response := c.prcsr.ShowFaceOnPoseImage(request)
	return response, nil
}

func (c *computerVisionGolfServer) IdentifyDTLPoseDifferences(ctx context.Context, request *cv.IdentifyDTLPoseDifferencesRequest) (*cv.IdentifyDTLPoseDifferencesResponse, error) {
	return nil, nil
}

func (c *computerVisionGolfServer) IdentifyFaceOnPoseDifferences(ctx context.Context, request *cv.IdentifyFaceOnPoseDifferencesRequest) (*cv.IdentifyFaceOnPoseDifferencesResponse, error) {
	return nil, nil
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
	responses := c.prcsr.ShowDTLPoseImagesFromVideo(requests)
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
			return nil
		}
		if err != nil {
			return err
		}
		requests = append(requests, in)
	}
	responses := c.prcsr.ShowFaceOnPoseImagesFromVideo(requests)
	for _, response := range responses {
		if err := stream.Send(response); err != nil {
			return err
		}
	}
	return nil
}
