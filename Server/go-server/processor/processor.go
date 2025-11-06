package processor

import (
	//"context"
	"log"
	//"sync"

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
	return nil, nil
}

func (p *Processor) GetFaceOnPoseSetupPoints(request *cv.GetFaceOnPoseSetupPointsRequest) (*cv.GetFaceOnPoseSetupPointsResponse, error) {
	getOpenPoseDataResponseCalibration, err := p.ocvmgr.GetOpenPoseData(request.CalibratedImage.CalibrationImage.Bytes)
	if err != nil {
		return nil, err
	}
	calibratedAxes, err := checkIfCalibrationImageIsGood(getOpenPoseDataResponseCalibration.Keypoints)
	if err != nil {
		return nil, err
	}
	getOpenPoseDataResponseImg, err := p.ocvmgr.GetOpenPoseData(request.CalibratedImage.Image.Bytes)
	if err != nil {
		return nil, err
	}
	sideBend := getSideBend(getOpenPoseDataResponseImg.Keypoints, calibratedAxes)
	log.Printf("Side bend is %f", sideBend)
	response := &cv.GetFaceOnPoseSetupPointsResponse{
		SetupPoints: &cv.FaceOnGolfSetupPoints{
			SideBend: sideBend,
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
