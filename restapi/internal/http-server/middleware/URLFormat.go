package logger

import (
	"github.com/gin-gonic/gin"
)

// URLFormat is a middleware that ensures the URL is properly formatted
func URLFormat() gin.HandlerFunc {
	return func(c *gin.Context) {
		if ext := c.Request.URL.Path[len(c.Request.URL.Path)-5:]; ext == ".json" {
			c.Header("Content-Type", "application/json")
		}
		c.Next()
	}
}