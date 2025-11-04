package processor

import (
	"context"
	"log"
	"sync"

	//cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	cvapi "github.com/sirfrank96/go-server/cv-api"
)

type Processor struct {
	ctx    context.Context
	mutex  sync.Mutex
	ocvmgr *cvapi.OpenCvApiManager
}

func NewProcessor() *Processor {
	p := &Processor{}
	p.ocvmgr = cvapi.NewOpenCvApiManager()
	log.Printf("New Processor")
	return p
}

func (p *Processor) StartOpenCVApiClient() {
	p.ocvmgr.StartOpenCVApiClient()
}

func (p *Processor) ShowDTLPoseImage(img []byte) []byte {
	return p.ocvmgr.GetOpenPoseImage(img)
}
