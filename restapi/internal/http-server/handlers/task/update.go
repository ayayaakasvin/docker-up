package task

import (
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	helper "restapi/internal/lib/helperfunctions"
	"restapi/internal/lib/sl"
	"restapi/internal/models/data"
	"restapi/internal/models/response"

	"github.com/gin-gonic/gin"
)

// UpdateTask implements TaskHandlers.
func (t TaskHandler) UpdateTask(c *gin.Context) {
	// load logger with necessary data
	const op = "handlers.task.TaskHandler.UpdateTask"
	logger := helper.LoadLogger(t.log, c, op)

	// fetch ID param
	taskID := helper.GetIDFromParams(c, helper.TaskIDKey)
	if taskID == -1 {
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
		return
	}

	// bind request
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(errorset.ErrBindRequest, sl.Err(err))
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
		return
	}

	logger.Info("decoded request", slog.Any(helper.ReqKey, req))

	// action with db
	err := t.db.UpdateTaskContent(int64(taskID), req.TaskContent)
	if err != nil {
		logger.Error("failed to update task", sl.Err(err))
		response.Error(c, http.StatusInternalServerError, "failed to update task")
		return
	}

	var data data.Data = data.NewData()
	data[helper.TaskIDKey] = taskID

	logger.Info("task updated successfully", slog.Int64(helper.TaskIDKey, taskID))
	response.Ok(c, http.StatusOK, data)
}
