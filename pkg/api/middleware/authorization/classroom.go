package middleware

import (
	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/middleware"
	"github.com/danmuck/the_cookie_jar/pkg/utils"

	"github.com/gin-gonic/gin"
)

func BoardCreationMiddleware(c *gin.Context) {
	// Grabbing the user and making sure it was successful
	user, err := database.GetUser(c.GetString("Username"))
	if err != nil {
		middleware.GiveStatusForbidden(c, "There was a problem retrieving your account.")
		return
	}

	// Grabbing classroom and making sure it exists
	classroom, err := database.GetClassroom(c.Param("ClassroomID"))
	if err != nil {
		middleware.GiveStatusForbidden(c, "There was a problem retrieving your classroom.")
		return
	}

	// Is user authorized to make a board
	if user.ID != classroom.ProfessorID && !utils.Contains(classroom.InstructorIDs, user.ID) {
		middleware.GiveStatusForbidden(c, "You cannot create a discussion board.")
		return
	}

	c.Next()
}
