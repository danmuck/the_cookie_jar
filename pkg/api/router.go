package api

import (
	"github.com/danmuck/the_cookie_jar/pkg/api/controllers"
	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/middleware"
	"github.com/danmuck/the_cookie_jar/pkg/utils"
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
	router.Static("/public/uploads", "./public/uploads")

	// Default user PFP
	database.AddDefaultMedia("default", "public/assets/default_pfp.jpg", "")

	router.GET("/", controllers.GET_Homepage)
	router.POST("/register", controllers.POST_UserRegister)
	router.POST("/login", controllers.POST_UserLogin)
	router.POST("/logout", controllers.POST_UserLogout)
	router.GET("/pfp/:account_id", controllers.GET_AccountPFP)

	router.GET("/account", middleware.UserAuthenticationMiddleware(), controllers.GET_Account)
	router.POST("/account-pfp-upload", middleware.UserAuthenticationMiddleware(), controllers.POST_AccountPFPUpload)
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

		gameRoutes := classroomRoutes.Group("/class-game")
		{
			gameRoutes.GET("/", controllers.GET_Game)
			gameRoutes.GET("/ws", controllers.GET_GameWebSocket)
		}
	}

	router.NoRoute(func(c *gin.Context) {
		utils.RouteError(c, "this content does not exist")
	})
	return router
}
