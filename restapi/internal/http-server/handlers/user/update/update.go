package update

import (
	"errors"
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	"restapi/internal/lib/password"
	"restapi/internal/lib/sl"
	"restapi/internal/models/status"

	"github.com/gin-gonic/gin"
)

type Request struct {
	ID 			int64 	`json:"id" binding:"required,gt=0"`
	Password	string	`json:"password" binding:"required,min=8"`
}

type Response struct {
	Status status.Status 	`json:"status"`
	UserID int64 			`json:"user_id,omitempty"`
}

type UserUpdater interface {
	UpdateUserPassword(id int64, password string) (error)
}

func UpdateUserPasswordHandler (log *slog.Logger, uu UserUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.user.update.UpdateUserHandler"
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

		if !password.IsValidPassword(req.Password) {
			log.Error("Invalid password")
			responseError(c, http.StatusBadRequest, "invalid password")
			return
		}

		err := uu.UpdateUserPassword(req.ID, req.Password)
		if err != nil {
			if errors.Is(err, errorset.ErrUserNotFound) {
				log.Error("user not found", sl.Err(err))
				responseError(c, http.StatusNotFound, "user not found")
				return
			}
			
			log.Error("failed to update user password", sl.Err(err))
			responseError(c, http.StatusInternalServerError, "failed to update user password")
			return
		}

		responseOk(c, req.ID)
	}
}

func responseOk(c *gin.Context, id int64) {
	c.JSON(http.StatusOK, Response{
		Status: status.OK(),
		UserID: id,
	})
}

func responseError(c *gin.Context, code int, errormsg string) {
	c.JSON(code, Response{
		Status: status.Error(errormsg),
	})
}