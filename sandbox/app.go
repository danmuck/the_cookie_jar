package sandbox

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type App struct {
	name    string
	version string

	Router *gin.Engine
}

func (a *App) GetInfo() string {
	return fmt.Sprintf("name: %s \nversion: %s \n", a.name, a.version)
}

func NewApp(opts ...string) *App {
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

		Router: gin.Default(),
	}

	err := app.Router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}
	return app
}
