package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	fmt.Println(">> loaded Logger [middleware] .. ")
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	fmt.Println(">> loaded AuthMiddleware [middleware] .. ")
	return func(c *gin.Context) {
		// Check if the user is authenticated
		if isAuthenticated(c) {
			c.Next()
			return
		}
		// User is not authenticated, return an error response
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}
func isAuthenticated(c *gin.Context) bool {
	fmt.Println(">> checking isAuthenticated .. ")
	for k, v := range c.Keys {
		fmt.Printf("k: %v \nv: %v", k, v)
	}
	// Check if the user is authenticated based on a JWT token, session, or any other mechanism
	// Return true if the user is authenticated, false otherwise
	return true
}
