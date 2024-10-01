package sandbox

import (
	"context"
	"fmt"
	"os"

	// "time"

	// "github.com/danmuck/the_cookie_jar/sandbox"
	"github.com/danmuck/the_cookie_jar/sandbox/db_types"
	"github.com/joho/godotenv"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	uri    string
	Client *mongo.Client
}

func (db *Database) UpdateUser(user db_types.User) {

	db_ := db.Client.Database("the_cookie_jar")
	db_.CreateCollection(context.TODO(), "users")
	users := db_.Collection("users")

	_, err := users.InsertOne(context.TODO(), user)
	if err != nil {
		fmt.Printf("insert error: %v", err)
	}
}

func NewDatabase() (*Database, error) {
	err := godotenv.Load(".env")
	uri := os.Getenv("MONGODB_URI")

	if err != nil {
		return nil, fmt.Errorf("cannot find file [.env]::[MONGODB_URI] %v", err)
	}
	db := &Database{
		uri:    uri,
		Client: nil,
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	// defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(db.uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Send a ping to confirm a successful connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping to database: %v", err)
	}

	fmt.Printf("\n\n\n uri: %v ", uri)
	fmt.Printf("\nPinged your deployment. You successfully connected to MongoDB! %v\n\n ", uri)

	db.Client = client
	fmt.Println("Database Created")
	return db, nil
}
