package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	fmt.Println(">> [middleware] loaded Logger .. ")

	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("status", "using Logger() middleware")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)
		fmt.Printf(">> [middleware] Latency: %v \n", latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
		fmt.Printf(">> [middleware] Status: %v \n", status)
	}
}
