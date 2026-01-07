package controller

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	cvclient "github.com/sirfrank96/go-server/cv-client"
	db "github.com/sirfrank96/go-server/db"
	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

type GolfKeypointsListener struct {
	skp.UnimplementedGolfKeypointsServiceServer
	cvmgr *cvclient.CvClientManager
	dbmgr *db.DbManager
}

func newGolfKeypointsListener(cvmgr *cvclient.CvClientManager, dbmgr *db.DbManager) *GolfKeypointsListener {
	return &GolfKeypointsListener{
		cvmgr: cvmgr,
		dbmgr: dbmgr,
	}
}

func (g *GolfKeypointsListener) UploadInputImage(ctx context.Context, request *skp.UploadInputImageRequest) (*skp.UploadInputImageResponse, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("unable to get userId from context")
	}
	if _, err := verifyUserExists(ctx, g.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	// put image into db
	inputImage := &db.InputImage{
		UserId:          userId,
		ImageType:       request.ImageType,
		InputImg:        request.Image,
		Description:     request.Description,
		Timestamp:       request.Timestamp.AsTime(), // UTC
		CalibrationInfo: *util.GetEmptyCalibrationInfo(),
	}
	inputImage, err := g.dbmgr.CreateInputImage(ctx, inputImage)
	if err != nil {
		return nil, fmt.Errorf("could not store input image: %w", err)
	}
	// return response
	response := &skp.UploadInputImageResponse{
		Success:      true,
		InputImageId: inputImage.Id.Hex(),
	}
	return response, nil
}

func (g *GolfKeypointsListener) ListInputImagesForUser(ctx context.Context, request *skp.ListInputImagesForUserRequest) (*skp.ListInputImagesForUserResponse, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if _, err := verifyUserExists(ctx, g.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	// get all input images for userid from db
	inputImgs, err := g.dbmgr.ReadInputImagesForUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("could not get images for user from db: %w", err)
	}
	var inputImgIds []string
	for _, inputImg := range inputImgs {
		inputImgIds = append(inputImgIds, inputImg.Id.Hex())
	}
	// return response
	response := &skp.ListInputImagesForUserResponse{
		Success:       true,
		InputImageIds: inputImgIds,
	}
	return response, nil
}

func (g *GolfKeypointsListener) ReadInputImage(ctx context.Context, request *skp.ReadInputImageRequest) (*skp.ReadInputImageResponse, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if _, err := verifyUserExists(ctx, g.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	// get inputimg with inputimgid from db
	inputImg, err := g.dbmgr.ReadInputImage(ctx, request.InputImageId)
	if err != nil {
		return nil, fmt.Errorf("could not read input image with id: %s: %w", request.InputImageId, err)
	}
	// return response
	response := &skp.ReadInputImageResponse{
		Success:         true,
		ImageType:       inputImg.ImageType,
		Image:           inputImg.InputImg,
		CalibrationType: inputImg.CalibrationInfo.CalibrationType,
		FeetLineMethod:  inputImg.CalibrationInfo.FeetLineMethod,
		Description:     inputImg.Description,
		Timestamp:       timestamppb.New(inputImg.Timestamp),
	}
	return response, nil
}

func (g *GolfKeypointsListener) DeleteInputImage(ctx context.Context, request *skp.DeleteInputImageRequest) (*skp.DeleteInputImageResponse, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if _, err := verifyUserExists(ctx, g.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	// delete inputimg with inputimgid in db
	err := g.dbmgr.DeleteInputImage(ctx, request.InputImageId)
	if err != nil {
		return nil, fmt.Errorf("could not delete input image with id: %s: %w", request.InputImageId, err)
	}
	// return response
	response := &skp.DeleteInputImageResponse{
		Success: true,
	}
	return response, nil
}

func (g *GolfKeypointsListener) CalibrateInputImage(ctx context.Context, request *skp.CalibrateInputImageRequest) (*skp.CalibrateInputImageResponse, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if _, err := verifyUserExists(ctx, g.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	calibrationInfo := &util.CalibrationInfo{
		ManualGenerated:              false,
		CalibrationImgAxes:           request.CalibrationImageAxes,
		CalibrationImgVanishingPoint: request.CalibrationImageVanishingPoint,
		CalibrationType:              request.CalibrationType,
		FeetLineMethod:               request.FeetLineMethod,
	}
	return g.calibrateInputImageHelper(ctx, calibrationInfo, request.InputImageId, request.GolfBall, request.ClubButt, request.ClubHead, nil, nil, request.ShoulderTilt, nil, nil)
}

func (g *GolfKeypointsListener) CalibrateInputImageManual(ctx context.Context, request *skp.CalibrateInputImageManualRequest) (*skp.CalibrateInputImageResponse, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if _, err := verifyUserExists(ctx, g.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	calibrationInfo := &util.CalibrationInfo{
		ManualGenerated: true,
		CalibrationType: request.CalibrationType,
		FeetLineMethod:  request.FeetLineMethod,
	}
	return g.calibrateInputImageHelper(ctx, calibrationInfo, request.InputImageId, request.GolfBall, request.ClubButt, request.ClubHead, request.HorizontalAxis, request.VerticalAxis, request.ShoulderTilt, request.FirstLineAtTarget, request.SecondLineAtTarget)
}

func (g *GolfKeypointsListener) CalculateGolfKeypoints(ctx context.Context, request *skp.CalculateGolfKeypointsRequest) (*skp.CalculateGolfKeypointsResponse, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if _, err := verifyUserExists(ctx, g.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	// get inputimage from db
	inputImage, err := g.dbmgr.ReadInputImage(ctx, request.InputImageId)
	if err != nil {
		return nil, fmt.Errorf("could not get input image with id: %s, error was %w", request.InputImageId, err)
	}
	// get pose image and data for input img
	getPoseAllResponse, err := g.cvmgr.GetPoseAll(inputImage.InputImg)
	if err != nil {
		return nil, fmt.Errorf("could not get pose all for image: %w", err)
	}
	// calculate golf setup points
	golfKeypoints := &db.GolfKeypoints{
		UserId:          userId,
		InputImageId:    request.InputImageId,
		OutputImg:       getPoseAllResponse.Image,
		OutputKeypoints: *getPoseAllResponse.PoseKeypoints,
	}
	// dtl setup points
	if inputImage.ImageType == skp.ImageType_DTL {
		golfKeypoints.DtlGolfSetupPoints = *CalculateDTLSetupPoints(ctx, getPoseAllResponse.PoseKeypoints, &inputImage.CalibrationInfo)
	} else { // face on setup points
		golfKeypoints.FaceonGolfSetupPoints = *CalculateFaceOnSetupPoints(ctx, getPoseAllResponse.PoseKeypoints, &inputImage.CalibrationInfo)
	}

	// store golfkeypoints in db
	_, err = g.dbmgr.CreateGolfKeypoints(ctx, golfKeypoints)
	if err != nil {
		return nil, fmt.Errorf("could not store golfkeypoints in db %w", err)
	}

	// return response
	response := &skp.CalculateGolfKeypointsResponse{
		Success:       true,
		OutputImage:   getPoseAllResponse.Image,
		GolfKeypoints: db.ConvertGolfKeypointsToCVGolfKeypoints(golfKeypoints),
	}
	return response, nil
}

func (g *GolfKeypointsListener) ReadGolfKeypoints(ctx context.Context, request *skp.ReadGolfKeypointsRequest) (*skp.ReadGolfKeypointsResponse, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if _, err := verifyUserExists(ctx, g.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	// find golf keypoints for associated input image id in db
	golfKeypoints, err := g.dbmgr.ReadGolfKeypointsForInputImage(ctx, request.InputImageId)
	if err != nil {
		return nil, fmt.Errorf("could not read golf keypoints from db for input image: %s, %w", request.InputImageId, err)
	}
	// return response
	response := &skp.ReadGolfKeypointsResponse{
		Success:       true,
		OutputImage:   golfKeypoints.OutputImg,
		GolfKeypoints: db.ConvertGolfKeypointsToCVGolfKeypoints(golfKeypoints),
	}
	return response, nil
}

func (g *GolfKeypointsListener) UpdateBodyKeypoints(ctx context.Context, request *skp.UpdateBodyKeypointsRequest) (*skp.UpdateBodyKeypointsResponse, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if _, err := verifyUserExists(ctx, g.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	// get inputimage from db
	inputImage, err := g.dbmgr.ReadInputImage(ctx, request.InputImageId)
	if err != nil {
		return nil, fmt.Errorf("could not get input image with id: %s, error was %w", request.InputImageId, err)
	}
	// get current golf keypoints for associated input image id in db
	golfKeypoints, err := g.dbmgr.ReadGolfKeypointsForInputImage(ctx, request.InputImageId)
	if err != nil {
		return nil, fmt.Errorf("could not read golf keypoints from db for input image: %s, %w", request.InputImageId, err)
	}
	golfKeypoints.OutputKeypoints = *db.UpdateOutputKeypointsFields(&golfKeypoints.OutputKeypoints, request.UpdatedBodyKeypoints)
	// recalculate golf setup points based on new keypoints
	// dtl setup points
	if inputImage.ImageType == skp.ImageType_DTL {
		golfKeypoints.DtlGolfSetupPoints = *CalculateDTLSetupPoints(ctx, &golfKeypoints.OutputKeypoints, &inputImage.CalibrationInfo)
	} else { // face on setup points
		golfKeypoints.FaceonGolfSetupPoints = *CalculateFaceOnSetupPoints(ctx, &golfKeypoints.OutputKeypoints, &inputImage.CalibrationInfo)
	}
	// update new golf keypoints in db
	updatedGolfKeypoints, err := g.dbmgr.UpdateGolfKeypointsForInputImage(ctx, request.InputImageId, golfKeypoints)
	if err != nil {
		return nil, fmt.Errorf("could not update golf keypoints from db for input image: %s, %w", request.InputImageId, err)
	}
	// return response
	response := &skp.UpdateBodyKeypointsResponse{
		Success:              true,
		UpdatedGolfKeypoints: db.ConvertGolfKeypointsToCVGolfKeypoints(updatedGolfKeypoints),
	}
	return response, nil
}

func (g *GolfKeypointsListener) DeleteGolfKeypoints(ctx context.Context, request *skp.DeleteGolfKeypointsRequest) (*skp.DeleteGolfKeypointsResponse, error) {
	// make sure user exists
	userId, ok := ctx.Value(util.UserIdKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if _, err := verifyUserExists(ctx, g.dbmgr, userId); err != nil {
		return nil, fmt.Errorf("could not verify user exists")
	}
	// delete golf keypoints for associated input image id in db
	err := g.dbmgr.DeleteGolfKeypointsForInputImage(ctx, request.InputImageId)
	if err != nil {
		return nil, fmt.Errorf("could not delete golf keypoints from db for input image: %s, %w", request.InputImageId, err)
	}
	// return response
	response := &skp.DeleteGolfKeypointsResponse{
		Success: true,
	}
	return response, nil
}
