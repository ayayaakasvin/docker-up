package logger

import (
	"github.com/gin-gonic/gin"
)

// URLFormat is a middleware that ensures the URL is properly formatted
func URLFormat() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		
		if len(path) >= 5 && path[len(path)-5:] == ".json" {
			c.Header("Content-Type", "application/json")
		}
		
		c.Next()
	}
}
