package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"github.com/danmuck/the_cookie_jar/pkg/api/middleware"
	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	fmt.Println("Raw password during registration:", password)

    user, err := models.NewUser(username, password)
    if err != nil {
        e := fmt.Sprintf("/register?error=%v", err)
        c.Redirect(http.StatusFound, e)
        return
    }

    var result *models.User
    users := database.GetCollection("users")
    
    err = users.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
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

    err := users.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
    if err != nil {
        c.Redirect(http.StatusFound, "/login?error=no_user")
        return
    }
	fmt.Println("Stored Hash:", result.Auth.Hash)
	fmt.Println("Password entered during login:", password)
	//If the password matches, generate an auth token
    if result.VerifyPassword(password) {
        token, err := middleware.GenToken()
        if err != nil {
            c.Redirect(http.StatusFound, "/login?error=authentication_error")
            return
        }
		//Setting expiration and storing the auth token in the database
        expiration := time.Now().Add(1 * time.Hour)
        _, err = users.UpdateOne(
            context.TODO(),
            bson.M{"username": username},
            bson.M{
                "$set": bson.M{
                    "auth.HashedToken":     middleware.HashToken(token),
                    "auth.TokenExpiration": expiration,
                },
            },
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store authentication token"})
            return
        }

        c.SetCookie("auth_token", token, int(time.Hour.Seconds()), "/", "localhost", false, true)
        c.Redirect(http.StatusFound, "/dashboard") // Change this to the correct path after login
    } else {
        c.Redirect(http.StatusFound, "/login?error=bad_password")
    }
}
