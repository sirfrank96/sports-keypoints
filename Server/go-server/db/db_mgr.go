package db

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	mongodb "go.mongodb.org/mongo-driver/mongo"
	mongoopts "go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbaddr = flag.String("dbaddr", "mongodb://localhost:27017", "the address to connect to")
)

type DbManager struct {
	mutex               sync.Mutex
	clientOptions       *mongoopts.ClientOptions
	client              *mongodb.Client
	db                  *mongodb.Database
	userCollection      *mongodb.Collection
	imageCollection     *mongodb.Collection
	imageInfoCollection *mongodb.Collection
}

func NewDbManager() *DbManager {
	d := &DbManager{}
	log.Printf("New Database Mgr")
	return d
}

func (d *DbManager) StartMongoDBClient(ctx context.Context) error {
	log.Printf("Starting MongoDB Client")
	flag.Parse()

	// Set client options
	d.clientOptions = mongoopts.Client().ApplyURI(*dbaddr)

	// Connect to MongoDB
	var err error
	d.client, err = mongodb.Connect(ctx, d.clientOptions)
	if err != nil {
		return fmt.Errorf("could not connect to mongodb %w", err)
	}

	// Check the connection
	err = d.client.Ping(ctx, nil)
	if err != nil {
		log.Printf("1.5\n")
		return fmt.Errorf("could not ping mongodb %w", err)
	}

	// Create Database
	d.db = d.client.Database("computervisiongolfdatabase")

	// Create Collections
	d.userCollection = d.db.Collection("users")
	d.imageCollection = d.db.Collection("images")
	d.imageInfoCollection = d.db.Collection("imageinfos")
	return nil
}

func (d *DbManager) StopMongoDBClient(ctx context.Context) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.client.Disconnect(ctx)
}
