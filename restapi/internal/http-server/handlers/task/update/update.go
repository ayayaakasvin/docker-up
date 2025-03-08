package update

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

type Request struct {
	TaskContent string `json:"taskContent" binding:"required"`
}

type TaskUpdater interface {
	UpdateTaskContent(taskId int64, content string) error
}

func UpdateTaskHandler(log *slog.Logger, tu TaskUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		// load logger with necessary data
		const op = "handlers.task.update.UpdateTaskHandler"
		helper.LoadLogger(&log, c, op)

		// fetch ID param
		taskID := helper.GetIDFromParams(c, helper.TaskIDKey)
		if taskID == -1 {
			response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
			return
		}

		// bind request
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(errorset.ErrBindRequest, sl.Err(err))
			response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
			return
		}

		log.Info("decoded request", slog.Any(helper.ReqKey, req))

		// action with db
		err := tu.UpdateTaskContent(int64(taskID), req.TaskContent)
		if err != nil {
			log.Error("failed to update task", sl.Err(err))
			response.Error(c, http.StatusInternalServerError, "failed to update task")
			return
		}

		var data data.Data = data.NewData()
		data[helper.TaskIDKey] = taskID

		log.Info("task updated successfully", slog.Int64(helper.TaskIDKey, taskID))
		response.Ok(c, http.StatusOK, data)
	}
}