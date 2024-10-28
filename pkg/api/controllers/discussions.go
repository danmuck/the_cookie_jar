package controllers

import (
	"fmt"
	"net/http"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/danmuck/the_cookie_jar/pkg/utils"
	"github.com/gin-gonic/gin"
)

// "/classrooms/77a43ace-c76f-48bc-9778-9bc91a6a8709/discussions/7c023439-4240-4342-8782-774a68073831/threads/d6323b64-6357-41da-a8da-e4d438f9c36a/"

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
		comment, _ := database.GetComment(commentID)
		all_comments = append(all_comments, *comment)
	}

	c.HTML(http.StatusOK, "discussion_board.tmpl", gin.H{
		"title":          "Discussion Board",
		"sub_title":      "Some Classroom Name Probably",
		"body":           "Welcome to ",
		"new_board":      "true",
		"classroom_id":   classroomID,
		"current_thread": thread,
		"all_comments":   all_comments,
		"username":       c.GetString("username"),
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

func POST_Comment(c *gin.Context) {
	user, _ := database.GetUser(c.GetString("username"))
	classroomID := c.Param("classroom_id")
	boardID := c.Param("board_id")
	threadID := c.Param("thread_id")
	title := c.PostForm("title")
	text := c.PostForm("text")

	database.AddComment(threadID, user.Username, text, title)
	e := fmt.Sprintf("/classrooms/%v/discussions/%v/threads/%v", classroomID, boardID, threadID)
	c.Redirect(http.StatusSeeOther, e)
}

func POST_CommentLike(c *gin.Context) {
	user, _ := database.GetUser(c.GetString("username"))
	classroomID := c.Param("classroom_id")
	boardID := c.Param("board_id")
	threadID := c.Param("thread_id")
	commentID := c.PostForm("comment_id")

	comment, err := database.GetComment(commentID)
	if err != nil {
		c.HTML(http.StatusOK, "discussion_board.tmpl", gin.H{
			"error": err,
		})
	}

	if utils.Contains(comment.LikedUsers, user.ID) {
		comment.LikedUsers = utils.RemoveItem(comment.LikedUsers, user.Username)
	} else {
		comment.LikedUsers = append(comment.LikedUsers, user.Username)
	}
	_ = database.UpdateComment(comment)

	e := fmt.Sprintf("/classrooms/%v/discussions/%v/threads/%v", classroomID, boardID, threadID)
	c.Redirect(http.StatusSeeOther, e)
}
