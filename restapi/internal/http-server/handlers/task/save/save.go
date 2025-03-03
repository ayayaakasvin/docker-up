package save

import (
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	"restapi/internal/lib/sl"
	"restapi/internal/models/status"

	"github.com/gin-gonic/gin"
)

// Request
type Request struct {
    UserID      int64   `json:"user_id" binding:"required,gt=0"`
    TaskContent string  `json:"task_content" binding:"required"`
}

// Response
type Response struct {
	Status status.Status    `json:"status"`
	UserID int64            `json:"user_id,omitempty"`
    TaskID int64            `json:"task_id,omitempty"`
}

type TaskSaver interface {
	SaveTask(userId int64, content string) (int64, error)
}

// SaveTaskHandler saves a new user
func SaveTaskHandler(log *slog.Logger, ts TaskSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.task.save.SaveTaskHandler"
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

		id, err := ts.SaveTask(req.UserID, req.TaskContent)
		if err != nil {
			log.Error("failed to save task", sl.Err(err))
            if err == errorset.ErrUserNotFound {
			    responseError(c, http.StatusConflict, errorset.ErrUserNotFound.Error()) 
				return
            }
			responseError(c, http.StatusInternalServerError, "failed to save task")
			return
		}
		if id == 0 {
			log.Error("unexpected task ID = 0 after saving task")
			responseError(c, http.StatusInternalServerError, "unexpected server error")
			return
		}

		log.Info("task saved successfully", slog.Int64("user_id", req.UserID))
		log.Info("task saved successfully", slog.Int64("task_id", id))
		responseOk(c, req.UserID, id)
	}
}

func responseOk(c *gin.Context, user_id, task_id int64) {
	c.JSON(http.StatusCreated, Response{
		Status: status.OK(),
		UserID: user_id,
        TaskID: task_id,
	})
}

func responseError(c *gin.Context, code int, errormsg string) {
	c.JSON(code, Response{
		Status: status.Error(errormsg),
	})
}