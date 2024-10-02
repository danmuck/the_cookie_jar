package routers

import (
	"net/http"

	"github.com/danmuck/the_cookie_jar/api/controllers"
	"github.com/danmuck/the_cookie_jar/api/middleware"
	"github.com/gin-gonic/gin"
)

func ServeHTML(router *gin.Engine) {
	router.LoadHTMLGlob("public/templates/**/*")
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
