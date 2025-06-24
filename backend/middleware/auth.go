package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Simple in-memory rate limiter implementation
var (
	requestCounts = make(map[string]int)
	lastCleared   = time.Now()
	mutex         sync.Mutex
	maxRequests   = 100 // Max requests per minute
)

// RateLimiter limits the number of requests per IP
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()

		mutex.Lock()

		// Reset counts every minute
		if time.Since(lastCleared) > time.Minute {
			requestCounts = make(map[string]int)
			lastCleared = time.Now()
		}

		// Increment request count for this IP
		requestCounts[ip]++
		count := requestCounts[ip]

		mutex.Unlock()

		// Check if rate limit exceeded
		if count > maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminAuth middleware for admin routes
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// In a real implementation, you would validate a token here
		token := c.GetHeader("Authorization")

		// Simple token check - in a real app, use JWT or similar
		if token != "admin-secret-token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
