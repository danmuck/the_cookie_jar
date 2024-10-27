package main

import (
	"fmt"

	"github.com/danmuck/the_cookie_jar/pkg/api"
	"github.com/gin-gonic/gin"
)

func main() {
	// Want Gin debug statements? Comment this line out.
	gin.SetMode(gin.ReleaseMode)

	router := api.BaseRouter()
	fmt.Println("-----------------------------------")
	fmt.Println("the_cookie_jar server is running...")
	fmt.Println("-----------------------------------")
	err := router.Run(":8080")

	// If a general error in the server occurs
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}
	fmt.Println("-----------------------------------")
	fmt.Println("the_cookie_jar server is closing...")
	fmt.Println("-----------------------------------")
}
