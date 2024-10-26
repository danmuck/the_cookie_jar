package controllers

import (
	"fmt"
	"net/http"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/gin-gonic/gin"
)

func ClassroomIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "classroom.tmpl", gin.H{
		"title":     "Classroom",
		"sub_title": "INDEX [tmp]",
		"body":      "New classroom stuff or information.",
		"new_board": "true",
	})
}

func GET_NewClassroom(c *gin.Context) {
	err := c.Query("error")
	if err != "" {
		c.HTML(http.StatusOK, "classroom.tmpl", gin.H{
			"error": err,
		})

	}
	c.HTML(http.StatusOK, "classroom.tmpl", gin.H{
		"title":         "Classroom -- New",
		"sub_title":     "Some Classroom Name Probably",
		"body":          "Welcome to the class",
		"new_classroom": "true",
	})
}

func POST_Classroom(c *gin.Context) {
	name := c.PostForm("name")
	username := c.GetString("username")
	if username != "" {
		err := database.AddClassroom(username, name)
		if err != nil {
			e := fmt.Sprintf("/classrooms/new?error=%v", err)
			c.Redirect(http.StatusNotFound, e)
		}
	}

	c.HTML(http.StatusOK, "classroom.tmpl", gin.H{
		"title":         name,
		"sub_title":     username,
		"body":          "Welcome to the class",
		"new_classroom": "false",
	})
}

func PUT_Classroom(c *gin.Context) {

}

func DELETE_Classroom(c *gin.Context) {

}
