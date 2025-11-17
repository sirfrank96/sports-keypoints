package keypointsserver

import (
	"context"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
)

type golfKeypointsServer struct {
	cv.UnimplementedGolfKeypointsServiceServer
	handler cv.GolfKeypointsServiceServer
}

func createNewGolfKeypointsServer(handler cv.GolfKeypointsServiceServer) *golfKeypointsServer {
	g := &golfKeypointsServer{}
	g.handler = handler
	return g
}

func (g *golfKeypointsServer) UploadInputImage(ctx context.Context, request *cv.UploadInputImageRequest) (*cv.UploadInputImageResponse, error) {
	if err := verifyUploadInputImageRequest(request); err != nil {
		return nil, err
	}
	return g.handler.UploadInputImage(ctx, request)
}

func (g *golfKeypointsServer) ListInputImagesForUser(ctx context.Context, request *cv.ListInputImagesForUserRequest) (*cv.ListInputImagesForUserResponse, error) {
	if err := verifyListInputImagesForUserRequest(request); err != nil {
		return nil, err
	}
	return g.handler.ListInputImagesForUser(ctx, request)
}

func (g *golfKeypointsServer) ReadInputImage(ctx context.Context, request *cv.ReadInputImageRequest) (*cv.ReadInputImageResponse, error) {
	if err := verifyReadInputImageRequest(request); err != nil {
		return nil, err
	}
	return g.handler.ReadInputImage(ctx, request)
}

func (g *golfKeypointsServer) DeleteInputImage(ctx context.Context, request *cv.DeleteInputImageRequest) (*cv.DeleteInputImageResponse, error) {
	if err := verifyDeleteInputImageRequest(request); err != nil {
		return nil, err
	}
	return g.handler.DeleteInputImage(ctx, request)
}

func (g *golfKeypointsServer) CalibrateInputImage(ctx context.Context, request *cv.CalibrateInputImageRequest) (*cv.CalibrateInputImageResponse, error) {
	if err := verifyCalibrateInputImageRequest(request); err != nil {
		return nil, err
	}
	return g.handler.CalibrateInputImage(ctx, request)
}

func (g *golfKeypointsServer) CalculateGolfKeypoints(ctx context.Context, request *cv.CalculateGolfKeypointsRequest) (*cv.CalculateGolfKeypointsResponse, error) {
	if err := verifyCalculateGolfKeypointsRequest(request); err != nil {
		return nil, err
	}
	return g.handler.CalculateGolfKeypoints(ctx, request)
}

func (g *golfKeypointsServer) ReadGolfKeypoints(ctx context.Context, request *cv.ReadGolfKeypointsRequest) (*cv.ReadGolfKeypointsResponse, error) {
	if err := verifyReadGolfKeypointsRequest(request); err != nil {
		return nil, err
	}
	return g.handler.ReadGolfKeypoints(ctx, request)
}

func (g *golfKeypointsServer) DeleteGolfKeypoints(ctx context.Context, request *cv.DeleteGolfKeypointsRequest) (*cv.DeleteGolfKeypointsResponse, error) {
	if err := verifyDeleteGolfKeypointsRequest(request); err != nil {
		return nil, err
	}
	return g.handler.DeleteGolfKeypoints(ctx, request)
}

/*func (g *golfKeypointsServer) CreateGolfKeypointsFromVideo(stream cv.GolfKeypointsService_CreateGolfKeypointsFromVideoServer) error {
	requests := []*cv.CreateGolfKeypointsRequest{}
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
	responses, err := g.handler.CreateGolfKeypointsFromVideo(requests)
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
*/
