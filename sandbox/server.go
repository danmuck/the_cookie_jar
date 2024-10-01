package sandbox

import (
	"fmt"
	"net/http"
	"os"

	"github.com/danmuck/the_cookie_jar/sandbox/db_types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/address"
)

type Server struct {
	addr         address.Address
	organization string
	org_key      string
	router       *gin.Engine
	db           *Database
}

func initServer(org_key string) Server {
	db, err := NewDatabase()
	if err != nil {
		panic(err)
	}
	s := Server{
		addr:         address.Address("localhost"),
		organization: "test_org",
		org_key:      org_key,

		db:     db,
		router: gin.Default(),
	}

	return s
}

func (s *Server) DB_AddUser(user db_types.User) {
	s.db.UpdateUser(user)

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

func RunServer(key string) (*Server, error) {
	s := initServer(key)
	err := s.router.SetTrustedProxies(nil)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, ">> Server Initialized .. \n>> addr: %v \n>> org: %v \n>> key: %v", s.addr, s.organization, s.org_key)

	s.router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return &s, nil
}
