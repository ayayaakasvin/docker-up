package delete

import (
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	helper "restapi/internal/lib/helperfunctions"
	"restapi/internal/lib/sl"
	"restapi/internal/models/response"

	"github.com/gin-gonic/gin"
)

type TaskDeleter interface {
	DeleteTask(taskId int64) error
}

// DeleteTaskHandler deletes user record from database
func DeleteTaskHandler(log *slog.Logger, td TaskDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// load logger with necessary data
		const op = "handlers.task.delete.DeleteTaskHandler"
		helper.LoadLogger(&log, c, op)

		// fetch ID param
		taskId := helper.GetIDFromParams(c, helper.TaskIDKey)
		if taskId == -1 {
			response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
			return
		}

		log.Info("decoded request", slog.Int64(helper.TaskIDKey, taskId))

		// action with db
		err := td.DeleteTask(taskId)
		if err != nil {
			log.Error("failed to delete task", sl.Err(err))
			response.Error(c, http.StatusInternalServerError, "failed to delete task")
			return
		}

		log.Info("task deleted successfully", slog.Int64(helper.TaskIDKey, taskId))
		response.Ok(c, http.StatusOK, nil)
	}
}