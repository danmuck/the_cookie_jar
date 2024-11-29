package controllers

import (
	"net/http"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/utils"
	"github.com/gin-gonic/gin"
)

func POST_CreateClassroom(c *gin.Context) {
	if c.PostForm("class-name") == "" {
		utils.RouteError(c, "class name cannot be empty")
		return
	}

	database.AddClassroom(c.GetString("username"), c.PostForm("class-name"))
	c.Redirect(http.StatusSeeOther, "/")
}

func GET_Classroom(c *gin.Context) {
	c.HTML(http.StatusOK, "class.html", gin.H{
		"IsLoggedIn":      true,
		"Username":        c.GetString("username"),
		"IsProfessor":     c.GetBool("isClassProfessor"),
		"ClassName":       c.GetString("className"),
		"SettingsMessage": c.Query("settingsMessage"),
	})
}

func POST_LeaveClassroom(c *gin.Context) {
	if c.GetBool("isClassProfessor") {
		c.Redirect(http.StatusSeeOther, "/"+c.Param("classroom_id")+"/?settingsMessage=You+are+the+professor.")
		return
	}

	database.UpdateClassroomStudents(c.Param("classroom_id"), c.GetString("username"), true)
	c.Redirect(http.StatusSeeOther, "/")
}

func POST_AddStudent(c *gin.Context) {
	if c.PostForm("username") == c.GetString("username") {
		c.Redirect(http.StatusSeeOther, "/"+c.Param("classroom_id")+"/?settingsMessage=Cannot+add+yourself.")
		return
	}

	err := database.UpdateClassroomStudents(c.Param("classroom_id"), c.PostForm("username"), false)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/"+c.Param("classroom_id")+"/?settingsMessage=Something+went+wrong.")
	} else {
		c.Redirect(http.StatusSeeOther, "/"+c.Param("classroom_id")+"/?settingsMessage=Added+the+user!")
	}
}

func POST_RemoveStudent(c *gin.Context) {
	if c.PostForm("username") == c.GetString("username") {
		c.Redirect(http.StatusSeeOther, "/"+c.Param("classroom_id")+"/?settingsMessage=Cannot+remove+yourself.")
		return
	}

	err := database.UpdateClassroomStudents(c.Param("classroom_id"), c.PostForm("username"), true)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/"+c.Param("classroom_id")+"/?settingsMessage=Something+went+wrong.")
	} else {
		c.Redirect(http.StatusSeeOther, "/"+c.Param("classroom_id")+"/?settingsMessage=Removed+the+user!")
	}
}
