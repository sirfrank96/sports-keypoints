package controller

import (
	"context"
	"log"

	db "github.com/sirfrank96/go-server/db"
	kpmgr "github.com/sirfrank96/go-server/keypoints-server"
	opencvclient "github.com/sirfrank96/go-server/opencv-client"
)

type Controller struct {
	ocvmgr *opencvclient.OpenCvClientManager
	dbmgr  *db.DbManager
	kpmgr  *kpmgr.KeypointsServerManager
}

func NewController() *Controller {
	p := &Controller{}
	p.ocvmgr = opencvclient.NewOpenCvClientManager()
	p.dbmgr = db.NewDbManager()
	p.kpmgr = kpmgr.NewKeypointsServerManager(newGolfKeypointsListener(p.ocvmgr, p.dbmgr), newUserListener(p.ocvmgr, p.dbmgr))
	log.Printf("New Controller")
	return p
}

func (c *Controller) StartOpenCvClient() error {
	return c.ocvmgr.StartOpenCvClient()
}

func (c *Controller) StartDatabaseClient(ctx context.Context) error {
	return c.dbmgr.StartMongoDBClient(ctx)
}

func (c *Controller) StartKeypointsServer() error {
	return c.kpmgr.StartKeypointsServer()
}

func (c *Controller) CloseOpenCvClient() error {
	return c.ocvmgr.CloseOpenCvClient()
}

func (c *Controller) CloseDatabaseClient(ctx context.Context) error {
	return c.dbmgr.CloseMongoDBClient(ctx)
}

func (c *Controller) StopKeypointsServer() error {
	return c.kpmgr.StopKeypointsServer()
}
