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

	"github.com/gin-gonic/gin"
)

// GetTaskByTaskID implements TaskHandlers.
func (t TaskHandler) GetTaskByTaskID(c *gin.Context) {
	// load logger with necessary data
	const op = "handlers.task.TaskHandler.GetTaskByTaskID"
	logger := helper.LoadLogger(t.log, c, op)

	// fetch ID param
	taskID := helper.GetIDFromParams(c, helper.TaskIDKey)
	if taskID == -1 {
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
		return
	}

	logger.Info("decoded request", slog.Any("req", taskID))

	// action with db
	task, err := t.db.GetTaskByTaskID(taskID)
	if err != nil {
		handleGettingTaskError(c, logger, err)
		return
	}

	var data data.Data = data.NewData()
	data[helper.TaskKey] = task

	logger.Info("task succesfully passed",
		slog.Int64(helper.UserIDKey, task.UserID),
		slog.Int64(helper.TaskIDKey, task.TaskID))

	response.Ok(c, http.StatusFound, data)
}

func handleGettingTaskError(c *gin.Context, log *slog.Logger, err error) {
	log.Error("failed to get task", sl.Err(err))
	if errors.Is(err, errorset.ErrTaskNotFound) {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Error(c, http.StatusInternalServerError, "failed to get task")
}
