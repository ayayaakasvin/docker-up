package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware adds a request ID to the context if missing
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID") // Fetch from request header
		if requestID == "" {
			requestID = uuid.New().String() // Generate new if missing
		}

		// Store in Gin context
		c.Set("RequestID", requestID)

		// Set response header for tracking
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next() // Continue request
	}
}
