package controller

import (
	"context"
	"fmt"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	db "github.com/sirfrank96/go-server/db"
	opencvclient "github.com/sirfrank96/go-server/opencv-client"
)

type UserListener struct {
	cv.UnimplementedUserServiceServer
	ocvmgr *opencvclient.OpenCvClientManager
	dbmgr  *db.DbManager
}

func newUserListener(ocvmgr *opencvclient.OpenCvClientManager, dbmgr *db.DbManager) *UserListener {
	return &UserListener{
		ocvmgr: ocvmgr,
		dbmgr:  dbmgr,
	}
}

func (u *UserListener) CreateUser(ctx context.Context, request *cv.CreateUserRequest) (*cv.CreateUserResponse, error) {
	hashedPassword, err := db.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}
	user := &db.User{
		Username: request.UserName,
		Password: hashedPassword,
		Email:    request.Email,
	}
	_, err = u.dbmgr.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("could not store user in db: %w", err)
	}
	response := &cv.CreateUserResponse{
		Success: true,
	}
	return response, nil
}

func (u *UserListener) RegisterUser(ctx context.Context, request *cv.RegisterUserRequest) (*cv.RegisterUserResponse, error) {
	user, err := u.dbmgr.ReadUserFromUsername(ctx, request.UserName)
	if err != nil {
		return nil, fmt.Errorf("could not fine user with username: %s, error: %w", request.UserName, err)
	}
	if !db.VerifyPasswordHash(user.Password, request.Password) {
		return nil, fmt.Errorf("passwords do not match, could not register user")
	}
	// TODO: IMPLEMENT JWT using unique user id
	response := &cv.RegisterUserResponse{
		Success:      true,
		SessionToken: user.Id.Hex(),
	}
	return response, nil
}

func (u *UserListener) ReadUser(ctx context.Context, request *cv.ReadUserRequest) (*cv.User, error) {
	return nil, nil
}

func (u *UserListener) UpdateUser(ctx context.Context, request *cv.UpdateUserRequest) (*cv.User, error) {
	return nil, nil
}

func (u *UserListener) DeleteUser(ctx context.Context, request *cv.DeleteUserRequest) (*cv.DeleteUserResponse, error) {
	return nil, nil
}

// implements cv.UnimplementedUserServiceServer
