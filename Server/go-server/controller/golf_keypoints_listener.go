package controller

import (
	"context"
	"fmt"

	db "github.com/sirfrank96/go-server/db"
	opencvclient "github.com/sirfrank96/go-server/opencv-client"
	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

type GolfKeypointsListener struct {
	skp.UnimplementedGolfKeypointsServiceServer
	ocvmgr *opencvclient.OpenCvClientManager
	dbmgr  *db.DbManager
}

func newGolfKeypointsListener(ocvmgr *opencvclient.OpenCvClientManager, dbmgr *db.DbManager) *GolfKeypointsListener {
	return &GolfKeypointsListener{
		ocvmgr: ocvmgr,
		dbmgr:  dbmgr,
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
	// TODO: Delete golf keypoint assocaited with input image
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
	// get inputimg with inputimgid from db
	inputImage, err := g.dbmgr.ReadInputImage(ctx, request.InputImageId)
	if err != nil {
		return nil, fmt.Errorf("could not read input image with id: %s: %w", request.InputImageId, err)
	}
	inputImage.CalibrationImgAxes = request.CalibrationImageAxes
	inputImage.CalibrationImgVanishingPoint = request.CalibrationImageVanishingPoint
	calibrationInfo := &util.CalibrationInfo{
		CalibrationType: request.CalibrationType,
		FeetLineMethod:  request.FeetLineMethod,
	}
	// put in golf ball/golf club points and warnings
	if request.GolfBall != nil {
		calibrationInfo.GolfBallPoint = *util.ConvertCvKeypointToPoint(request.GolfBall)
	} else {
		calibrationInfo.GolfBallWarning = util.WarningImpl{
			Severity: util.MINOR,
			Message:  "no golf ball identified, may not be able to provide all setup points",
		}
	}
	if request.ClubButt != nil {
		calibrationInfo.ClubButtPoint = *util.ConvertCvKeypointToPoint(request.ClubButt)
	} else {
		calibrationInfo.GolfClubWarning = util.WarningImpl{
			Severity: util.MINOR,
			Message:  "no club butt identified, may not be able to provide all setup points",
		}
	}
	if request.ClubHead != nil {
		calibrationInfo.ClubHeadPoint = *util.ConvertCvKeypointToPoint(request.ClubHead)
	} else {
		calibrationInfo.GolfClubWarning = util.WarningImpl{
			Severity: util.MINOR,
			Message:  "no club head identified, may not be able to provide all setup points",
		}
	}
	// dtl calibration via calibration images
	if inputImage.ImageType == skp.ImageType_DTL {
		// axes calibration
		if calibrationInfo.CalibrationType != skp.CalibrationType_NO_CALIBRATION {
			getOpenPoseDataResponse, err := g.ocvmgr.GetOpenPoseData(inputImage.CalibrationImgAxes)
			if err != nil {
				return nil, fmt.Errorf("could not get openpose data for calibration image axes %w", err)
			}
			var warning util.Warning
			calibrationInfo, warning = util.VerifyCalibrationImageAxes(getOpenPoseDataResponse.Keypoints, calibrationInfo)
			if warning != nil {
				return nil, fmt.Errorf("could not verify calibration image axes: %s", warning.Error())
			}
			// vanishing point calibration
			if calibrationInfo.CalibrationType != skp.CalibrationType_AXES_CALIBRATION_ONLY {
				getOpenPoseDataResponse, err := g.ocvmgr.GetOpenPoseData(inputImage.CalibrationImgVanishingPoint)
				if err != nil {
					return nil, fmt.Errorf("could not get openpose data for calibration image vanishingpoint %w", err)
				}
				calibrationInfo, warning = util.VerifyCalibrationImageVanishingPoint(getOpenPoseDataResponse.Keypoints, calibrationInfo)
				if warning != nil {
					return nil, fmt.Errorf("could not verify calibration image axes: %s", warning.Error())
				}
				// no vanishing point calibration
			} else {
				calibrationInfo.VanishingPointCalibrationWarning = util.WarningImpl{
					Severity: util.MINOR,
					Message:  "no vanishing point calibration, may not be able to provide all setup points",
				}
			}
			// no axes or vanishing point calibration
		} else {
			calibrationInfo.AxesCalibrationWarning = util.WarningImpl{
				Severity: util.MINOR,
				Message:  "no axes or vanishing point calibration, may not be able to provide all setup points",
			}
		}
		// face on calibration via calibration image
	} else {
		// axes calibration
		if calibrationInfo.CalibrationType != skp.CalibrationType_NO_CALIBRATION {
			getOpenPoseDataResponse, err := g.ocvmgr.GetOpenPoseData(inputImage.CalibrationImgAxes)
			if err != nil {
				return nil, fmt.Errorf("could not get openpose data for calibration image axes %w", err)
			}
			var warning util.Warning
			calibrationInfo, warning = util.VerifyCalibrationImageAxes(getOpenPoseDataResponse.Keypoints, calibrationInfo)
			if warning != nil {
				return nil, fmt.Errorf("could not verify calibration image axes: %s", warning.Error())
			}
			// no axes calibration
		} else {
			calibrationInfo.AxesCalibrationWarning = util.WarningImpl{
				Severity: util.MINOR,
				Message:  "no axes calibration, may not be able to provide all setup points",
			}
		}
	}
	inputImage.CalibrationInfo = *calibrationInfo
	// update inputimg with inputimgid in db
	_, err = g.dbmgr.UpdateInputImage(ctx, request.InputImageId, inputImage)
	if err != nil {
		return nil, fmt.Errorf("could not update input image with id: %s with calibration info: %w", request.InputImageId, err)
	}
	// return response
	response := &skp.CalibrateInputImageResponse{
		Success: true,
	}
	return response, nil
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

	// TODO: Use GETOPENPOSEALL
	// get openpose image for input img
	getOpenPoseImageResponse, err := g.ocvmgr.GetOpenPoseImage(inputImage.InputImg)
	if err != nil {
		return nil, fmt.Errorf("could not get open pose image for image: %w", err)
	}
	// get openpose data for input img
	getOpenPoseDataResponse, err := g.ocvmgr.GetOpenPoseData(inputImage.InputImg)
	if err != nil {
		return nil, fmt.Errorf("could not get open pose data for image: %w", err)
	}

	// calculate golf setup points
	golfKeypoints := &db.GolfKeypoints{
		UserId:          userId,
		InputImageId:    request.InputImageId,
		OutputImg:       getOpenPoseImageResponse.Image,
		OutputKeypoints: *getOpenPoseDataResponse.Keypoints,
	}
	// dtl setup points
	if inputImage.ImageType == skp.ImageType_DTL {
		golfKeypoints.DtlGolfSetupPoints = *CalculateDTLSetupPoints(ctx, getOpenPoseDataResponse.Keypoints, &inputImage.CalibrationInfo)
	} else { // face on setup points
		golfKeypoints.FaceonGolfSetupPoints = *CalculateFaceOnSetupPoints(ctx, getOpenPoseDataResponse.Keypoints, &inputImage.CalibrationInfo)
	}

	// store golfkeypoints in db
	_, err = g.dbmgr.CreateGolfKeypoints(ctx, golfKeypoints)
	if err != nil {
		return nil, fmt.Errorf("could not store golfkeypoints in db %w", err)
	}

	// return response
	response := &skp.CalculateGolfKeypointsResponse{
		Success:       true,
		OutputImage:   getOpenPoseImageResponse.Image,
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
	golfKeypoints.OutputKeypoints = *db.UpdateOutputKeypoints(&golfKeypoints.OutputKeypoints, request.UpdatedBodyKeypoints)
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
