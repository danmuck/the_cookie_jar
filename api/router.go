package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/danmuck/the_cookie_jar/api/controllers"
	"github.com/danmuck/the_cookie_jar/api/middleware"
	"github.com/gin-gonic/gin"
)

func wd() {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	fmt.Println("Current Directory:", dir)

	// Walk through the directory and print the structure
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Print the file or directory
		fmt.Println(path)
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}
func ServeHTML(router *gin.Engine) {
	// wd()
	router.LoadHTMLGlob("/root/public/templates/*")
	router.GET("/users/posts", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/users/info", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Info",
		})
	})
}

func BaseRouter() *gin.Engine {
	router := gin.Default()
	go ServeHTML(router)
	// Public routes
	public := router.Group("/users")
	public.Use(middleware.Logger())
	{
		public.POST("/register", controllers.PingPong)
		public.POST("/login", controllers.PingPong)
	}
	// Protected routes
	protected := router.Group("/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/", controllers.Root)
		protected.POST("/users", controllers.Root)
		protected.GET("/users/:id", controllers.Root)
		protected.PUT("/users/:id", controllers.Root)
		protected.DELETE("/users/:id", controllers.Root)
	}

	return router
}
