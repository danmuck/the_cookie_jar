package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectMongoDB() *mongo.Client {

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

	fmt.Printf("\n\n\n uri: %v ", uri)
	fmt.Printf("\nPinged your deployment. You successfully connected to MongoDB! %v\n\n ", uri)

	fmt.Printf("%v Created", name)

	return client
}

func PingPong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func Root(c *gin.Context) {
	username := c.Param("username")
	o := fmt.Sprintf("User: %v", username)
	c.String(http.StatusOK, o)
}
