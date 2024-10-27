package authorization

import (
	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

/*
Relies on: `UserAuthenticationMiddleware()`
*/
func ClassroomVerificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Grabbing the user and making sure it was successful
		user, err := database.GetUser(c.GetString("username"))
		if err != nil {
			middleware.GiveStatusForbidden(c, "There was a problem retrieving your account.")
			return
		}

		// Grabbing classroom and making sure it exists
		classroom, err := database.GetClassroom(c.Param("classroom_id"))
		if err != nil || !classroom.ContainsUserID(user.ID) {
			middleware.GiveStatusForbidden(c, "There was a problem retrieving your classroom.")
			return
		}

		c.Next()
	}
}
