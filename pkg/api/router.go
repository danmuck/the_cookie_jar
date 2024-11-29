package api

import (
	"github.com/danmuck/the_cookie_jar/pkg/api/controllers"
	"github.com/danmuck/the_cookie_jar/pkg/api/middleware"
	"github.com/gin-gonic/gin"
)

func BaseRouter() *gin.Engine {
	router := gin.Default()

	// Middleware that will be used by ALL routes
	router.Use(middleware.DefaultMiddleware())

	// Loading our templates and CSS stylesheets
	router.LoadHTMLGlob("/root/public/templates/*")
	router.Static("/public/styles", "./public/styles")
	router.Static("/public/assets", "./public/assets")
	router.Static("/public/scripts", "./public/scripts")

	router.GET("/", controllers.GET_Homepage)
	router.POST("/register", controllers.POST_UserRegister)
	router.POST("/login", controllers.POST_UserLogin)
	router.POST("/logout", controllers.POST_UserLogout)

	router.GET("/account", middleware.UserAuthenticationMiddleware(), controllers.GET_Account)
	router.POST("/create-classroom", middleware.UserAuthenticationMiddleware(), controllers.POST_CreateClassroom)

	classroomRoutes := router.Group("/:classroom_id", middleware.UserAuthenticationMiddleware(), middleware.ClassroomVerificationMiddleware())
	{
		classroomRoutes.GET("/", controllers.GET_Classroom)
		classroomRoutes.POST("/leave", controllers.POST_LeaveClassroom)
		classroomRoutes.POST("/add", controllers.POST_AddStudent)
		classroomRoutes.POST("/remove", controllers.POST_RemoveStudent)

		threadsRoutes := classroomRoutes.Group("/discussion-board")
		{
			threadsRoutes.GET("/", controllers.GET_Threads)
			threadsRoutes.GET("/ws", controllers.GET_ThreadsWebSocket)
			commentRoutes := threadsRoutes.Group("/:thread_id")
			{
				commentRoutes.GET("/", controllers.GET_Comments)
				commentRoutes.GET("/ws", controllers.GET_CommentsWebSocket)
			}
		}
	}

	return router

	/*
		router.POST("/media", controllers.UploadMedia)
		router.GET("/media/:id", controllers.GetMedia)
		router.DELETE("/media/:id", controllers.DeleteMedia)

		// Non-authenticated public routes
		public := router.Group("/")
		{
			public.GET("/", controllers.Index)
			public.GET("/tmp", controllers.DevIndex)

			public.GET("/register", controllers.GET_UserRegistration)
			public.POST("/register", controllers.POST_UserRegistration)
			public.GET("/login", controllers.GET_UserLogin)
			public.POST("/login", controllers.POST_UserLogin)
			public.POST("/logout", middleware.UserAuthenticationMiddleware(), controllers.POST_UserLogout)
			public.GET("/classrooms/discussions", middleware.UserAuthenticationMiddleware(), func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, os.Getenv("dev_url")) })
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
								thread.GET("/", controllers.GET_Thread)
								thread.POST("/comment", controllers.POST_Comment)
								thread.POST("/like", controllers.POST_CommentLike)
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
				c.HTML(http.StatusOK, "dev_index.tmpl", gin.H{
					"routes": t,
				})
			})
		}

		return router*/
}
