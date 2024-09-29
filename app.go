package main

// import "github.com/gin-gonic/gin"
import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	id       []byte
	username string
	status   string
}

func (u *User) getId() []byte {
	// go naming conventions make methods that start with a lowercase letter are private
	return u.id
}

func NewUser() User {
	return User{
		id:       []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09},
		username: "Big",
		status:   "Chillin",
	}
}

type App struct {
	name    string
	version string

	router *gin.Engine
}

func (a *App) GetInfo() string {
	return fmt.Sprintf("name: %s \nversion: %s \n", a.name, a.version)
}
func newApp(opts ...string) *App {
	// args are a slice of strings since go does not implement optional args
	// access them by index
	if len(opts) > 2 {
		fmt.Println(">> args: name, version \n(all else ignored)")
	}

	app := &App{
		name: opts[0],
		// functions expressed this way are essentially lambda functions
		// they are executed in place (note: () follows the closing brace)
		version: func() string {
			if len(opts) > 1 {
				return opts[1]
			} else {
				return "dev"
			}
		}(),

		router: gin.Default(),
	}

	return app
}

func main() {

	new_user := NewUser()
	var tmp int = 12
	fmt.Fprintf(os.Stderr, "Sup, %s --%s age: %d \n", new_user.username, new_user.status, tmp)
	fmt.Println("Using goroutines, iterate the bytes of [new_user.id] and print them ..\n .. they are indexed in the order the goroutines were created")

	for i, id := range new_user.id {
		go func() {
			// by calling a function using [go nameOfFunction()] a goroutine is created
			// goroutines run concurrently
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("%d: %d \n", i, id)
		}()
	}

	// sleep to wait for the goroutines to finish
	time.Sleep(2 * time.Second)
	fmt.Println("id: ", new_user.getId())

	// initialize App
	app := newApp("the_cookie_jar")
	fmt.Println(app.GetInfo())
	// add route to app router
	app.router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Listen and Server in 0.0.0.0:8080
	app.router.Run(":8080")
}
