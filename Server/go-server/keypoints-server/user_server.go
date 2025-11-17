package keypointsserver

import (
	"context"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

type userServer struct {
	skp.UnimplementedUserServiceServer
	handler skp.UserServiceServer
}

func createNewUserServer(handler skp.UserServiceServer) *userServer {
	u := &userServer{}
	u.handler = handler
	return u
}

func (u *userServer) CreateUser(ctx context.Context, request *skp.CreateUserRequest) (*skp.CreateUserResponse, error) {
	if err := verifyCreateUserRequest(request); err != nil {
		return nil, err
	}
	return u.handler.CreateUser(ctx, request)
}

func (u *userServer) RegisterUser(ctx context.Context, request *skp.RegisterUserRequest) (*skp.RegisterUserResponse, error) {
	if err := verifyRegisterUserRequest(request); err != nil {
		return nil, err
	}
	return u.handler.RegisterUser(ctx, request)
}

func (u *userServer) ReadUser(ctx context.Context, request *skp.ReadUserRequest) (*skp.User, error) {
	if err := verifyReadUserRequest(request); err != nil {
		return nil, err
	}
	return u.handler.ReadUser(ctx, request)
}

func (u *userServer) UpdateUser(ctx context.Context, request *skp.UpdateUserRequest) (*skp.User, error) {
	if err := verifyUpdateUserRequest(request); err != nil {
		return nil, err
	}
	return u.handler.UpdateUser(ctx, request)
}

func (u *userServer) DeleteUser(ctx context.Context, request *skp.DeleteUserRequest) (*skp.DeleteUserResponse, error) {
	if err := verifyDeleteUserRequest(request); err != nil {
		return nil, err
	}
	return u.handler.DeleteUser(ctx, request)
}
