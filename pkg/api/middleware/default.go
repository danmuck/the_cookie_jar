package middleware

import (
	"github.com/gin-gonic/gin"
)

/*
All routed paths should use this middleware.
*/
func DefaultMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If user is logged in get their username, if not no biggie
		token, _ := c.Cookie("jwt_token")
		user, err := isAuthenticated(token)
		if err == nil {
			c.Set("username", user.Username)
		}

		c.Header("X-Content-Type-Options", "nosniff")
		c.Next()
	}
}
