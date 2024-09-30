package sandbox

import (
	"context"
	"fmt"
	// "time"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	uri    string
	Client *mongo.Client
}

// func (db *Database) test_insert(msg string) {

// }

func NewDatabase() (*Database, error) {
	db := &Database{
		uri:    "mongodb://database:27017/",
		Client: nil,
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	// defer cancel()

	// serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// opts := options.Client().ApplyURI(db.uri).SetServerAPIOptions(serverAPI)

	// client, err := mongo.Connect(ctx, opts)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to database: %v", err)
	// }

	// err = client.Ping(ctx, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to ping to database: %v", err)
	// }

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(db.uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	// var result bson.M
	// if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
	// 	panic(err)
	// }
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping to database: %v", err)
	}
	fmt.Println("\n\nPinged your deployment. You successfully connected to MongoDB!\n\n ")

	db.Client = client
	fmt.Println("Database Created")
	return db, nil
}
