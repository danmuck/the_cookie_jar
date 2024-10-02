package controllers

import (
	"net/http"

	"github.com/danmuck/the_cookie_jar/api/models"
	"github.com/gin-gonic/gin"
)

func UpdateUsername(c *gin.Context) {
	var user models.User
	username := c.Param("username")
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	user.ID = username

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
	c.String(http.StatusOK, "pong")
}

func AddUser(c *gin.Context) {
	c.String(http.StatusOK, "add user controller")
}

func DeleteUser(c *gin.Context) {
	c.String(http.StatusOK, "delete user controller")
}

func LookupUser(c *gin.Context) {
	var user models.User
	username := c.Param("username")
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	user.ID = username

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}
