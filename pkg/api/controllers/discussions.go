package controllers

import (
	"fmt"
	"net/http"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/gin-gonic/gin"
)

func DiscussionIndex(c *gin.Context) {
	classroomID := c.Param("classroom_id")
	c.HTML(http.StatusOK, "discussion_board.tmpl", gin.H{
		"title":        "Discussion Board",
		"sub_title":    "Some Classroom Name Probably",
		"body":         "Welcome to ",
		"new_board":    "true",
		"classroom_id": classroomID,
	})
}

// Threads
func GET_Thread(c *gin.Context) {
	classroomID := c.Param("classroom_id")
	threadID := c.Param("thread_id")
	thread, err := database.GetThread(threadID)
	if err != nil {
		c.HTML(http.StatusOK, "discussion_board.tmpl", gin.H{
			"error": err,
		})
	}
	all_comments := make([]models.Comment, 0)
	for _, commentID := range thread.CommentIDs {
		comment, err := database.GetComment(commentID)
		if err != nil {
			all_comments = append(all_comments, *comment)
		}
	}
	c.HTML(http.StatusOK, "discussion_board.tmpl", gin.H{
		"title":          "Discussion Board",
		"sub_title":      "Some Classroom Name Probably",
		"body":           "Welcome to ",
		"new_board":      "true",
		"classroom_id":   classroomID,
		"current_thread": thread,
		"all_comments":   all_comments,
	})
}

// Discussions
func GET_NewDiscussion(c *gin.Context) {
	err := c.Query("error")
	if err != "" {
		c.HTML(http.StatusOK, "discussion_board.tmpl", gin.H{
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

	c.HTML(http.StatusOK, "discussion_board.tmpl", gin.H{
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
