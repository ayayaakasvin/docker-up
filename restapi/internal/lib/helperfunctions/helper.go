package helper

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	TaskIDKey 	= "taskId"
	UserIDKey 	= "userId"
	TaskKey 	= "task"
	TasksKey 	= "tasks"
	UserKey 	= "user"
	ReqKey 		= "request"
	UsernameKey = "username"
)

func GetIDFromParams(c *gin.Context, idkey string) int64 {
	if idkey == "" {
		return -1
	}

	taskIDString := c.Param(idkey)
	if taskIDString == "" {
		return -1
	}
	taskID, err := strconv.ParseInt(taskIDString, 10, 0)
	if err != nil {
		return -1
	}

	return taskID
}

func LoadLogger(log **slog.Logger, c *gin.Context, operation string)  {
	requestID, exists := c.Get("RequestID")
	if !exists {
		requestID = "unknown"
	}

	*log = (*log).With(
		slog.String("op", operation),
		slog.String("request_id", requestID.(string)),
	)
}