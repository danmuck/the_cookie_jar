package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Grabs desired collection from MongoDB database.
*/
func GetCollection(name string) *mongo.Collection {
	client, db := connectToDatabase()
	return client.Database(db).Collection(name)
}

/*
Sets up a connection to the database, performs a test ping, and then returns
MongoDB client and database name (as string).
*/
func connectToDatabase() (*mongo.Client, string) {
	// Loading environment file
	godotenv.Load(".env")

	// Setting Stable API version to 1 for server
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Creating a new client and connecting to server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %v", err))
	}

	// Send a ping to confirm a successful connection
	if err = client.Ping(context.TODO(), nil); err != nil {
		panic(fmt.Errorf("failed to ping database: %v", err))
	}

	return client, os.Getenv("DB_NAME")
}
