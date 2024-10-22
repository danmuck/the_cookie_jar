package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

/*
	TODO:
	if user is not logged in, redirect to login page

*/

func POST_User(c *gin.Context) {
	username := c.Param("username")
	password := c.Param("password")
	if password == "" {
		password = "pass@!word"
	}

	user, err := models.NewUser(username, password)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":    err.Error(),
            "type":     "POST",
            "who":      username,
            "password": password,
        })
        return
    }
	var result *models.User
	users := database.GetCollection("users")
	err = users.FindOne(context.TODO(), gin.H{"username": username}).Decode(&result)
	if err != nil {
		_, err = users.InsertOne(context.TODO(), user)
		if err != nil {
			fmt.Printf("insert error: %v", err)
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
			"user":    user,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "User exists."})
}

func DEL_User(c *gin.Context) {
	id := c.Query("id")
	coll := database.GetCollection("users")
	filter := bson.M{"_id": id}

	var result models.User
	err := coll.FindOneAndDelete(context.TODO(), filter).Decode(&result)
	if err != nil {
		c.String(http.StatusNotFound, "User does not exist")
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
	coll := database.GetCollection("users")
	filter := bson.M{"username": username}

	var result models.User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User found successfully",
		"type":    "GET",
		"who":     username,
		"user":    result,
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
			user.UpdatePassword(value)
		case "org":
			user.Org = value
		case "status":
			user.UpdateStatus(value)
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
