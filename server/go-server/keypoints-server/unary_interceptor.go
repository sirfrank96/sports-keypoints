package keypointsserver

import (
	"context"
	"fmt"
	"log"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"

	"google.golang.org/grpc"
)

func sessionUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("unary interceptor: ", info.FullMethod)
	switch info.FullMethod {
	case "/sports_keypoints_proto.UserService/CreateUser":
		// no-op for now
	case "/sports_keypoints_proto.UserService/RegisterUser":
		// no-op for now
	case "/sports_keypoints_proto.UserService/ReadUser":
		userId, err := getUserIdFromSessionToken(req.(*skp.ReadUserRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.UserService/UpdateUser":
		userId, err := getUserIdFromSessionToken(req.(*skp.UpdateUserRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.UserService/DeleteUser":
		userId, err := getUserIdFromSessionToken(req.(*skp.DeleteUserRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/UploadInputImage":
		userId, err := getUserIdFromSessionToken(req.(*skp.UploadInputImageRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/ListInputImagesForUser":
		userId, err := getUserIdFromSessionToken(req.(*skp.ListInputImagesForUserRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/ReadInputImage":
		userId, err := getUserIdFromSessionToken(req.(*skp.ReadInputImageRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/DeleteInputImage":
		userId, err := getUserIdFromSessionToken(req.(*skp.DeleteInputImageRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/CalibrateInputImage":
		userId, err := getUserIdFromSessionToken(req.(*skp.CalibrateInputImageRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/CalculateGolfKeypoints":
		userId, err := getUserIdFromSessionToken(req.(*skp.CalculateGolfKeypointsRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/ReadGolfKeypoints":
		userId, err := getUserIdFromSessionToken(req.(*skp.ReadGolfKeypointsRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/UpdateBodyKeypoints":
		userId, err := getUserIdFromSessionToken(req.(*skp.UpdateBodyKeypointsRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/DeleteGolfKeypoints":
		userId, err := getUserIdFromSessionToken(req.(*skp.DeleteGolfKeypointsRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	}
	return handler(ctx, req)
}

func getUserIdFromSessionToken(sessionToken string) (string, error) {
	if sessionToken == "" {
		return "", fmt.Errorf("no session token provided")
	}
	claims, err := util.VerifyJWTSessionToken(sessionToken)
	if err != nil {
		return "", err
	}
	userId, err := util.GetUserIdFromClaims(claims)
	if err != nil {
		return "", err
	}
	return userId, nil
}
