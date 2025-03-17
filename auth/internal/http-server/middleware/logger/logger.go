package logger

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func New(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		entry := log.With(
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("remote_addr", c.Request.RemoteAddr),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("request_id", c.GetString("request_id")),	
		)

		log.Info("logger middleware enabled")

		t1 := time.Now()

		c.Next()

		entry.Info("request completed",
			slog.Int("status", c.Writer.Status()),
			slog.Int("bytes", c.Writer.Size()),
			slog.String("duration", time.Since(t1).String()),
		)
	}
}