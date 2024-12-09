package middleware

import (
	"sync"
	"time"

	"github.com/danmuck/the_cookie_jar/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Limiter struct {
	RequestCount     int
	FirstRequestTime time.Time

	IsBlocked   bool
	BlockedTime time.Time
}

/*
Creates default limiter object when limiter is reset or IP requests for the
first time.
*/
func createDefaultLimiter() Limiter {
	return Limiter{
		RequestCount:     1,
		FirstRequestTime: time.Now(),
		IsBlocked:        false,
		BlockedTime:      time.Now(),
	}
}

// All open WebSockets, also a mutex to prevent race conditions
var ipLimits = make(map[string]Limiter)
var ipLimitsMutex sync.Mutex

// Per how many general requests to clear idle IP limiter mappings
const clearInterval = 2000

// How many requests have we processed so far
var requestsCounted = 0

/*
Protects against DoS attacks.
*/
func RateLimitingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIp := c.ClientIP()

		ipLimitsMutex.Lock()
		limiter, exists := ipLimits[clientIp]
		defer ipLimitsMutex.Unlock()

		if !exists {
			limiter = createDefaultLimiter()
		} else {
			timeSinceFirstRequest := time.Since(limiter.FirstRequestTime).Seconds()
			timeSinceBlock := time.Since(limiter.BlockedTime).Seconds()

			// If they are not blocked, and it's been 10 seconds, reset the
			// number of requests they have to 1.
			// If they are blocked, and it's been 30 seconds, reset them so
			// they are unblocked.
			// Otherwise, just add to their request count and see if we should
			// block them
			if !limiter.IsBlocked && timeSinceFirstRequest > 10 {
				ipLimits[clientIp] = createDefaultLimiter()
				limiter = ipLimits[clientIp]
			} else if limiter.IsBlocked && timeSinceBlock > 30 {
				ipLimits[clientIp] = createDefaultLimiter()
				limiter = ipLimits[clientIp]
			} else {
				limiter.RequestCount += 1

				// If more than 50 requests has occured within 10 seconds
				if limiter.RequestCount > 50 && timeSinceFirstRequest <= 10 {
					limiter.IsBlocked = true
					limiter.BlockedTime = time.Now()
				}
			}
		}

		ipLimits[clientIp] = limiter
		if limiter.IsBlocked {
			utils.RouteIPLimit(c, 30, "seconds")
			c.Abort()
			return
		}

		requestsCounted += 1
		// Seeing if we need to clear idle IP limits
		if requestsCounted == clearInterval {
			requestsCounted = 0

			var deleteKeys []string
			for ip, limiter := range ipLimits {
				timeSinceFirstRequest := time.Since(limiter.FirstRequestTime).Seconds()

				if timeSinceFirstRequest > 120 {
					deleteKeys = append(deleteKeys, ip)
				}
			}
			for _, key := range deleteKeys {
				delete(ipLimits, key)
			}
		}

		c.Next()
	}
}
