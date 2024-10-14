package api

import (
	"net/http"

	"github.com/danmuck/the_cookie_jar/pkg/api/controllers"
	"github.com/danmuck/the_cookie_jar/pkg/api/middleware"
	"github.com/gin-gonic/gin"
)

func ServeHTML_demo(router *gin.Engine) {
	router.GET("/users/posts", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":     "User Posts [tmp]",
			"sub_title": "Some Recent Post Maybe",
			"body":      "Some post text from a user that was in their recent post",
		})
	})
	router.GET("/users/info/:username", func(c *gin.Context) {
		username := c.Param("username")
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":     "User Info",
			"sub_title": username,
			"body":      "Some user info.",
		})
	})
}

func BaseRouter() *gin.Engine {
	// init router
	router := gin.Default()
	// load templates (will error if none exist at path)
	router.LoadHTMLGlob("/root/public/templates/*")    // load templates
	router.Static("/public/styles", "./public/styles") // load css stylesheets

	go ServeHTML_demo(router)

	// Public routes
	public := router.Group("/")
	public.Use(middleware.Logger())
	{
		public.GET("/", controllers.Index)
		public.POST("/", controllers.Index)
		public.GET("/register", controllers.GET_UserRegistration)
		public.POST("/register", controllers.POST_UserRegistration)
		public.GET("/login", controllers.GET_UserLogin)
		public.POST("/login", controllers.POST_UserLogin)
	}
	// Protected routes
	protected := router.Group("/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/", controllers.PingPong)

		protected.GET("/:username", controllers.GET_Username)
		protected.POST("/:username", controllers.POST_User)
		protected.PUT("/:id", controllers.PUT_User)
		protected.DELETE("/:username", controllers.DEL_User)
	}

	dev := router.Group("/dev")
	{
		dev.GET("/routes", func(c *gin.Context) {
			routes := router.Routes()
			type tmp struct {
				Method string `json:"Method"`
				Path   string `json:"Path"`
			}
			var t []tmp
			for _, route := range routes {
				r := tmp{
					Path:   route.Path,
					Method: route.Method,
				}
				t = append(t, r)
			}
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"routes": t,
			})
		})
	}

	return router
}
