package delete

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
}

// Response
type Response struct {
	Status status.Status    `json:"status"`
}

type TaskDeleter interface {
	DeleteTask(task_id int64) error
}

// SaveTaskHandler saves a new user
func DeleteTaskHandler(log *slog.Logger, td TaskDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.task.delete.DeleteTaskHandler"
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

		err := td.DeleteTask(req.TaskID)
		if err != nil {
			log.Error("failed to delete task", sl.Err(err))
			responseError(c, http.StatusInternalServerError, "failed to delete task")
			return
		}

		log.Info("task deleted successfully", slog.Int64("task_id", req.TaskID))
		responseOk(c)
	}
}

func responseOk(c *gin.Context) {
	c.JSON(http.StatusCreated, Response{
		Status: status.OK(),
	})
}

func responseError(c *gin.Context, code int, errormsg string) {
	c.JSON(code, Response{
		Status: status.Error(errormsg),
	})
}