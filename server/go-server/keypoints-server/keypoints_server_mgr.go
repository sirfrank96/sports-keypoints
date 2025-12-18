// Receive an image/video via client rest apis
// (CRUD REST APIs for users and also gRPC apis for images/videos and more specific bidirectional streaming client apis related to computervision (ie. identifyPoseDifferences, showPose, etc.))
// Converts received data into objects for server_mgr to process and figure out what needs to be done and what needs computervision stuff
package keypointsserver

import (
	"flag"
	"fmt"
	"log"
	"net"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
)

var (
	//tls = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	//certFile = flag.String("cert_file", "", "The TLS cert file")
	//keyFile = flag.String("key_file", "", "The TLS key file")
	port = flag.Int("port", 50052, "The server port")
)

type KeypointsServerManager struct {
	grpcServer          *grpc.Server
	userServer          *userServer
	golfKeypointsServer *golfKeypointsServer
}

func NewKeypointsServerManager(golfKeypointsHandler skp.GolfKeypointsServiceServer, userHandler skp.UserServiceServer) *KeypointsServerManager {
	k := &KeypointsServerManager{}
	k.userServer = createNewUserServer(userHandler)
	k.golfKeypointsServer = createNewGolfKeypointsServer(golfKeypointsHandler)
	log.Printf("New keypoints_server_mgr")
	return k
}

func (k *KeypointsServerManager) StartKeypointsServer() error {
	log.Printf("Starting keypoints and user servers")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	//var opts []grpc.ServerOption
	k.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(sessionUnaryInterceptor))
	skp.RegisterGolfKeypointsServiceServer(k.grpcServer, k.golfKeypointsServer)
	skp.RegisterUserServiceServer(k.grpcServer, k.userServer)
	k.grpcServer.Serve(lis)
	return nil
}

func (k *KeypointsServerManager) StopKeypointsServer() error {
	k.grpcServer.GracefulStop()
	return nil
}
