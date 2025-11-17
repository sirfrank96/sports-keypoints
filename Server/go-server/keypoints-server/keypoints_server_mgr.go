// Receive an image/video via client rest apis
// (CRUD REST APIs for users and also gRPC apis for images/videos and more specific bidirectional streaming client apis related to computervision (ie. identifyPoseDifferences, showPose, etc.))
// Converts received data into objects for server_mgr to process and figure out what needs to be done and what needs opencv/openpose stuff
package keypointsserver

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
)

var (
	//tls = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	//certFile = flag.String("cert_file", "", "The TLS cert file")
	//keyFile = flag.String("key_file", "", "The TLS key file")
	//jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port = flag.Int("port", 50052, "The server port")
)

type KeypointsServerManager struct {
	grpcServer          *grpc.Server
	userServer          *userServer
	golfKeypointsServer *golfKeypointsServer
}

func NewKeypointsServerManager(golfKeypointsHandler cv.GolfKeypointsServiceServer, userHandler cv.UserServiceServer) *KeypointsServerManager {
	k := &KeypointsServerManager{}
	k.userServer = createNewUserServer(userHandler)
	k.golfKeypointsServer = createNewGolfKeypointsServer(golfKeypointsHandler)
	log.Printf("New keypoints_server_mgr")
	return k
}

func (k *KeypointsServerManager) StartKeypointsServer() error {
	log.Printf("Starting keypoints and user servers")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	//var opts []grpc.ServerOption
	k.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(sessionUnaryInterceptor))
	cv.RegisterGolfKeypointsServiceServer(k.grpcServer, k.golfKeypointsServer)
	cv.RegisterUserServiceServer(k.grpcServer, k.userServer)
	k.grpcServer.Serve(lis)
	return nil
}

func (k *KeypointsServerManager) StopKeypointsServer() error {
	k.grpcServer.GracefulStop()
	return nil
}

func sessionUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("unary interceptor: ", info.FullMethod)
	switch info.FullMethod {
	case "/computer_vision_sports_proto.GolfKeypointsService/uploadInputImage":
		// TODO: pull userid out of jwt, put userid key into util
		ctx = context.WithValue(ctx, "userid", req.(*cv.UploadInputImageRequest).SessionToken)
	case "/computer_vision_sports_proto.GolfKeypointsService/calibrateInputImage":
		// TODO: pull userid out of jwt
		ctx = context.WithValue(ctx, "userid", req.(*cv.CalibrateInputImageRequest).SessionToken)
	case "/computer_vision_sports_proto.GolfKeypointsService/calculateGolfKeypoints":
		// TODO: pull userid out of jwt
		ctx = context.WithValue(ctx, "userid", req.(*cv.CalculateGolfKeypointsRequest).SessionToken)
	}
	return handler(ctx, req)
}
