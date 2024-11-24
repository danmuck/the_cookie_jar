package authorization

import (
	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

/*
Relies on: `ClassroomVerificationMiddleware()`
*/
func BoardVerificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Grabbing board and making sure it exists
		_, err := database.GetBoard(c.Param("board_id"))
		if err != nil {
			middleware.GiveStatusForbidden(c, "There was a problem retrieving your board.")
			return
		}

		c.Next()
	}
}

/*
Relies on: `ClassroomVerificationMiddleware()`
*/
func BoardCreationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Grabbing the user and making sure it was successful
		user, err := database.GetUser(c.GetString("username"))
		if err != nil {
			middleware.GiveStatusForbidden(c, "There was a problem retrieving your account.")
			return
		}

		// Grabbing classroom and making sure it exists
		classroom, err := database.GetClassroom(c.Param("classroom_id"))
		if err != nil {
			middleware.GiveStatusForbidden(c, "There was a problem retrieving your classroom.")
			return
		}

		// Is user authorized to make a board
		if classroom.IsUserIDPrivileged(user.ID) {
			middleware.GiveStatusForbidden(c, "You cannot create a discussion board.")
			return
		}

		c.Next()
	}
}

/*
Relies on: `BoardVerificationMiddleware()`
*/
func ThreadVerificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Grabbing thread and making sure it exists
		_, err := database.GetThread(c.Param("thread_id"))
		if err != nil {
			middleware.GiveStatusForbidden(c, "There was a problem retrieving your thread.")
			return
		}

		c.Next()
	}
}
