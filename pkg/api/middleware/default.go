package middleware

import (
	"github.com/gin-gonic/gin"
)

/*
All routed paths should use this middleware.
*/
func DefaultMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Next()
	}
}
