package sandbox

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	uri    string
	Client *mongo.Client
}

func NewDatabase() (*Database, error) {
	db := &Database{
		uri:    "mongodb://0.0.0.0:27017/",
		Client: nil,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	opts := options.Client().ApplyURI(db.uri)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping to database: %v", err)
	}

	db.Client = client
	fmt.Println("Database Created")
	return db, nil
}
