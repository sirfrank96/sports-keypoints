package controller

/*
import (
	"context"
	"fmt"
	"log"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	opencvclient "github.com/sirfrank96/go-server/opencv-client"
	db "github.com/sirfrank96/go-server/db"
)

// TODO: RENAME PROCESSOR TO CONTROLLER? ORGANIZER? SUPERVISOR? AND HAVE ALL THE STARTSERVICES Functions start from here (except client_api_mgr starts)
//TODO: Init a supervisor and pass it as a param to client_api_mgr?
type Processor struct {
	//mutex  sync.Mutex
	ocvmgr *cvapi.OpenCvClientManager
	dbmgr  *db.DbManager
	//kpmgr *kpserver.KeypointsServerManager
}

func NewProcessor() *Processor {
	p := &Processor{}
	p.ocvmgr = cvapi.NewOpenCvApiManager()
	p.dbmgr = db.NewDbManager()
	log.Printf("New Processor")
	return p
}

func (p *Processor) StartOpenCvClient() error {
	return p.ocvmgr.StartOpenCvClient()
}

func (p *Processor) CloseOpenCvClient() error {
	return p.ocvmgr.CloseOpenCvClient()
}

func (p *Processor) StartDatabaseClient() error {
	return p.dbmgr.StartMongoDBClient(p.ctx)
}

func (p *Processor) StopDatabaseClient() error {
	return p.dbmgr.StopMongoDBClient(p.ctx)
}

func (p *Processor) ReadImageInfo(ctx context.Context, request *cv.ReadImageInfoRequest) (*cv.ImageInfo, error) {
	if err := p.verifyRequest(ctx, request.UserId, request.ImgId, nil); err != nil {
		return nil, fmt.Errorf("unable to verify request: %w", err)
	}
	imageInfo, err := p.dbmgr.ReadImageInfo(ctx, request.ImgId)
	if err != nil {
		return nil, fmt.Errorf("unable to read image info %w", err)
	}
	return imageInfo, nil
}

func (p *Processor) GetDTLPoseImage(ctx context.Context, request *cv.GetDTLPoseImageRequest) (*cv.GetDTLPoseImageResponse, error) {
	var image *cv.Image
	if request.Image != nil {
		image = request.Image
	}
	if err := p.verifyRequest(ctx, request.UserId, request.ImgId, image); err != nil {
		return nil, fmt.Errorf("unable to verify request: %w", err)
	}
	var inputImgId string
	var outputImage *cv.Image
	if request.ImgId != "" { // If img id is provided, fetch from db
		var err error
		image, err := p.verifyInputImageExists(ctx, request.ImgId)
		if err != nil {
			return nil, err
		}
		inputImgId = image.Id
		imageInfo, err := p.dbmgr.ReadImageInfo(ctx, inputImgId)
		if err != nil {
			return nil, err
		}
		// TODO: If output img isnt there, then run opencv api
		if imageInfo.OutputImg == nil {
			return nil, fmt.Errorf("ouput image not yet generated")
		} else {
			outputImage = imageInfo.OutputImg
		}
	} else { // If image is provided, process request
		getOpenPoseImageResponse, err := p.ocvmgr.GetOpenPoseImage(request.Image.Bytes)
		if err != nil {
			return nil, err
		}
		// Store results in db
		imageInfo := &cv.ImageInfo{
			UserId:    request.UserId,
			ImageType: cv.ImageType_IMAGE_TYPE_DTL,
			OutputImg: getOpenPoseImageResponse.Image,
		}
		imageInfo, err = p.storeImageResults(ctx, request.Image, imageInfo)
		if err != nil {
			return nil, fmt.Errorf("could not store results in in db %w", err)
		}
		inputImgId = imageInfo.InputImgId
		outputImage = getOpenPoseImageResponse.Image
	}
	// Return results
	response := &cv.GetDTLPoseImageResponse{
		InputImgId: inputImgId,
		Image:      outputImage,
	}
	return response, nil
}

func (p *Processor) GetFaceOnPoseImage(ctx context.Context, request *cv.GetFaceOnPoseImageRequest) (*cv.GetFaceOnPoseImageResponse, error) {
	var image *cv.Image
	if request.Image != nil {
		image = request.Image
	}
	if err := p.verifyRequest(ctx, request.UserId, request.ImgId, image); err != nil {
		return nil, fmt.Errorf("unable to verify request: %w", err)
	}
	var inputImgId string
	var outputImage *cv.Image
	if request.ImgId != "" { // If img id is provided, fetch from db
		var err error
		image, err := p.verifyInputImageExists(ctx, request.ImgId)
		if err != nil {
			return nil, err
		}
		inputImgId = image.Id
		imageInfo, err := p.dbmgr.ReadImageInfo(ctx, inputImgId)
		if err != nil {
			return nil, err
		}
		// TODO: If output img isnt there, then run opencv api
		outputImage = imageInfo.OutputImg
	} else { // If image is provided, process request
		getOpenPoseImageResponse, err := p.ocvmgr.GetOpenPoseImage(request.Image.Bytes)
		if err != nil {
			return nil, err
		}
		// Store results in db
		imageInfo := &cv.ImageInfo{
			UserId:    request.UserId,
			ImageType: cv.ImageType_IMAGE_TYPE_FACE_ON,
			OutputImg: getOpenPoseImageResponse.Image,
		}
		imageInfo, err = p.storeImageResults(ctx, request.Image, imageInfo)
		if err != nil {
			return nil, fmt.Errorf("could not store results in in db %w", err)
		}
		inputImgId = imageInfo.InputImgId
		outputImage = getOpenPoseImageResponse.Image
	}
	// Return results
	response := &cv.GetFaceOnPoseImageResponse{
		InputImgId: inputImgId,
		Image:      outputImage,
	}
	return response, nil
}

// 3 code paths: just grab image and imageinfo and return, update imageinfo, and create new image info
func (p *Processor) GetDTLPoseSetupPoints(ctx context.Context, request *cv.GetDTLPoseSetupPointsRequest) (*cv.GetDTLPoseSetupPointsResponse, error) {
	var inputImg *cv.Image
	// TODO: Move these to client-api-mgr verify functions
	if request.CalibratedImage != nil && request.CalibratedImage.Image != nil {
		inputImg = request.CalibratedImage.Image
	}
	if err := p.verifyRequest(ctx, request.UserId, request.ImgId, inputImg); err != nil {
		return nil, fmt.Errorf("unable to verify request: %w", err)
	}
	var inputImgId string
	var imageInfo *cv.ImageInfo
	// If img id is provided, fetch from db
	if request.ImgId != "" {
		var err error
		inputImg, err = p.verifyInputImageExists(ctx, request.ImgId)
		if err != nil {
			return nil, err
		}
		inputImgId = inputImg.Id
		imageInfo, err = p.dbmgr.ReadImageInfo(ctx, inputImgId)
		if err != nil {
			return nil, err
		}
		// if dtlsetuppoints is there then return, otherwise do openpose api processing
		if imageInfo.DtlGolfSetupPoints != nil {
			return &cv.GetDTLPoseSetupPointsResponse{
				InputImgId:  inputImgId,
				SetupPoints: imageInfo.DtlGolfSetupPoints,
			}, nil
		}
	}
	getOpenPoseDataResponseCalibrationAxes, err := p.ocvmgr.GetOpenPoseData(request.CalibratedImage.CalibrationImageAxes.Bytes)
	if err != nil {
		return nil, err
	}
	getOpenPoseDataResponseCalibrationVanishingPoint, err := p.ocvmgr.GetOpenPoseData(request.CalibratedImage.CalibrationImageVanishingPoint.Bytes)
	if err != nil {
		return nil, err
	}
	calibrationInfo, warning := VerifyDTLCalibrationImages(getOpenPoseDataResponseCalibrationAxes.Keypoints, getOpenPoseDataResponseCalibrationVanishingPoint.Keypoints, request.CalibratedImage.FeetLineMethod)
	if warning != nil {
		return nil, fmt.Errorf("could not verify dtl calibration images: %w", warning)
	}
	getOpenPoseDataResponse, err := p.ocvmgr.GetOpenPoseData(inputImg.Bytes)
	if err != nil {
		return nil, err
	}
	spineAngle, warning := GetSpineAngle(getOpenPoseDataResponse.Keypoints, calibrationInfo)
	var spineAngleWarning string
	if warning != nil {
		spineAngleWarning = warning.Error()
	}
	log.Printf("Spine angle is %f", spineAngle)
	feetAlignment, warning := GetFeetAlignment(getOpenPoseDataResponse.Keypoints, calibrationInfo)
	var feetAlignmentWarning string
	if warning != nil {
		feetAlignmentWarning = warning.Error()
	}
	log.Printf("Feet alignment is %f", feetAlignment)
	dtlSetupPoints := &cv.DTLGolfSetupPoints{
		SpineAngle: &cv.Double{
			Data:    spineAngle,
			Warning: spineAngleWarning,
		},
		FeetAlignment: &cv.Double{
			Data:    feetAlignment,
			Warning: feetAlignmentWarning,
		},
	}
	if imageInfo != nil {
		imageInfo.CalibrationImgAxes =
		imageInfo.OutputKeypoints = getOpenPoseDataResponse.Keypoints
		imageInfo.DtlGolfSetupPoints = dtlSetupPoints
		imageInfo, err = p.dbmgr.UpdateImageInfo()
	} else {
		// Store results in db
		imageInfo = &cv.ImageInfo{
			UserId:                       request.UserId,
			ImageType:                    cv.ImageType_IMAGE_TYPE_DTL,
			CalibrationImgAxes:           request.CalibratedImage.CalibrationImageAxes,
			CalibrationImgVanishingPoint: request.CalibratedImage.CalibrationImageVanishingPoint,
			OutputKeypoints:              getOpenPoseDataResponse.Keypoints,
			DtlGolfSetupPoints:           dtlSetupPoints,
		}
		imageInfo, err = p.storeImageResults(ctx, request.CalibratedImage.Image, imageInfo)
		if err != nil {
			return nil, fmt.Errorf("could not store results in in db %w", err)
		}
		inputImgId = imageInfo.InputImgId
	}

	response := &cv.GetDTLPoseSetupPointsResponse{
		InputImgId:  inputImgId,
		SetupPoints: dtlSetupPoints,
	}
	return response, nil
}

func (p *Processor) GetFaceOnPoseSetupPoints(ctx context.Context, request *cv.GetFaceOnPoseSetupPointsRequest) (*cv.GetFaceOnPoseSetupPointsResponse, error) {
	var image *cv.Image
	if request.CalibratedImage != nil && request.CalibratedImage.Image != nil {
		image = request.CalibratedImage.Image
	}
	if err := p.verifyRequest(ctx, request.UserId, request.ImgId, image); err != nil {
		return nil, fmt.Errorf("unable to verify request: %w", err)
	}
	var inputImgId string
	var faceOnSetupPoints *cv.FaceOnGolfSetupPoints
	if request.ImgId != "" { // If img id is provided, fetch from db
		var err error
		image, err := p.verifyInputImageExists(ctx, request.ImgId)
		if err != nil {
			return nil, err
		}
		inputImgId = image.Id
		imageInfo, err := p.dbmgr.ReadImageInfo(ctx, inputImgId)
		if err != nil {
			return nil, err
		}
		// TODO: If setup points isnt there, then run opencv api
		faceOnSetupPoints = imageInfo.FaceOnGolfSetupPoints
	} else { // If image is provided, process request
		getOpenPoseDataResponseCalibration, err := p.ocvmgr.GetOpenPoseData(request.CalibratedImage.CalibrationImageAxes.Bytes)
		if err != nil {
			return nil, err
		}
		calibrationInfo, warning := VerifyFaceOnCalibrationImage(getOpenPoseDataResponseCalibration.Keypoints, request.CalibratedImage.FeetLineMethod)
		if warning != nil {
			return nil, fmt.Errorf("could not verify face on calibration images: %w", warning)
		}
		getOpenPoseDataResponseImg, err := p.ocvmgr.GetOpenPoseData(request.CalibratedImage.Image.Bytes)
		if err != nil {
			return nil, err
		}
		sideBend, warning := GetSideBend(getOpenPoseDataResponseImg.Keypoints, calibrationInfo)
		var sideBendWarning string
		if warning != nil {
			sideBendWarning = warning.Error()
		}
		log.Printf("Side bend is %f", sideBend)
		lFootFlare, warning := GetLeftFootFlare(getOpenPoseDataResponseImg.Keypoints, calibrationInfo)
		var lFootFlareWarning string
		if warning != nil {
			lFootFlareWarning = warning.Error()
		}
		log.Printf("Left foot flare is %f", lFootFlare)
		rFootFlare, warning := GetRightFootFlare(getOpenPoseDataResponseImg.Keypoints, calibrationInfo)
		var rFootFlareWarning string
		if warning != nil {
			rFootFlareWarning = warning.Error()
		}
		log.Printf("Right foot flare is %f", rFootFlare)
		faceOnSetupPoints = &cv.FaceOnGolfSetupPoints{
			SideBend: &cv.Double{
				Data:    sideBend,
				Warning: sideBendWarning,
			},
			LFootFlare: &cv.Double{
				Data:    lFootFlare,
				Warning: lFootFlareWarning,
			},
			RFootFlare: &cv.Double{
				Data:    rFootFlare,
				Warning: rFootFlareWarning,
			},
		}
		// Store results in db
		imageInfo := &cv.ImageInfo{
			UserId:                request.UserId,
			ImageType:             cv.ImageType_IMAGE_TYPE_FACE_ON,
			CalibrationImgAxes:    request.CalibratedImage.CalibrationImageAxes,
			OutputKeypoints:       getOpenPoseDataResponseImg.Keypoints,
			FaceOnGolfSetupPoints: faceOnSetupPoints,
		}
		imageInfo, err = p.storeImageResults(ctx, request.CalibratedImage.Image, imageInfo)
		if err != nil {
			return nil, fmt.Errorf("could not store results in in db %w", err)
		}
		inputImgId = imageInfo.InputImgId
	}
	response := &cv.GetFaceOnPoseSetupPointsResponse{
		InputImgId:  inputImgId,
		SetupPoints: faceOnSetupPoints,
	}
	return response, nil
}

func (p *Processor) GetDTLPoseImagesFromVideo(requests []*cv.GetDTLPoseImageRequest) ([]*cv.GetDTLPoseImageResponse, error) {
	images := [][]byte{}
	for _, request := range requests {
		img := request.Image.Bytes
		images = append(images, img)
	}
	openPoseResponses, err := p.ocvmgr.GetOpenPoseImagesFromFromVideo(images)
	if err != nil {
		return nil, err
	}
	responses := []*cv.GetDTLPoseImageResponse{}
	for _, openPoseResponse := range openPoseResponses {
		response := &cv.GetDTLPoseImageResponse{
			Image: openPoseResponse.Image,
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (p *Processor) GetFaceOnPoseImagesFromVideo(requests []*cv.GetFaceOnPoseImageRequest) ([]*cv.GetFaceOnPoseImageResponse, error) {
	images := [][]byte{}
	for _, request := range requests {
		img := request.Image.Bytes
		images = append(images, img)
	}
	openPoseResponses, err := p.ocvmgr.GetOpenPoseImagesFromFromVideo(images)
	if err != nil {
		return nil, err
	}
	responses := []*cv.GetFaceOnPoseImageResponse{}
	for _, openPoseResponse := range openPoseResponses {
		response := &cv.GetFaceOnPoseImageResponse{
			Image: openPoseResponse.Image,
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (p *Processor) CreateUser(ctx context.Context, request *cv.CreateUserRequest) (*cv.User, error) {
	user, err := p.dbmgr.CreateUser(ctx, &cv.User{
		UserName: request.UserName,
		Email:    request.Email,
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *Processor) ReadUser(ctx context.Context, request *cv.ReadUserRequest) (*cv.User, error) {
	return p.dbmgr.ReadUser(ctx, request.Id)
}

func (p *Processor) UpdateUser(ctx context.Context, request *cv.UpdateUserRequest) (*cv.User, error) {
	updatedUser, err := p.dbmgr.UpdateUser(ctx, request.Id, &cv.User{
		UserName: request.UserName,
		Email:    request.Email,
	})
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (p *Processor) DeleteUser(ctx context.Context, request *cv.DeleteUserRequest) (*cv.DeleteUserResponse, error) {
	if err := p.dbmgr.DeleteUser(ctx, request.Id); err != nil {
		return nil, err
	}
	response := &cv.DeleteUserResponse{
		Success: true,
	}
	return response, nil
}

func (p *Processor) verifyRequest(ctx context.Context, userId string, imgId string, image *cv.Image) error {
	// Make sure user exists
	if _, err := p.verifyUserExists(ctx, userId); err != nil {
		return err
	}
	// Make sure at least one of img id or image are provided
	if imgId == "" && image == nil {
		return fmt.Errorf("please provide either an image id or an actual image in request")
	}
	return nil
}

func (p *Processor) verifyUserExists(ctx context.Context, userId string) (*cv.User, error) {
	user, err := p.dbmgr.ReadUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("could not find user %s: %w", userId, err)
	}
	return user, nil
}

func (p *Processor) verifyInputImageExists(ctx context.Context, imgId string) (*cv.Image, error) {
	image, err := p.dbmgr.ReadInputImage(ctx, imgId)
	if err != nil {
		return nil, fmt.Errorf("could not find input image %s: %w", imgId, err)
	}
	return image, nil
}

func (p *Processor) storeImageResults(ctx context.Context, inputImage *cv.Image, imageInfo *cv.ImageInfo) (*cv.ImageInfo, error) {
	inputImage, err := p.dbmgr.CreateInputImage(ctx, inputImage)
	if err != nil {
		return nil, fmt.Errorf("could not store input image: %w", err)
	}
	imageInfo.InputImgId = inputImage.Id
	imageInfo, err = p.dbmgr.CreateImageInfo(ctx, imageInfo)
	if err != nil {
		return nil, fmt.Errorf("could not store image info: %w", err)
	}
	return imageInfo, nil
}

/*func (p *Processor) IdentifyDTLPoseDifferences(request *cv.IdentifyDTLPoseDifferencesRequest) (*cv.IdentifyDTLPoseDifferencesResponse, error) {
	return nil, nil
}

func (p *Processor) IdentifyFaceOnPoseDifferences(request *cv.IdentifyFaceOnPoseDifferencesRequest) (*cv.IdentifyFaceOnPoseDifferencesResponse, error) {
	return nil, nil
}*/
