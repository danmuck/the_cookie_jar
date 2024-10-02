package database

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/danmuck/the_cookie_jar/api/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client Database
)

type Database struct {
	uri     string
	name    string
	Client  *mongo.Client
	err     error
	start   sync.Once
	timeout time.Duration
}

func (db *Database) UpdateUser(user models.User) {
	db_ := db.Client.Database(db.name)
	db_.CreateCollection(context.TODO(), "users")
	users := db_.Collection("users")

	_, err := users.InsertOne(context.TODO(), user)
	if err != nil {
		fmt.Printf("insert error: %v", err)
	}
}
func (db *Database) LookupUser(username string) *models.User {
	coll := db.Client.Database(db.name).Collection("users")
	filter := bson.M{"username": username}

	var result models.User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			db.UpdateUser(result)
			return &result
		}
		return nil
	}
	return nil
}

func GetClient() *Database {
	client.start.Do(initMongoDB)
	return &client
}

func initMongoDB() {
	db := &Database{
		Client:  nil,
		err:     nil,
		timeout: 30 * time.Second,
	}

	err := godotenv.Load(".env")
	uri := os.Getenv("MONGODB_URI")
	name := os.Getenv("DB_NAME")
	if err != nil {
		db.err = fmt.Errorf("cannot find file [.env]::[MONGODB_URI, DB_NAME] %v", err)
	}
	db.uri = uri
	db.name = name

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(db.uri).SetServerAPIOptions(serverAPI)
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		db.err = fmt.Errorf("failed to connect to database: %v", err)
	}

	// Send a ping to confirm a successful connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		db.err = fmt.Errorf("failed to ping to database: %v", err)
	}

	fmt.Printf("\n\n\n uri: %v ", uri)
	fmt.Printf("\nPinged your deployment. You successfully connected to MongoDB! %v\n\n ", uri)

	db.Client = client
	fmt.Printf("%v Created", db.name)
}
