package controller

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	cvclient "github.com/sirfrank96/go-server/cv-client"
	db "github.com/sirfrank96/go-server/db"
	"github.com/sirfrank96/go-server/draw"
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
		Calibrated:      false,
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
		Success:     true,
		ImageType:   inputImg.ImageType,
		Image:       inputImg.InputImg,
		Description: inputImg.Description,
		Calibrated:  inputImg.Calibrated,
		Timestamp:   timestamppb.New(inputImg.Timestamp),
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
	return g.calibrateInputImageHelper(ctx, calibrationInfo, request.InputImageId, nil, nil, nil, nil)
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
	return g.calibrateInputImageHelper(ctx, calibrationInfo, request.InputImageId, request.HorizontalAxis, request.VerticalAxis, request.FirstLineAtTarget, request.SecondLineAtTarget)
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
	// get pose data for input img
	getPoseDataResponse, err := g.cvmgr.GetPoseData(inputImage.InputImg)
	if err != nil {
		return nil, fmt.Errorf("could not get pose data for image: %w", err)
	}
	// init GolfKeypoints obj to be stored in db
	golfKeypoints := &db.GolfKeypoints{
		UserId:         userId,
		InputImageId:   request.InputImageId,
		BodyDatapoints: *getPoseDataResponse.Datapoints,
	}
	// put in golf specific data points
	if request.GolfSpecificDatapoints != nil {
		golfKeypoints.GolfSpecificDatapoints = *request.GolfSpecificDatapoints
	}
	// dtl setup points
	if inputImage.ImageType == skp.ImageType_DTL {
		golfKeypoints.DtlGolfSetupPoints = *CalculateDTLSetupPoints(ctx, getPoseDataResponse.Datapoints, request.GolfSpecificDatapoints, &inputImage.CalibrationInfo)
	} else { // face on setup points
		golfKeypoints.FaceonGolfSetupPoints = *CalculateFaceOnSetupPoints(ctx, getPoseDataResponse.Datapoints, request.GolfSpecificDatapoints, &inputImage.CalibrationInfo)
	}
	// draw datapoints with skeleton on image
	outputImg, err := draw.DrawGolfSkeleton(ctx, inputImage.InputImg, getPoseDataResponse.Datapoints, request.GolfSpecificDatapoints)
	if err != nil {
		return nil, fmt.Errorf("could not draw golf skeleton on image: %v", err)
	}
	golfKeypoints.OutputImg = outputImg
	// store golfkeypoints in db
	_, err = g.dbmgr.CreateGolfKeypoints(ctx, golfKeypoints)
	if err != nil {
		return nil, fmt.Errorf("could not store golfkeypoints in db %v", err)
	}
	// return response
	response := &skp.CalculateGolfKeypointsResponse{
		Success:       true,
		OutputImage:   outputImg,
		GolfKeypoints: db.ConvertGolfKeypointsToSkpGolfKeypoints(golfKeypoints),
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
		GolfKeypoints: db.ConvertGolfKeypointsToSkpGolfKeypoints(golfKeypoints),
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

func (g *GolfKeypointsListener) UpdateBodyDatapoints(ctx context.Context, request *skp.UpdateBodyDatapointsRequest) (*skp.UpdateBodyDatapointsResponse, error) {
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
	golfKeypoints.BodyDatapoints = *db.UpdateBodyDatapointsFields(&golfKeypoints.BodyDatapoints, request.UpdatedBodyDatapoints)
	// recalculate golf setup points based on new datapoints
	// dtl setup points
	if inputImage.ImageType == skp.ImageType_DTL {
		golfKeypoints.DtlGolfSetupPoints = *CalculateDTLSetupPoints(ctx, &golfKeypoints.BodyDatapoints, &golfKeypoints.GolfSpecificDatapoints, &inputImage.CalibrationInfo)
	} else { // face on setup points
		golfKeypoints.FaceonGolfSetupPoints = *CalculateFaceOnSetupPoints(ctx, &golfKeypoints.BodyDatapoints, &golfKeypoints.GolfSpecificDatapoints, &inputImage.CalibrationInfo)
	}
	// redraw skeleton based on new datapoints
	updatedOutputImg, err := draw.DrawGolfSkeleton(ctx, inputImage.InputImg, &golfKeypoints.BodyDatapoints, &golfKeypoints.GolfSpecificDatapoints)
	if err != nil {
		return nil, fmt.Errorf("could not redraw golf skeleton on image: %v", err)
	}
	golfKeypoints.OutputImg = updatedOutputImg
	// update new golf keypoints in db
	updatedGolfKeypoints, err := g.dbmgr.UpdateGolfKeypointsForInputImage(ctx, request.InputImageId, golfKeypoints)
	if err != nil {
		return nil, fmt.Errorf("could not update golf keypoints from db for input image: %s, %w", request.InputImageId, err)
	}
	// return response
	response := &skp.UpdateBodyDatapointsResponse{
		Success:              true,
		UpdatedOutputImage:   updatedOutputImg,
		UpdatedGolfKeypoints: db.ConvertGolfKeypointsToSkpGolfKeypoints(updatedGolfKeypoints),
	}
	return response, nil
}
