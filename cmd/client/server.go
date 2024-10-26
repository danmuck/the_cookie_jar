package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/danmuck/the_cookie_jar/pkg/api"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/address"
)

type Server struct {
	addr         address.Address
	organization string
	org_key      string
	router       *gin.Engine
	plugins      map[string]*gin.Engine
}

func initServer(org_key string) Server {
	s := Server{
		addr:         address.Address("localhost"),
		organization: "test_org",
		org_key:      org_key,

		router:  api.BaseRouter(),
		plugins: make(map[string]*gin.Engine),
	}

	return s
}

func (s *Server) Serve() error {
	// Listen and Server @ localhost:8080
	// NOTE: this is mapped on Dockerfile and is served to browser at 8080
	err := s.router.Run(":8080")
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

	for {
		i := int64(rand.Intn(len(test_users)))
		rs := test_users[i]
		path := fmt.Sprintf("http://localhost:8080/users/%s", rs)
		request, _ := http.NewRequest("POST", path, nil)
		client := &http.Client{}
		client.Do(request)
		time.Sleep(300 * time.Second)
	}
}
