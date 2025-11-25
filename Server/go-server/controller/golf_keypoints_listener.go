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
	// dtl calibration
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
		// face on calibration
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
		spineAngle, warning := GetSpineAngle(getOpenPoseDataResponse.Keypoints, &inputImage.CalibrationInfo)
		var spineAngleWarning string
		if warning != nil {
			spineAngleWarning = warning.Error()
		}
		fmt.Printf("Spine angle is %f", spineAngle)
		feetAlignment, warning := GetFeetAlignment(getOpenPoseDataResponse.Keypoints, &inputImage.CalibrationInfo)
		var feetAlignmentWarning string
		if warning != nil {
			feetAlignmentWarning = warning.Error()
		}
		fmt.Printf("Feet alignment is %f", feetAlignment)
		// TODO: Add heel and toe alignment based on FeetLineMethod
		kneeBend, warning := GetKneeBend(getOpenPoseDataResponse.Keypoints)
		var kneeBendWarning string
		if warning != nil {
			kneeBendWarning = warning.Error()
		}
		fmt.Printf("Knee bend is %f", kneeBend)
		shoulderAlignment, warning := GetShoulderAlignment(getOpenPoseDataResponse.Keypoints, &inputImage.CalibrationInfo)
		var shoulderAlignmentWarning string
		if warning != nil {
			shoulderAlignmentWarning = warning.Error()
		}
		fmt.Printf("Shoulder alignment is %f", shoulderAlignment)

		dtlGolfSetupPoints := &skp.DTLGolfSetupPoints{
			SpineAngle: &skp.Double{
				Data:    spineAngle,
				Warning: spineAngleWarning,
			},
			FeetAlignment: &skp.Double{
				Data:    feetAlignment,
				Warning: feetAlignmentWarning,
			},
			KneeBend: &skp.Double{
				Data:    kneeBend,
				Warning: kneeBendWarning,
			},
			ShoulderAlignment: &skp.Double{
				Data:    shoulderAlignment,
				Warning: shoulderAlignmentWarning,
			},
		}
		golfKeypoints.DtlGolfSetupPoints = *dtlGolfSetupPoints
		// face on setup points
	} else {
		sideBend, warning := GetSideBend(getOpenPoseDataResponse.Keypoints, &inputImage.CalibrationInfo)
		var sideBendWarning string
		if warning != nil {
			sideBendWarning = warning.Error()
		}
		fmt.Printf("Side bend is %f", sideBend)
		lFootFlare, warning := GetLeftFootFlare(getOpenPoseDataResponse.Keypoints, &inputImage.CalibrationInfo)
		var lFootFlareWarning string
		if warning != nil {
			lFootFlareWarning = warning.Error()
		}
		fmt.Printf("Left foot flare is %f", lFootFlare)
		rFootFlare, warning := GetRightFootFlare(getOpenPoseDataResponse.Keypoints, &inputImage.CalibrationInfo)
		var rFootFlareWarning string
		if warning != nil {
			rFootFlareWarning = warning.Error()
		}
		fmt.Printf("Right foot flare is %f", rFootFlare)
		stanceWidth, warning := GetStanceWidth(getOpenPoseDataResponse.Keypoints)
		var stanceWidthWarning string
		if warning != nil {
			stanceWidthWarning = warning.Error()
		}
		fmt.Printf("Stance width is %f", stanceWidth)
		shoulderTilt, warning := GetShoulderTilt(getOpenPoseDataResponse.Keypoints, &inputImage.CalibrationInfo)
		var shoulderTiltWarning string
		if warning != nil {
			shoulderTiltWarning = warning.Error()
		}
		fmt.Printf("Shoulder tilt is %f", shoulderTilt)
		faceOnGolfSetupPoints := &skp.FaceOnGolfSetupPoints{
			SideBend: &skp.Double{
				Data:    sideBend,
				Warning: sideBendWarning,
			},
			LFootFlare: &skp.Double{
				Data:    lFootFlare,
				Warning: lFootFlareWarning,
			},
			RFootFlare: &skp.Double{
				Data:    rFootFlare,
				Warning: rFootFlareWarning,
			},
			StanceWidth: &skp.Double{
				Data:    stanceWidth,
				Warning: stanceWidthWarning,
			},
			ShoulderTilt: &skp.Double{
				Data:    shoulderTilt,
				Warning: shoulderTiltWarning,
			},
		}
		golfKeypoints.FaceonGolfSetupPoints = *faceOnGolfSetupPoints
	}

	// store golfkeypoints in db
	_, err = g.dbmgr.CreateGolfKeypoints(ctx, golfKeypoints)
	if err != nil {
		return nil, fmt.Errorf("could not store golfkeypoints in db %w", err)
	}

	// return response
	response := &skp.CalculateGolfKeypointsResponse{
		Success:       true,
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
		Success: true,
		GolfKeypoints: &skp.GolfKeypoints{
			OutputImg:             golfKeypoints.OutputImg,
			DtlGolfSetupPoints:    &golfKeypoints.DtlGolfSetupPoints,
			FaceonGolfSetupPoints: &golfKeypoints.FaceonGolfSetupPoints,
		},
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
