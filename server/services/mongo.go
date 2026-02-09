package services

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var db *mongo.Database

func ConnectMongo() error {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return errors.New("MONGO_URI missing")
	}
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "jobtracker"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	// Ping
	if err := c.Ping(ctx, nil); err != nil {
		return err
	}

	client = c
	db = client.Database(dbName)

	// Indexes
	return ensureIndexes(ctx)
}

func DisconnectMongo() {
	if client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = client.Disconnect(ctx)
}

func UsersCollection() *mongo.Collection {
	return db.Collection("users")
}
