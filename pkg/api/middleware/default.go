package middleware

import (
	"github.com/gin-gonic/gin"
)

/*
All routed paths should use this middleware.
*/
func DefaultMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Header("X-Content-Type-Options", "nosniff")

		// If user is logged in get their username, if not no biggie
		token, _ := c.Cookie("jwt_token")
		user, err := isAuthenticated(token)
		if err == nil {
			c.Set("username", user.Username)
		} else {
			c.Set("username", "")
		}

		c.Next()
	}
}
