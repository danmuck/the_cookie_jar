package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func get_collection(coll string) *mongo.Collection {
	client, db := connectMongoDB()
	collection := client.Database(db).Collection(coll)
	return collection
}

func connectMongoDB() (*mongo.Client, string) {

	err := godotenv.Load(".env")
	uri := os.Getenv("MONGODB_URI")
	name := os.Getenv("DB_NAME")
	timeout := 20 * time.Second
	if err != nil {
		panic(fmt.Errorf("cannot find file [.env]::[MONGODB_URI, DB_NAME] %v", err))
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %v", err))
	}

	// Send a ping to confirm a successful connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(fmt.Errorf("failed to ping to database: %v", err))
	}

	fmt.Printf(">> [db] Pinged your deployment. You successfully connected to MongoDB! %v \n", uri)

	return client, name
}
