package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/danmuck/the_cookie_jar/api"
	"github.com/danmuck/the_cookie_jar/api/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/address"
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

func (db *Database) InsertUser(user models.User) {
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
			db.InsertUser(result)
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

	client.timeout = 30 * time.Second
	err := godotenv.Load(".env")
	uri := os.Getenv("MONGODB_URI")
	name := os.Getenv("DB_NAME")
	if err != nil {
		client.err = fmt.Errorf("cannot find file [.env]::[MONGODB_URI, DB_NAME] %v", err)
	}
	client.uri = uri
	client.name = name

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(client.uri).SetServerAPIOptions(serverAPI)
	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()

	cl, err := mongo.Connect(ctx, opts)
	if err != nil {
		client.err = fmt.Errorf("failed to connect to database: %v", err)
	}

	// Send a ping to confirm a successful connection
	err = cl.Ping(context.TODO(), nil)
	if err != nil {
		client.err = fmt.Errorf("failed to ping to database: %v", err)
	}

	fmt.Printf("\n\n\n uri: %v ", uri)
	fmt.Printf("\nPinged your deployment. You successfully connected to MongoDB! %v\n\n ", uri)

	client.Client = cl
	fmt.Printf("%v Created", client.name)
}

type Server struct {
	addr         address.Address
	organization string
	org_key      string
	router       *gin.Engine
	db           *Database
	plugins      map[string]*gin.Engine
}

func initServer(org_key string) Server {
	db := GetClient()
	s := Server{
		addr:         address.Address("localhost"),
		organization: "test_org",
		org_key:      org_key,

		db:      db,
		router:  api.BaseRouter(),
		plugins: make(map[string]*gin.Engine),
	}

	return s
}

func (s *Server) db_AddUser(user models.User) error {
	coll := client.Client.Database("the_cookie_jar").Collection("users")
	filter := bson.M{"username": user.GetUsername()}

	var result models.User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.db.InsertUser(user)
			return nil
		}
		return err
	}

	return fmt.Errorf("username taken: %v", result.Username)

}

func (s *Server) Serve() error {
	// Listen and Server @ localhost:8080
	// NOTE: this is mapped on Dockerfile and is served to browser at 8080
	err := s.router.Run(":6669")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("\nServer connecting .. \n ")

	server := initServer("test_org_key")
	err := server.router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, ">> Server Initialized .. \n>> addr: %v \n>> org: %v \n>> key: %v", server.addr, server.organization, server.org_key)

	server.router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	go server.Serve()

	test_users := []string{"dan_m", "michael_r", "michael_ys", "saqib_m", "cordell_h"}
	other_users := []string{"Guest", "Professor", "TA", "Admin", "Student"}
	for {

		i := int64(rand.Intn(len(test_users)))
		rs := test_users[i]
		maker := *models.NewUser(rs)
		// fmt.Printf(">> waiting to insert random maker (%v with %v: \"%v\") \n>> .. ctrl-c to quit .. \n",
		// 	maker.GetUsername(), maker.Org, maker.GetStatus_String())

		time.Sleep(10 * time.Second)
		err = server.db_AddUser(maker)
		if err != nil {
			// fmt.Printf(">> error: %v \n", err.Error())

			i = int64(rand.Intn(len(other_users)))
			rs = other_users[i]
			other := *models.NewUser(rs)
			server.db_AddUser(other)
		}
	}

}
