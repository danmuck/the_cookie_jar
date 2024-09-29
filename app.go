package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danmuck/the_cookie_jar/sandbox"
	"github.com/gin-gonic/gin"
)

func main() {

	new_user := sandbox.NewUser()
	var tmp int = 12
	fmt.Fprintf(os.Stderr, "\n\nSup, %s --%s age: %d \n\n", new_user.GetUsername(), new_user.GetStatus_String(), tmp)
	fmt.Println("Using goroutines, iterate the bytes of [new_user.id] and print them ..\n .. they are indexed in the order the goroutines were created .. ")

	for i, id := range new_user.GetId() {
		go func() {
			// by calling a function using [go nameOfFunction()] a goroutine is created
			// goroutines run concurrently
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("%d: %d \n", i, id)
		}()
	}

	// sleep to wait for the goroutines to finish
	time.Sleep(2 * time.Second)
	fmt.Println("id: ", new_user.GetId())

	// initialize App passing <name> string only, taking default <version>
	app := sandbox.NewApp("the_cookie_jar")
	fmt.Println(app.GetInfo())
	// add route to app router
	app.Router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	fmt.Println("\nServer connecting .. ctrl-c to quit\n ")
	// Listen and Server in 0.0.0.0:8080
	app.Router.Run(":6669")

	database, err := sandbox.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	db := database.Client.Database("the_cookie_jar")
	fmt.Println("Database: ", db.Name())

}
