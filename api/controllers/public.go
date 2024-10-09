package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danmuck/the_cookie_jar/api/database"
	"github.com/danmuck/the_cookie_jar/api/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func PingPong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func RouteIndex(c *gin.Context) {
	if c.Query("new_user") == "true" {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":     "Welcome to the_cookie_jar API!",
			"sub_title": "Learning Management System",
			"body":      "Thanks for registering",
		})
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":           "Welcome to the_cookie_jar API!",
		"sub_title":       "Learning Management System",
		"body":            "TODO",
		"register_button": "true",
	})
}

func Index(c *gin.Context) {
	if c.Query("new_user") == "true" {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":     "Welcome to the_cookie_jar API!",
			"sub_title": "Learning Management System",
			"body":      "Thanks for registering",
		})
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":           "Welcome to the_cookie_jar API!",
		"sub_title":       "Learning Management System",
		"body":            "TODO",
		"register_button": "true",
	})
}

func GET_UserRegistration(c *gin.Context) {
	err := c.Query("error")

	c.HTML(http.StatusOK, "register.tmpl", gin.H{
		"error": err,
	})
}

func POST_UserRegistration(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	pw_bytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pw_bytes, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	c.SetCookie("username", username, 3600, "/", "localhost", false, true)
	c.SetCookie("password", string(hash), 3600, "/", "localhost", false, true)

	var user *models.User = models.NewUser(username, string(hash))
	var result *models.User
	users := database.GetCollection("users")
	err = users.FindOne(context.TODO(), gin.H{"username": username}).Decode(&result)
	if err != nil {
		_, err = users.InsertOne(context.TODO(), user)
		if err != nil {
			e := fmt.Sprintf("/register?error=%v", err)
			c.Redirect(http.StatusFound, e)
			return
		}
		c.Redirect(http.StatusFound, "/?new_user=true")
		return
	}
	c.Redirect(http.StatusFound, "/register?error=username_taken")
}

func GET_UserLogin(c *gin.Context) {
	err := c.Query("error")

	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title":     "Welcome!",
		"sub_title": "Login please!",
		"error":     err,
	})
}

func POST_UserLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var result *models.User
	users := database.GetCollection("users")
	err := users.FindOne(context.TODO(), gin.H{"username": username}).Decode(&result)
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=no_user")
		return
	}
	if result.VerifyPassword(password) {
		c.Redirect(http.StatusFound, "/?login=true")
		c.SetCookie("username", result.Username, 3600, "/", "localhost", false, true)
		c.SetCookie("password", result.Auth.Hash, 3600, "/", "localhost", false, true)
		return

	}
	c.Redirect(http.StatusFound, "/login?error=bad_password")
}
