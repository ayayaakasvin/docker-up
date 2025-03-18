package task

import (
	"log/slog"
	"restapi/internal/storage"

	"github.com/gin-gonic/gin"
)

type TaskHandlers interface {
	DeleteTask(c *gin.Context)
	GetTaskByTaskID(c *gin.Context)
	GetTasksByUserID(c *gin.Context)
	UpdateTask(c *gin.Context)
	SaveTask(c *gin.Context)
}

type TaskHandler struct {
	log *slog.Logger
	db  storage.Storage
}

func NewTaskHandler(log *slog.Logger, db storage.Storage) TaskHandlers {
	return TaskHandler{
		log: log,
		db: db,
	}
}

type request struct {
	TaskContent string `json:"taskContent" binding:"required"`
}
