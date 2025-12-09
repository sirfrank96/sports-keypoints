package db

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	mongodb "go.mongodb.org/mongo-driver/mongo"
	mongoopts "go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbaddr = flag.String("dbaddr", "mongodb://localhost:27017", "the address to connect to")
)

type DbManager struct {
	mutex                  sync.Mutex
	clientOptions          *mongoopts.ClientOptions
	client                 *mongodb.Client
	db                     *mongodb.Database
	userCollection         *mongodb.Collection
	inputImageCollection   *mongodb.Collection
	golfKeypointCollection *mongodb.Collection
}

func NewDbManager() *DbManager {
	d := &DbManager{}
	log.Printf("New Database Mgr")
	return d
}

func (d *DbManager) StartMongoDBClient(ctx context.Context) error {
	log.Printf("Starting MongoDB Client")
	flag.Parse()
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = *dbaddr
	}
	// Set client options
	d.clientOptions = mongoopts.Client().ApplyURI(mongoURI)
	// Connect to MongoDB
	var err error
	d.client, err = mongodb.Connect(ctx, d.clientOptions)
	if err != nil {
		return fmt.Errorf("could not connect to mongodb %w", err)
	}
	// Check the connection
	err = d.client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not ping mongodb %w", err)
	}
	// Create Database
	d.db = d.client.Database("golfkeypointsdatabase")
	// Create Collections
	d.userCollection = d.db.Collection("users")
	d.inputImageCollection = d.db.Collection("inputimages")
	d.golfKeypointCollection = d.db.Collection("golfkeypoints")
	return nil
}

func (d *DbManager) CloseMongoDBClient(ctx context.Context) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.client.Disconnect(ctx)
}
