package update

import (
	"log/slog"
	"net/http"

	"restapi/internal/lib/sl"
	"restapi/internal/models/status"

	"github.com/gin-gonic/gin"
)

// Request
type Request struct {
	TaskID      int64   `json:"task_id" binding:"required,gt=0"`
    TaskContent string  `json:"task_content" binding:"required"`
}

// Response
type Response struct {
	Status status.Status    `json:"status"`
    TaskID int64            `json:"task_id,omitempty"`
}

type TaskUpdater interface {
	UpdateTaskContent(task_id int64, content string) error
}

// SaveTaskHandler saves a new user
func UpdateTaskHandler(log *slog.Logger, tu TaskUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.task.update.UpdateTaskHandler"
		requestID, exists := c.Get("request_id")
		if !exists {
			requestID = "unknown"
		}

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", requestID.(string)),
		)

		var req Request
		// Reads the body of the request and binds it to the Request struct
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("failed to bind request", sl.Err(err))
			responseError(c, http.StatusBadRequest, "failed to bind request")
			return
		}

		log.Info("decoded request", slog.Any("req", req))

		err := tu.UpdateTaskContent(req.TaskID, req.TaskContent)
		if err != nil {
			log.Error("failed to update task", sl.Err(err))
			responseError(c, http.StatusInternalServerError, "failed to update task")
			return
		}

		log.Info("task updated successfully", slog.Int64("task_id", req.TaskID))
		responseOk(c, req.TaskID)
	}
}

func responseOk(c *gin.Context, task_id int64) {
	c.JSON(http.StatusOK, Response{
		Status: status.OK(),
        TaskID: task_id,
	})
}

func responseError(c *gin.Context, code int, errormsg string) {
	c.JSON(code, Response{
		Status: status.Error(errormsg),
	})
}