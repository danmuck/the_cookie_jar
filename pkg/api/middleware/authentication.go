package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	fmt.Println(">> [middleware] loaded AuthMiddleware .. ")
	return func(c *gin.Context) {
		// Check if the user is authenticated
		if isAuthenticated(c) {
			c.Header("X-Content-Type-Options", "nosniff")
			c.Next()
			return
		}
		// User is not authenticated, return an error response
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}

func isAuthenticated(c *gin.Context) bool {
	fmt.Println(">> [middleware] checking isAuthenticated .. ")
	for k, v := range c.Keys {
		fmt.Printf(">> [middleware] k: %v \nv: %v", k, v)
	}
	// Check if the user is authenticated based on a JWT token, session, or any other mechanism
	// Return true if the user is authenticated, false otherwise
	return true
}
