package controllers

import (
	"fmt"
	"net/http"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/gin-gonic/gin"
)

func DiscussionIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "discussions.tmpl", gin.H{
		"title":     "Discussion Board",
		"sub_title": "Some Classroom Name Probably",
		"body":      "Welcome to ",
		"new_board": "true",
	})
}

func GET_NewDiscussion(c *gin.Context) {
	err := c.Query("error")
	if err != "" {
		c.HTML(http.StatusOK, "discussions.tmpl", gin.H{
			"error": err,
		})
	}
}

func POST_Discussion(c *gin.Context) {
	name := c.PostForm("name")
	classroomID := c.GetString("ClassroomID")
	if classroomID != "" {
		err := database.AddBoard(classroomID, name)
		if err != nil {
			e := fmt.Sprintf("/classrooms/discussions/%v/new?error=%v", classroomID, err)
			c.Redirect(http.StatusNotFound, e)
		}
	}

	c.HTML(http.StatusOK, "discussions.tmpl", gin.H{
		"title":     "Discussion Board",
		"sub_title": "Some Classroom Name Probably",
		"body":      "Welcome to ",
		"new_board": "true",
	})
}

func PUT_Discussion(c *gin.Context) {

}

func DELETE_Discussion(c *gin.Context) {

}
