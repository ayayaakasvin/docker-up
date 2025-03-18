package task

import (
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	helper "restapi/internal/lib/helperfunctions"
	"restapi/internal/lib/sl"
	"restapi/internal/models/response"

	"github.com/gin-gonic/gin"
)

// DeleteTask implements TaskHandlers.
func (t TaskHandler) DeleteTask(c *gin.Context) {
	// load logger with necessary data
	const op = "handlers.task.TaskHandler.DeleteTask"
	logger := helper.LoadLogger(t.log, c, op)

	// fetch ID param
	taskId := helper.GetIDFromParams(c, helper.TaskIDKey)
	if taskId == -1 {
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
		return
	}

	logger.Info("decoded request", slog.Int64(helper.TaskIDKey, taskId))

	// action with db
	err := t.db.DeleteTask(taskId)
	if err != nil {
		logger.Error("failed to delete task", sl.Err(err))
		response.Error(c, http.StatusInternalServerError, "failed to delete task")
		return
	}

	logger.Info("task deleted successfully", slog.Int64(helper.TaskIDKey, taskId))
	response.Ok(c, http.StatusOK, nil)
}
