package logger

import (
	"log/slog"
	"github.com/gin-gonic/gin"
)

func RequestIDLoggerMiddleware(log *slog.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID, _ := c.Get("RequestID")
        if requestID == nil {
            requestID = "unknown"
        }

        // Add the request ID to the logger
        ctxLog := log.With(slog.String("request_id", requestID.(string)))
        c.Set("logger", ctxLog)

        c.Next()
    }
}
