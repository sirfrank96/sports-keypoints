package controller

import (
	"context"
	"log"

	cvclient "github.com/sirfrank96/go-server/cv-client"
	db "github.com/sirfrank96/go-server/db"
	kpserver "github.com/sirfrank96/go-server/keypoints-server"
)

type Controller struct {
	cvmgr *cvclient.CvClientManager
	dbmgr *db.DbManager
	kpmgr *kpserver.KeypointsServerManager
}

func NewController() *Controller {
	p := &Controller{}
	p.cvmgr = cvclient.NewCvClientManager()
	p.dbmgr = db.NewDbManager()
	p.kpmgr = kpserver.NewKeypointsServerManager(newGolfKeypointsListener(p.cvmgr, p.dbmgr), newUserListener(p.cvmgr, p.dbmgr))
	log.Printf("New Controller")
	return p
}

func (c *Controller) StartCvClient() error {
	return c.cvmgr.StartCvClient()
}

func (c *Controller) StartDatabaseClient(ctx context.Context) error {
	return c.dbmgr.StartMongoDBClient(ctx)
}

func (c *Controller) StartKeypointsServer() error {
	return c.kpmgr.StartKeypointsServer()
}

func (c *Controller) CloseCvClient() error {
	return c.cvmgr.CloseCvClient()
}

func (c *Controller) CloseDatabaseClient(ctx context.Context) error {
	return c.dbmgr.CloseMongoDBClient(ctx)
}

func (c *Controller) StopKeypointsServer() error {
	return c.kpmgr.StopKeypointsServer()
}
