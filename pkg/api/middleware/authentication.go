package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/gin-gonic/gin"
)

func GiveStatusForbidden(c *gin.Context, body string) {
	c.HTML(http.StatusForbidden, "index.tmpl", gin.H{
		"title":     "Welcome to the_cookie_jar API!",
		"sub_title": "Learning Management System",
		"body":      body,
		"error":     "Unauthorized",
	})
	c.Abort()
}

func isAuthenticated(token string) (*models.User, error) {
	// TO-DO: CHECK IF THE TOKEN IS EXPIRED

	// Hashing the given token
	hasher := sha256.New()
	hasher.Write([]byte(token))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	// Trying to grab a user that has that hashed token
	var user *models.User
	err := database.GetCollection("users").FindOne(context.TODO(), gin.H{"auth.hashed_token": tokenHash}).Decode(&user)
	return user, err
}

func UserAuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("jwt_token")
		if err != nil {
			GiveStatusForbidden(c, "There was a problem verifying your account, please try logging in again.")
			return
		}

		user, err := isAuthenticated(token)
		if err != nil {
			GiveStatusForbidden(c, "There was a problem verifying your account, please try logging in again.")
			return
		}

		c.Set("username", user.Username)
		c.Next()
	}
}
