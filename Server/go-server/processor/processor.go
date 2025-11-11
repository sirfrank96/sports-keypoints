package processor

import (
	"fmt"
	"log"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	cvapi "github.com/sirfrank96/go-server/cv-api"
)

type Processor struct {
	//ctx    context.Context
	//mutex  sync.Mutex
	ocvmgr *cvapi.OpenCvApiManager
}

func NewProcessor() *Processor {
	p := &Processor{}
	p.ocvmgr = cvapi.NewOpenCvApiManager()
	log.Printf("New Processor")
	return p
}

func (p *Processor) StartOpenCvApiClient() error {
	return p.ocvmgr.StartOpenCvApiClient()
}

func (p *Processor) CloseOpenCvApiClient() error {
	return p.ocvmgr.CloseOpenCvApiClient()
}

func (p *Processor) ShowDTLPoseImage(request *cv.ShowDTLPoseImageRequest) (*cv.ShowDTLPoseImageResponse, error) {
	getOpenPoseImageResponse, err := p.ocvmgr.GetOpenPoseImage(request.Image.Bytes)
	if err != nil {
		return nil, err
	}
	response := &cv.ShowDTLPoseImageResponse{
		Image: getOpenPoseImageResponse.Image,
	}
	return response, nil
}

func (p *Processor) ShowFaceOnPoseImage(request *cv.ShowFaceOnPoseImageRequest) (*cv.ShowFaceOnPoseImageResponse, error) {
	getOpenPoseImageResponse, err := p.ocvmgr.GetOpenPoseImage(request.Image.Bytes)
	if err != nil {
		return nil, err
	}
	response := &cv.ShowFaceOnPoseImageResponse{
		Image: getOpenPoseImageResponse.Image,
	}
	return response, nil
}

func (p *Processor) GetDTLPoseSetupPoints(request *cv.GetDTLPoseSetupPointsRequest) (*cv.GetDTLPoseSetupPointsResponse, error) {
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
	getOpenPoseDataResponseImg, err := p.ocvmgr.GetOpenPoseData(request.CalibratedImage.Image.Bytes)
	if err != nil {
		return nil, err
	}
	spineAngle, warning := GetSpineAngle(getOpenPoseDataResponseImg.Keypoints, calibrationInfo)
	var spineAngleWarning string
	if warning != nil {
		spineAngleWarning = warning.Error()
	}
	log.Printf("Spine angle is %f", spineAngle)
	feetAlignment, warning := GetFeetAlignment(getOpenPoseDataResponseImg.Keypoints, calibrationInfo)
	var feetAlignmentWarning string
	if warning != nil {
		feetAlignmentWarning = warning.Error()
	}
	log.Printf("Feet alignment is %f", feetAlignment)
	response := &cv.GetDTLPoseSetupPointsResponse{
		SetupPoints: &cv.DTLGolfSetupPoints{
			SpineAngle: &cv.Double{
				Data:    spineAngle,
				Warning: spineAngleWarning,
			},
			FeetAlignment: &cv.Double{
				Data:    feetAlignment,
				Warning: feetAlignmentWarning,
			},
		},
	}
	return response, nil
}

func (p *Processor) GetFaceOnPoseSetupPoints(request *cv.GetFaceOnPoseSetupPointsRequest) (*cv.GetFaceOnPoseSetupPointsResponse, error) {
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
	response := &cv.GetFaceOnPoseSetupPointsResponse{
		SetupPoints: &cv.FaceOnGolfSetupPoints{
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
		},
	}
	return response, nil
}

func (p *Processor) IdentifyDTLPoseDifferences(request *cv.IdentifyDTLPoseDifferencesRequest) (*cv.IdentifyDTLPoseDifferencesResponse, error) {
	return nil, nil
}

func (p *Processor) IdentifyFaceOnPoseDifferences(request *cv.IdentifyFaceOnPoseDifferencesRequest) (*cv.IdentifyFaceOnPoseDifferencesResponse, error) {
	return nil, nil
}

func (p *Processor) ShowDTLPoseImagesFromVideo(requests []*cv.ShowDTLPoseImageRequest) ([]*cv.ShowDTLPoseImageResponse, error) {
	images := [][]byte{}
	for _, request := range requests {
		img := request.Image.Bytes
		images = append(images, img)
	}
	openPoseResponses, err := p.ocvmgr.GetOpenPoseImagesFromFromVideo(images)
	if err != nil {
		return nil, err
	}
	responses := []*cv.ShowDTLPoseImageResponse{}
	for _, openPoseResponse := range openPoseResponses {
		response := &cv.ShowDTLPoseImageResponse{
			Image: openPoseResponse.Image,
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (p *Processor) ShowFaceOnPoseImagesFromVideo(requests []*cv.ShowFaceOnPoseImageRequest) ([]*cv.ShowFaceOnPoseImageResponse, error) {
	images := [][]byte{}
	for _, request := range requests {
		img := request.Image.Bytes
		images = append(images, img)
	}
	openPoseResponses, err := p.ocvmgr.GetOpenPoseImagesFromFromVideo(images)
	if err != nil {
		return nil, err
	}
	responses := []*cv.ShowFaceOnPoseImageResponse{}
	for _, openPoseResponse := range openPoseResponses {
		response := &cv.ShowFaceOnPoseImageResponse{
			Image: openPoseResponse.Image,
		}
		responses = append(responses, response)
	}
	return responses, nil
}
