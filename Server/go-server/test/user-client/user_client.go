package uclient

import (
	"context"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	//testutil "github.com/sirfrank96/go-server/test/test-util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//TODO: RETURN ERRORS INSTEAD OF LOG.FATALF

// Middle arg is a close function, should be called by calling function
func InitUserServiceGrpcClient(serveraddr string) (skp.UserServiceClient, func() error, error) {
	// Set up a connection to the cv_api server.
	conn, err := grpc.NewClient(serveraddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, conn.Close, err
	}
	// Init ComputerVisionGolf grpc client
	uclient := skp.NewUserServiceClient(conn)
	return uclient, conn.Close, nil
}

func CreateUser(ctx context.Context, uclient skp.UserServiceClient, userName string, password string, email string) (*skp.CreateUserResponse, error) {
	request := &skp.CreateUserRequest{
		UserName: userName,
		Password: password,
		Email:    email,
	}
	return uclient.CreateUser(ctx, request)
}

func RegisterUser(ctx context.Context, uclient skp.UserServiceClient, userName string, password string) (*skp.RegisterUserResponse, error) {
	request := &skp.RegisterUserRequest{
		UserName: userName,
		Password: password,
	}
	return uclient.RegisterUser(ctx, request)
}
