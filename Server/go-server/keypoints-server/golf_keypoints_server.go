package keypointsserver

import (
	"context"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

type golfKeypointsServer struct {
	skp.UnimplementedGolfKeypointsServiceServer
	handler skp.GolfKeypointsServiceServer
}

func createNewGolfKeypointsServer(handler skp.GolfKeypointsServiceServer) *golfKeypointsServer {
	g := &golfKeypointsServer{}
	g.handler = handler
	return g
}

func (g *golfKeypointsServer) UploadInputImage(ctx context.Context, request *skp.UploadInputImageRequest) (*skp.UploadInputImageResponse, error) {
	if err := verifyUploadInputImageRequest(request); err != nil {
		return nil, err
	}
	return g.handler.UploadInputImage(ctx, request)
}

func (g *golfKeypointsServer) ListInputImagesForUser(ctx context.Context, request *skp.ListInputImagesForUserRequest) (*skp.ListInputImagesForUserResponse, error) {
	if err := verifyListInputImagesForUserRequest(request); err != nil {
		return nil, err
	}
	return g.handler.ListInputImagesForUser(ctx, request)
}

func (g *golfKeypointsServer) ReadInputImage(ctx context.Context, request *skp.ReadInputImageRequest) (*skp.ReadInputImageResponse, error) {
	if err := verifyReadInputImageRequest(request); err != nil {
		return nil, err
	}
	return g.handler.ReadInputImage(ctx, request)
}

func (g *golfKeypointsServer) DeleteInputImage(ctx context.Context, request *skp.DeleteInputImageRequest) (*skp.DeleteInputImageResponse, error) {
	if err := verifyDeleteInputImageRequest(request); err != nil {
		return nil, err
	}
	return g.handler.DeleteInputImage(ctx, request)
}

func (g *golfKeypointsServer) CalibrateInputImage(ctx context.Context, request *skp.CalibrateInputImageRequest) (*skp.CalibrateInputImageResponse, error) {
	if err := verifyCalibrateInputImageRequest(request); err != nil {
		return nil, err
	}
	return g.handler.CalibrateInputImage(ctx, request)
}

func (g *golfKeypointsServer) CalculateGolfKeypoints(ctx context.Context, request *skp.CalculateGolfKeypointsRequest) (*skp.CalculateGolfKeypointsResponse, error) {
	if err := verifyCalculateGolfKeypointsRequest(request); err != nil {
		return nil, err
	}
	return g.handler.CalculateGolfKeypoints(ctx, request)
}

func (g *golfKeypointsServer) ReadGolfKeypoints(ctx context.Context, request *skp.ReadGolfKeypointsRequest) (*skp.ReadGolfKeypointsResponse, error) {
	if err := verifyReadGolfKeypointsRequest(request); err != nil {
		return nil, err
	}
	return g.handler.ReadGolfKeypoints(ctx, request)
}

func (g *golfKeypointsServer) DeleteGolfKeypoints(ctx context.Context, request *skp.DeleteGolfKeypointsRequest) (*skp.DeleteGolfKeypointsResponse, error) {
	if err := verifyDeleteGolfKeypointsRequest(request); err != nil {
		return nil, err
	}
	return g.handler.DeleteGolfKeypoints(ctx, request)
}

/*func (g *golfKeypointsServer) CreateGolfKeypointsFromVideo(stream skp.GolfKeypointsService_CreateGolfKeypointsFromVideoServer) error {
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
