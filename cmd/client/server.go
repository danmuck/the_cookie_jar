package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/danmuck/the_cookie_jar/api/database"
	"github.com/danmuck/the_cookie_jar/api/models"
	"github.com/danmuck/the_cookie_jar/api/routers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/address"
)

type Server struct {
	addr         address.Address
	organization string
	org_key      string
	router       *gin.Engine
	db           *database.Database
	plugins      map[string]*gin.Engine
}

func initServer(org_key string) Server {
	db := database.GetClient()
	s := Server{
		addr:         address.Address("localhost"),
		organization: "test_org",
		org_key:      org_key,

		db:      db,
		router:  routers.BaseRouter(),
		plugins: make(map[string]*gin.Engine),
	}

	return s
}

func (s *Server) db_AddUser(user models.User) error {
	coll := s.db.Client.Database("the_cookie_jar").Collection("users")
	filter := bson.M{"username": user.GetUsername()}

	var result models.User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.db.UpdateUser(user)
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

	test_users := []string{"Dan M.", "Michael R.", "Michael Y.", "Saqib M.", "Cordell H."}
	other_users := []string{"Guest", "Professor", "TA", "Admin", "Student"}
	for {

		i := int64(rand.Intn(len(test_users)))
		rs := test_users[i]
		maker := *models.NewUser(rs)
		fmt.Printf(">> waiting to insert random maker (%v with %v: \"%v\") \n>> .. ctrl-c to quit .. \n",
			maker.GetUsername(), maker.Org, maker.GetStatus_String())

		time.Sleep(10 * time.Second)
		err = server.db_AddUser(maker)
		if err != nil {
			fmt.Printf(">> error: %v \n", err.Error())

			i = int64(rand.Intn(len(other_users)))
			rs = other_users[i]
			other := *models.NewUser(rs)
			server.db_AddUser(other)
		}
	}

}
