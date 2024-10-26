package controllers

import (
	"fmt"
	"net/http"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/gin-gonic/gin"
)

func DiscussionIndex(c *gin.Context) {
	classroomID := c.Param("classroom_id")
	c.HTML(http.StatusOK, "discussions.tmpl", gin.H{
		"title":        "Discussion Board",
		"sub_title":    "Some Classroom Name Probably",
		"body":         "Welcome to ",
		"new_board":    "true",
		"classroom_id": classroomID,
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
	classroomID := c.Param("classroom_id")

	classroom, err := database.GetClassroom(classroomID)
	if err != nil {
		e := fmt.Sprintf("/classrooms/%v/discussions/new?error=%v", classroomID, err)
		c.Redirect(http.StatusTemporaryRedirect, e)
	}
	err = database.AddBoard(classroom.ID, name)
	if err != nil {
		e := fmt.Sprintf("/classrooms/%v/discussions/new?error=%v", classroom.ID, err)
		c.Redirect(http.StatusTemporaryRedirect, e)
	}

	c.HTML(http.StatusOK, "discussions.tmpl", gin.H{
		"title":        classroom.Name,
		"sub_title":    classroom.ID,
		"body":         "Welcome to class",
		"new_board":    "false",
		"classroom_id": classroom.ID,
	})
}

func PUT_Discussion(c *gin.Context) {

}

func DELETE_Discussion(c *gin.Context) {

}
