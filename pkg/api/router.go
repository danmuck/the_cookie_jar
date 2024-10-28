package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/danmuck/the_cookie_jar/pkg/api/controllers"
	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/middleware"
	"github.com/danmuck/the_cookie_jar/pkg/api/middleware/authorization"

	"github.com/gin-gonic/gin"
)

func DefaultClassroomSetup() {
	database.AddUser("admin", "password")
	user, _ := database.GetUser("admin")

	if len(user.ClassroomIDs) == 0 {
		database.AddClassroom(user.Username, "dev_class")
		user, _ = database.GetUser(user.Username)
	}
	classroom, _ := database.GetClassroom(user.ClassroomIDs[0])

	if len(classroom.BoardIDs) == 0 {
		database.AddBoard(classroom.ID, "dev_discussion")
		classroom, _ = database.GetClassroom(user.ClassroomIDs[0])
	}
	board, _ := database.GetBoard(classroom.BoardIDs[0])

	if len(board.ThreadIDs) == 0 {
		database.AddThread(board.ID, "dev_thread")
		board, _ = database.GetBoard(classroom.BoardIDs[0])
	}
	thread, _ := database.GetThread(board.ThreadIDs[0])

	if len(thread.CommentIDs) == 0 {
		database.AddComment(thread.ID, user.Username, "Welcome!", "Welcome to this thread. This is a default thread created for grading purposes. Feel free to comment and like! Refresh the page to see new comments.")
		database.AddComment(thread.ID, user.Username, "Welcome Part 2", "This is the second comment in our default development DefaultClassroomSetup()")
	}

	os.Setenv("dev_url", fmt.Sprintf("/classrooms/%v/discussions/%v/threads/%v", classroom.ID, board.ID, thread.ID))
	os.Setenv("dev_class_id", classroom.ID)
}

func BaseRouter() *gin.Engine {
	router := gin.Default()
	DefaultClassroomSetup()
	// Loading our templates and CSS stylesheets
	router.LoadHTMLGlob("/root/public/templates/*")
	router.Static("/public/styles", "./public/styles")
	router.Static("/public/assets", "./public/assets")
	router.StaticFile("/public/functions.js", "./public/functions.js")

	// Middleware that will be used by ALL routes
	router.Use(middleware.DefaultMiddleware())

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

	return router
}
