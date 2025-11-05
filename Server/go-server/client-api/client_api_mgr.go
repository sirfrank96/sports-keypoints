// Receive an image/video via client rest apis
// (CRUD REST APIs for users and also gRPC apis for images/videos and more specific bidirectional streaming client apis related to computervision (ie. identifyPoseDifferences, showPose, etc.))
// Converts received data into objects for server_mgr to process and figure out what needs to be done and what needs opencv/openpose stuff

package clientapi

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	//"sync"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	processor "github.com/sirfrank96/go-server/processor"

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

type ClientApiManager struct {
	ctx context.Context
	//mutex      sync.Mutex
	prcsr      *processor.Processor
	grpcServer *grpc.Server
}

func NewClientApiManager(ctx context.Context) *ClientApiManager {
	c := &ClientApiManager{}
	c.ctx = ctx
	c.prcsr = processor.NewProcessor()
	log.Printf("Starting client_api_mgr")
	return c
}

func (c *ClientApiManager) StartGolfComputerVisionServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	c.grpcServer = grpc.NewServer(opts...)
	cv.RegisterComputerVisionGolfServiceServer(c.grpcServer, createNewComputerVisionGolfServer(c.ctx, c.prcsr))
	c.grpcServer.Serve(lis)
}

func (c *ClientApiManager) StartOpenCvApiClient() {
	c.prcsr.StartOpenCvApiClient()
}

func (c *ClientApiManager) StopGolfComputerVisionServer() {
	c.grpcServer.GracefulStop()
}

func (c *ClientApiManager) CloseOpenCvApiClient() {
	c.prcsr.CloseOpenCvApiClient()
}
