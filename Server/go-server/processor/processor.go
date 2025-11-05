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

func (p *Processor) StartOpenCvApiClient() {
	p.ocvmgr.StartOpenCvApiClient()
}

func (p *Processor) CloseOpenCvApiClient() {
	p.ocvmgr.CloseOpenCvApiClient()
}

func (p *Processor) ShowDTLPoseImage(request *cv.ShowDTLPoseImageRequest) *cv.ShowDTLPoseImageResponse {
	getOpenPoseImageResponse := p.ocvmgr.GetOpenPoseImage(request.Image.Bytes)
	return &cv.ShowDTLPoseImageResponse{
		Image: getOpenPoseImageResponse.Image,
	}
}

func (p *Processor) ShowFaceOnPoseImage(request *cv.ShowFaceOnPoseImageRequest) *cv.ShowFaceOnPoseImageResponse {
	getOpenPoseImageResponse := p.ocvmgr.GetOpenPoseImage(request.Image.Bytes)
	return &cv.ShowFaceOnPoseImageResponse{
		Image: getOpenPoseImageResponse.Image,
	}
}

func (p *Processor) ShowDTLPoseImagesFromVideo(requests []*cv.ShowDTLPoseImageRequest) []*cv.ShowDTLPoseImageResponse {
	images := [][]byte{}
	for _, request := range requests {
		img := request.Image.Bytes
		images = append(images, img)
	}
	openPoseResponses := p.ocvmgr.GetOpenPoseImagesFromFromVideo(images)
	responses := []*cv.ShowDTLPoseImageResponse{}
	for _, openPoseResponse := range openPoseResponses {
		response := &cv.ShowDTLPoseImageResponse{
			Image: openPoseResponse.Image,
		}
		responses = append(responses, response)
	}
	return responses
}

func (p *Processor) ShowFaceOnPoseImagesFromVideo(requests []*cv.ShowFaceOnPoseImageRequest) []*cv.ShowFaceOnPoseImageResponse {
	images := [][]byte{}
	for _, request := range requests {
		img := request.Image.Bytes
		images = append(images, img)
	}
	openPoseResponses := p.ocvmgr.GetOpenPoseImagesFromFromVideo(images)
	responses := []*cv.ShowFaceOnPoseImageResponse{}
	for _, openPoseResponse := range openPoseResponses {
		response := &cv.ShowFaceOnPoseImageResponse{
			Image: openPoseResponse.Image,
		}
		responses = append(responses, response)
	}
	return responses
}
