package get

import (
	"errors"
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	"restapi/internal/lib/sl"
	"restapi/internal/models"
	"restapi/internal/models/status"

	"github.com/gin-gonic/gin"
)

type Request struct {
	UserID int64 `json:"user_id" binding:"required,gt=0"`
}

type Response struct {
	Status status.Status `json:"status"`
	Tasks []*models.Task `json:"tasks,omitempty"`
}

type TasksGetter interface {
	GetTasksByUserID(id int64) ([]*models.Task, error)
}

func GetTasksHandler (log *slog.Logger, tg TasksGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.task.get.GetTasksHandler"
		requestID, exists := c.Get("RequestID")
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

		tasksSlice, err := tg.GetTasksByUserID(req.UserID)
		if err != nil {
			if errors.Is(err, errorset.ErrUserNotFound) {
				log.Error(err.Error(), sl.Err(err))
				responseError(c, http.StatusNotFound, err.Error())
				return
			}
			
			log.Error("failed to get user", sl.Err(err))
			responseError(c, http.StatusInternalServerError, "failed to get user")
			return
		}

		if len(tasksSlice) == 0 {
			responseError(c, http.StatusNotFound, errorset.ErrTaskNotFound.Error())
			return
		}

		responseOk(c, tasksSlice)
	}
}

func responseOk(c *gin.Context, taskSlice []*models.Task) {
	c.JSON(http.StatusOK, Response{
		Status: status.OK(),
		Tasks: taskSlice,
	})
}

func responseError(c *gin.Context, code int, errormsg string) {
	c.JSON(code, Response{
		Status: status.Error(errormsg),
	})
}