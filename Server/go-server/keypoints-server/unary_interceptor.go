package keypointsserver

import (
	"context"
	"log"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"

	"google.golang.org/grpc"
)

func sessionUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("unary interceptor: ", info.FullMethod)
	switch info.FullMethod {
	case "/sports_keypoints_proto.UserService/createUser":
		// no-op for now
	case "/sports_keypoints_proto.UserService/registerUser":
		// no-op for now
	case "/sports_keypoints_proto.UserService/readUser":
		userId, err := getUserIdFromSessionToken(req.(*skp.ReadUserRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.UserService/updateUser":
		userId, err := getUserIdFromSessionToken(req.(*skp.UpdateUserRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.UserService/deleteUser":
		userId, err := getUserIdFromSessionToken(req.(*skp.DeleteUserRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/uploadInputImage":
		userId, err := getUserIdFromSessionToken(req.(*skp.UploadInputImageRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		log.Printf("Userid from sessiontoken is: %s", userId)
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/listInputImagesForUser":
		userId, err := getUserIdFromSessionToken(req.(*skp.ListInputImagesForUserRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/readInputImage":
		userId, err := getUserIdFromSessionToken(req.(*skp.ReadInputImageRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/deleteInputImage":
		userId, err := getUserIdFromSessionToken(req.(*skp.DeleteInputImageRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/calibrateInputImage":
		userId, err := getUserIdFromSessionToken(req.(*skp.CalibrateInputImageRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/calculateGolfKeypoints":
		userId, err := getUserIdFromSessionToken(req.(*skp.CalculateGolfKeypointsRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/readGolfKeypoints":
		userId, err := getUserIdFromSessionToken(req.(*skp.ReadGolfKeypointsRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	case "/sports_keypoints_proto.GolfKeypointsService/deleteGolfKeypoints":
		userId, err := getUserIdFromSessionToken(req.(*skp.DeleteGolfKeypointsRequest).SessionToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, util.UserIdKey, userId)
	}
	return handler(ctx, req)
}

func getUserIdFromSessionToken(sessionToken string) (string, error) {
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
