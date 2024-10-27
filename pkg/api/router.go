package api

import (
	"net/http"

	"github.com/danmuck/the_cookie_jar/pkg/api/controllers"
	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/middleware"
	"github.com/danmuck/the_cookie_jar/pkg/api/middleware/authorization"

	"github.com/gin-gonic/gin"
)

func DefaultClassroomSetup() {
	database.AddUser("admin", "password")
	database.AddClassroom("admin", "dev_class")

	database.AddBoard()
}

func BaseRouter() *gin.Engine {
	router := gin.Default()

	// Loading our templates and CSS stylesheets
	router.LoadHTMLGlob("/root/public/templates/*")
	router.Static("/public/styles", "./public/styles")
	router.Static("/public/assets", "./public/assets")

	// Middleware that will be used by ALL routes
	router.Use(middleware.DefaultMiddleware())

	// Non-authenticated public routes
	public := router.Group("/")
	{
		public.GET("/", controllers.Index)
		public.GET("/tmp", controllers.TestIndex)

		public.GET("/register", controllers.GET_UserRegistration)
		public.POST("/register", controllers.POST_UserRegistration)
		public.GET("/login", controllers.GET_UserLogin)
		public.POST("/login", controllers.POST_UserLogin)
	}

	// Authenticated user-data routes
	protected := router.Group("/users", middleware.UserAuthenticationMiddleware())
	{
		protected.GET("/", controllers.PingPong)
		protected.GET("/:username", controllers.GET_Username)
		protected.POST("/:username", controllers.POST_User)
		protected.PUT("/:id", controllers.PUT_User)
		protected.DELETE("/:username", controllers.DEL_User)
	}

	// '.../classrooms'
	classrooms := router.Group("/classrooms", middleware.UserAuthenticationMiddleware())
	{
		classrooms.GET("/", controllers.ClassroomIndex)
		classrooms.POST("/new", controllers.POST_Classroom)

		// '.../classrooms/ID'
		classroom := classrooms.Group("/:classroom_id", authorization.ClassroomVerificationMiddleware())
		{
			classroom.GET("/", controllers.ClassroomIndex)

			// '.../classrooms/ID/discussions'
			boards := classroom.Group("/discussions")
			{
				boards.GET("/", controllers.DiscussionIndex)
				boards.POST("/new", controllers.POST_Discussion, authorization.BoardCreationMiddleware())

				// '.../classrooms/ID/discussions/ID'
				board := boards.Group("/:board_id", authorization.BoardVerificationMiddleware())
				{
					board.GET("/")

					// '.../classrooms/ID/discussions/ID/threads
					threads := board.Group("/threads")
					{
						threads.GET("/")
						threads.POST("/new")

						// '.../classrooms/ID/discussions/ID/threads/ID
						thread := threads.Group("/:thread_id", authorization.ThreadVerificationMiddleware())
						{
							thread.GET("/")
							threads.POST("/comment")
							threads.POST("/like")
						}
					}
				}
			}

			// '.../classrooms/ID/settings'
			settings := classroom.Group("/settings")
			{
				settings.GET("/")
				settings.POST("/addstudent")
				settings.POST("/addinstructor")
			}
		}
	}

	// NOTE: Not for production use.
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
