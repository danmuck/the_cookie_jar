package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/danmuck/the_cookie_jar/pkg/utils"
	"github.com/gin-gonic/gin"
)

func isAuthenticated(token string) (*models.User, error) {
	// TO-DO: CHECK IF THE TOKEN IS EXPIRED

	// Hashing the given token
	hasher := sha256.New()
	hasher.Write([]byte(token))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	// Trying to grab a user that has that hashed token
	var user *models.User
	err := database.GetCollection("users").FindOne(context.TODO(), gin.H{"Auth.AuthTokenHash": tokenHash}).Decode(&user)
	return user, err
}

func UserAuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("jwt_token")
		if err != nil {
			utils.RouteError(c, "there was a problem verifying your account, please try again")
			c.Abort()
			return
		}

		user, err := isAuthenticated(token)
		if err != nil {
			utils.RouteError(c, "there was a problem verifying your account, please try again")
			c.Abort()
			return
		}

		c.Set("username", user.Username)
		c.Next()
	}
}
