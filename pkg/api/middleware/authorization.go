package middleware

import (
	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/utils"
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
			utils.RouteError(c, "there was a problem")
			c.Abort()
			return
		}

		// Grabbing classroom and making sure it exists
		classroom, err := database.GetClassroom(c.Param("classroom_id"))
		if err != nil || !classroom.ContainsUserID(user.Username) {
			utils.RouteError(c, "this content does not exist")
			c.Abort()
			return
		}

		// Is the user authenticating as the professor?
		if classroom.IsUserIDPrivileged(user.Username) {
			c.Set("isClassProfessor", true)
		} else {
			c.Set("isClassProfessor", false)
		}

		c.Set("className", classroom.Name)
		c.Next()
	}
}

func ClassroomPrivilegedVerificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.GetBool("isClassProfessor") {
			utils.RouteError(c, "invalid permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
