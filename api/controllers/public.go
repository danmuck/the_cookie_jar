package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danmuck/the_cookie_jar/api/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func PingPong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":     "Index",
		"sub_title": "Future Homepage",
		"body":      "Some post text from a user that was in their recent post",
	})
}

func ServeUserRegistration(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tmpl", nil)
}

func UserRegistration(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	pw_bytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pw_bytes, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	c.SetCookie("username", username, 3600, "/", "localhost", false, true)
	c.SetCookie("password", string(hash), 3600, "/", "localhost", false, true)

	note := fmt.Sprintf("username: %v password: %v", username, password)

	var user *models.User = models.NewUser(username, string(hash))
	var result *models.User
	users := get_collection("users")
	err = users.FindOne(context.TODO(), gin.H{"username": username}).Decode(&result)
	if err != nil {
		_, err = users.InsertOne(context.TODO(), user)
		if err != nil {
			fmt.Printf("insert error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User added successfully",
			"note":    note,
			"type":    "POST",
			"user":    user,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": ("exists:" + note)})
	// c.HTML(http.StatusOK, "register.tmpl", gin.H{
	// 	"username": "SUCCESS",
	// 	"password": "password",
	// 	"error":    nil,
	// })
}
