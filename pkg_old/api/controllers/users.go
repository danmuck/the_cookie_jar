package controllers

import (
	"context"
	"net/http"
	"os"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func POST_User(c *gin.Context) {
	username := c.Param("username")
	password := c.Param("password")
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "cannot use blank password",
			"result": "",
		})
		return
	}
	err := database.AddUser(username, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    err.Error(),
			"type":     "POST",
			"who":      username,
			"password": password,
		})
		return
	}

	user, err := database.GetUser(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    err.Error(),
			"type":     "POST",
			"who":      username,
			"password": password,
		})
		return
	}
	user.ClassroomIDs = append(user.ClassroomIDs, os.Getenv("dev_class_id"))
	err = database.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    err.Error(),
			"type":     "POST",
			"who":      username,
			"password": password,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User added successfully",
		"type":    "POST",
	})
}

func DEL_User(c *gin.Context) {
	id := c.Query("id")
	coll := database.GetCollection("users")
	filter := bson.M{"_id": id}

	var result models.User
	err := coll.FindOneAndDelete(context.TODO(), filter).Decode(&result)
	if err != nil {
		c.String(http.StatusNotFound, "User does not exist")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"type":    "DELETE",
		"who":     result.Username,
		"user":    result,
	})
}

func GET_Username(c *gin.Context) {
	username := c.Param("username")
	user, err := database.GetUser(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User found successfully",
		"type":    "GET",
		"who":     username,
		"user":    user,
	})
}

func PUT_User(c *gin.Context) {
	id := c.Param("id")
	coll := database.GetCollection("users")
	filter := bson.M{"_id": id}

	var user models.User
	err := coll.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	queries := []string{"username", "password", "org", "status"}
	for _, param := range queries {
		value := c.Query(param)
		if value == "" {
			continue
		}
		switch param {
		case "username":
			user.Username = value
		case "password":
			user.Auth.PasswordHash = ""
		default:
		}
	}

	result, err := coll.ReplaceOne(context.TODO(), filter, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"result": user,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"result":  result,
		"user":    user,
	})
}
