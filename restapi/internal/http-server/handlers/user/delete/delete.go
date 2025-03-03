package delete

import (
	"errors"
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	"restapi/internal/lib/sl"
	"restapi/internal/models/status"

	"github.com/gin-gonic/gin"
)

type Request struct {
	ID 			int64 	`json:"id" binding:"required,gt=0"`
}

type Response struct {
	Status status.Status 	`json:"status"`
}

type UserDeleter interface {
	DeleteUser(id int64) (error)
}

func DeleteUserHandler (log *slog.Logger, ud UserDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.user.delete.DeleteUserHandler"
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

		err := ud.DeleteUser(req.ID)
		if err != nil {
			if errors.Is(err, errorset.ErrUserNotFound) {
				log.Error("user not found", sl.Err(err))
				responseError(c, http.StatusNotFound, "user not found")
				return
			}
			
			log.Error("failed to delete user", sl.Err(err))
			responseError(c, http.StatusInternalServerError, "failed to delete user")
			return
		}

		responseOk(c)
	}
}

func responseOk(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Status: status.OK(),
	})
}

func responseError(c *gin.Context, code int, errormsg string) {
	c.JSON(code, Response{
		Status: status.Error(errormsg),
	})
}