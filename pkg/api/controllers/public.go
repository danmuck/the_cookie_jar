package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/gin-gonic/gin"
)

func PingPong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func RouteIndex(c *gin.Context) {
	if c.Query("new_user") == "true" {
		c.HTML(http.StatusOK, "dev_index.tmpl", gin.H{
			"title":     "Welcome to the_cookie_jar API!",
			"sub_title": "Learning Management System",
			"body":      "Thanks for registering",
		})
		return
	}
	c.HTML(http.StatusOK, "dev_index.tmpl", gin.H{
		"title":           "Welcome to the_cookie_jar API!",
		"sub_title":       "Learning Management System",
		"body":            "TODO",
		"register_button": "true",
	})
}

func DevIndex(c *gin.Context) {
	if c.Query("new_user") == "true" {
		c.HTML(http.StatusOK, "dev_index.tmpl", gin.H{
			"title":     "Welcome to the_cookie_jar API!",
			"sub_title": "Learning Management System",
			"body":      "Thanks for registering",
		})
		return
	}
	c.HTML(http.StatusOK, "dev_index.tmpl", gin.H{
		"title":           "Welcome to the_cookie_jar API!",
		"sub_title":       "Learning Management System",
		"body":            "TODO",
		"register_button": "true",
	})
}
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":           "Welcome to the_cookie_jar API!",
		"sub_title":       "Learning Management System",
		"body":            "TODO",
		"register_button": "true",
		"username":        c.GetString("username"),
	})
}

func GET_UserRegistration(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func POST_UserRegistration(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	password_confirm := c.PostForm("password_confirm")

	if password != password_confirm {
		e := fmt.Sprintf("/register?error=%v", "passwords do not match")
		c.Redirect(http.StatusFound, e)
		return
	}

	err := database.AddUser(username, password)
	if err != nil {
		e := fmt.Sprintf("/register?error=%v", err)
		c.Redirect(http.StatusFound, e)
		return
	}

	user, err := database.GetUser(username)
	if err != nil {
		e := fmt.Sprintf("/register?error=%v", err)
		c.Redirect(http.StatusFound, e)
		return
	}
	user.ClassroomIDs = append(user.ClassroomIDs, os.Getenv("dev_class_id"))
	err = database.UpdateUser(user)
	if err != nil {
		e := fmt.Sprintf("/register?error=%v", err)
		c.Redirect(http.StatusFound, e)
		return
	}

	classroom, err := database.GetClassroom(os.Getenv("dev_class_id"))
	if err != nil {
		e := fmt.Sprintf("/register?error=%v", err)
		c.Redirect(http.StatusFound, e)
		return
	}
	classroom.StudentIDs = append(classroom.StudentIDs, user.ID)
	err = database.UpdateClassroom(classroom)
	if err != nil {
		e := fmt.Sprintf("/register?error=%v", err)
		c.Redirect(http.StatusFound, e)
		return
	}

	c.Redirect(http.StatusFound, "/?new_user=true")
}

func GET_UserLogin(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/")

}

func POST_UserLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// If the password matches, generate an auth token
	err := database.VerifyPassword(username, password)
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=bad_password")
	}

	token, err := database.GenerateAuthToken(username)
	if err != nil {
		e := fmt.Sprintf("/login?error=%v", err)
		c.Redirect(http.StatusFound, e)
		return
	}

	c.SetCookie("jwt_token", token, int(time.Hour.Seconds()), "/", "/", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func POST_UserLogout(c *gin.Context) {
	user, err := database.GetUser(c.GetString("username"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"result": user,
		})
		return
	}

	user.Auth.AuthTokenHash = ""
	err = database.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"result": user,
		})
		return
	}

	c.SetCookie("jwt_token", "", 1, "/", "/", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}
