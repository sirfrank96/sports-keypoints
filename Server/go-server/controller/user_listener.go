package controller

import (
	"context"
	"fmt"

	cvclient "github.com/sirfrank96/go-server/cv-client"
	db "github.com/sirfrank96/go-server/db"
	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

type UserListener struct {
	skp.UnimplementedUserServiceServer
	cvmgr *cvclient.CvClientManager
	dbmgr *db.DbManager
}

func newUserListener(cvmgr *cvclient.CvClientManager, dbmgr *db.DbManager) *UserListener {
	return &UserListener{
		cvmgr: cvmgr,
		dbmgr: dbmgr,
	}
}

func (u *UserListener) CreateUser(ctx context.Context, request *skp.CreateUserRequest) (*skp.CreateUserResponse, error) {
	// check if user exists already
	user, err := u.dbmgr.ReadUserFromUsername(ctx, request.UserName)
	if user != nil && err == nil {
		return nil, fmt.Errorf("user with username: %s already exists", request.UserName)
	}
	hashedPassword, err := db.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}
	user = &db.User{
		Username: request.UserName,
		Password: hashedPassword,
		Email:    request.Email,
	}
	_, err = u.dbmgr.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("could not store user in db: %w", err)
	}
	response := &skp.CreateUserResponse{
		Success: true,
	}
	return response, nil
}

func (u *UserListener) RegisterUser(ctx context.Context, request *skp.RegisterUserRequest) (*skp.RegisterUserResponse, error) {
	user, err := u.dbmgr.ReadUserFromUsername(ctx, request.UserName)
	if err != nil {
		return nil, fmt.Errorf("could not fine user with username: %s, error: %w", request.UserName, err)
	}
	if !db.VerifyPasswordHash(user.Password, request.Password) {
		return nil, fmt.Errorf("passwords do not match, could not register user")
	}
	sessionToken, err := util.CreateJWTSessionToken(user.Id.Hex())
	if err != nil {
		return nil, fmt.Errorf("could not create session token from id: %w", err)
	}
	response := &skp.RegisterUserResponse{
		Success:      true,
		SessionToken: sessionToken,
	}
	return response, nil
}

func (u *UserListener) ReadUser(ctx context.Context, request *skp.ReadUserRequest) (*skp.User, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if _, err := verifyUserExists(ctx, u.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	// find user with associated user id in db
	user, err := u.dbmgr.ReadUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %w", err)
	}
	// return response
	response := &skp.User{
		UserName: user.Username,
		Email:    user.Email,
	}
	return response, nil
}

func (u *UserListener) UpdateUser(ctx context.Context, request *skp.UpdateUserRequest) (*skp.User, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if _, err := verifyUserExists(ctx, u.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	// find user with associated user id in db
	currUser, err := u.dbmgr.ReadUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %w", err)
	}
	// update user with new fields
	newPassword := ""
	if request.Password != "" {
		newPassword, err = db.HashPassword(request.Password)
		if err != nil {
			return nil, fmt.Errorf("could not hash new password")
		}
	}
	newUser := &db.User{Username: request.UserName, Password: newPassword, Email: request.Email}
	updatedFieldsUser := db.UpdateUserFields(currUser, newUser)
	updatedUser, err := u.dbmgr.UpdateUser(ctx, userId, updatedFieldsUser)
	if err != nil {
		return nil, fmt.Errorf("could not update user in db: %w", err)
	}
	// return response
	response := &skp.User{
		UserName: updatedUser.Username,
		Email:    updatedUser.Email,
	}
	return response, nil
}

func (u *UserListener) DeleteUser(ctx context.Context, request *skp.DeleteUserRequest) (*skp.DeleteUserResponse, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if _, err := verifyUserExists(ctx, u.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	// delete user with associated user id in db
	err := u.dbmgr.DeleteUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("could not delete user: %w", err)
	}
	// return response
	response := &skp.DeleteUserResponse{
		Success: true,
	}
	return response, nil
}
