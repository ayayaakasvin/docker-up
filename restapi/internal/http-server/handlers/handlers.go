package handlers

import (
	"log/slog"
	"restapi/internal/http-server/handlers/task"
	"restapi/internal/http-server/handlers/user"
	"restapi/internal/storage"
)

type Handlers struct {
	Task task.TaskHandlers
	User user.UserHandlers
}

func NewHandlers(db storage.Storage, log *slog.Logger) *Handlers {
	return &Handlers{
		Task: task.NewTaskHandler(log, db),
		User: user.NewUserHandler(log, db),
	}
}
