package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danmuck/the_cookie_jar/api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

/*
	TODO:
	if user is not logged in, redirect to login page
	PUT update user

*/

func POST_user(c *gin.Context) {
	username := c.Param("username")
	password := c.Param("password")
	if password == "" {
		password = "pass@!word"
	}

	dev := fmt.Sprintf("username: %v password: %v", username, password)

	var user *models.User = models.NewUser(username, password)
	var result *models.User
	users := get_collection("users")
	err := users.FindOne(context.TODO(), gin.H{"username": username}).Decode(&result)
	if err != nil {
		_, err = users.InsertOne(context.TODO(), user)
		if err != nil {
			fmt.Printf("insert error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User added successfully",
			"who":     dev,
			"type":    "POST",

			"user": user,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "User exists."})
}

func DEL_user(c *gin.Context) {
	id := c.Query("id")

	coll := get_collection("users")
	filter := bson.M{"_id": id}

	var result models.User
	err := coll.FindOneAndDelete(context.TODO(), filter).Decode(&result)
	if err != nil {
		c.String(http.StatusNotFound, "User does not exist")
	}

	dev := fmt.Sprintf("[DEL_user] username: %v", result.Username)
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"who":     dev,
		"type":    "DELETE",
		"user":    result,
	})
}

func GET_username(c *gin.Context) {
	username := c.Param("username")
	dev := fmt.Sprintf("[GET_username] username: %v", username)

	// logic to look up user from mongodb
	coll := get_collection("users")
	filter := bson.M{"username": username}

	var result models.User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User found successfully",
		"who":     dev,
		"type":    "GET",
		"user":    result,
	})
}

func PUT_username(c *gin.Context) {
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
	// c.String(http.StatusOK, "pong")
}
