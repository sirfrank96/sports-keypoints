package keypointsserver

import (
	"context"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
)

type userServer struct {
	cv.UnimplementedUserServiceServer
	handler cv.UserServiceServer
}

func createNewUserServer(handler cv.UserServiceServer) *userServer {
	u := &userServer{}
	u.handler = handler
	return u
}

func (u *userServer) CreateUser(ctx context.Context, request *cv.CreateUserRequest) (*cv.CreateUserResponse, error) {
	if err := verifyCreateUserRequest(request); err != nil {
		return nil, err
	}
	return u.handler.CreateUser(ctx, request)
}

func (u *userServer) RegisterUser(ctx context.Context, request *cv.RegisterUserRequest) (*cv.RegisterUserResponse, error) {
	if err := verifyRegisterUserRequest(request); err != nil {
		return nil, err
	}
	return u.handler.RegisterUser(ctx, request)
}

func (u *userServer) ReadUser(ctx context.Context, request *cv.ReadUserRequest) (*cv.User, error) {
	if err := verifyReadUserRequest(request); err != nil {
		return nil, err
	}
	return u.handler.ReadUser(ctx, request)
}

func (u *userServer) UpdateUser(ctx context.Context, request *cv.UpdateUserRequest) (*cv.User, error) {
	if err := verifyUpdateUserRequest(request); err != nil {
		return nil, err
	}
	return u.handler.UpdateUser(ctx, request)
}

func (u *userServer) DeleteUser(ctx context.Context, request *cv.DeleteUserRequest) (*cv.DeleteUserResponse, error) {
	if err := verifyDeleteUserRequest(request); err != nil {
		return nil, err
	}
	return u.handler.DeleteUser(ctx, request)
}
