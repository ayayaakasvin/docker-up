package task

import (
	"errors"
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	helper "restapi/internal/lib/helperfunctions"
	"restapi/internal/lib/sl"
	"restapi/internal/models/data"
	"restapi/internal/models/response"
	"restapi/internal/models/task"

	"github.com/gin-gonic/gin"
)

// GetTasksByUserID implements TaskHandlers.
func (t TaskHandler) GetTasksByUserID(c *gin.Context) {
	// load logger with necessary data
	const op = "handlers.task.TaskHandler.GetTasksByUserID"
	logger := helper.LoadLogger(t.log, c, op)

	// fetch ID param
	userId := helper.GetIDFromParams(c, helper.UserIDKey)
	if userId == -1 {
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
		return
	}

	logger.Info("decoded request", slog.Any(helper.UserIDKey, userId))

	// action with db
	tasksSlice, err := t.db.GetTasksByUserID(userId)
	if err != nil || len(tasksSlice) == 0 {
		handleGettingTasksError(c, logger, err, tasksSlice, userId)
		return
	}

	var data data.Data = data.NewData()
	data[helper.TasksKey] = tasksSlice

	logger.Info("tasks succesfully passed", slog.Int64(helper.UserIDKey, userId))
	response.Ok(c, http.StatusOK, data)
}

func handleGettingTasksError(c *gin.Context, log *slog.Logger, err error, tasksSlice []*task.Task, userId int64) {
	if errors.Is(err, errorset.ErrUserNotFound) {
		log.Error(err.Error(), sl.Err(err))
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	if len(tasksSlice) == 0 {
		log.Warn("no tasks found for user", slog.Int64(helper.UserIDKey, userId))
		response.Error(c, http.StatusNotFound, errorset.ErrTaskNotFound.Error())
		return
	}
}
