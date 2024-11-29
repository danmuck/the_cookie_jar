package controllers

import (
	"net/http"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/danmuck/the_cookie_jar/pkg/utils"
	"github.com/gin-gonic/gin"
)

func GET_Homepage(c *gin.Context) {
	if c.GetString("username") != "" {
		createdClassList := make([]models.Classroom, 0)
		joinedClassList := make([]models.Classroom, 0)

		user, err := database.GetUser(c.GetString("username"))
		if err != nil {
			utils.RouteError(c, "there was a problem")
			return
		}

		for _, classroomId := range user.ClassroomIDs {
			classroom, err := database.GetClassroom(classroomId)
			if err != nil {
				utils.RouteError(c, "there was a problem")
				return
			}

			if classroom.ProfessorID == user.Username {
				createdClassList = append(createdClassList, *classroom)
			} else {
				joinedClassList = append(joinedClassList, *classroom)
			}
		}

		c.HTML(http.StatusOK, "classlist.html", gin.H{
			"IsLoggedIn":       true,
			"Username":         user.Username,
			"CreatedClassList": createdClassList,
			"JoinedClassList":  joinedClassList,
		})
		return
	}

	c.HTML(http.StatusOK, "login_register.html", gin.H{
		"IsLoggedIn":             false,
		"SuccessfullyRegistered": c.Query("register") == "true",
	})
}
