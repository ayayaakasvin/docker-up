package update

import (
	"errors"
	"log/slog"
	"net/http"
	"restapi/internal/errorset"
	helper "restapi/internal/lib/helperfunctions"
	"restapi/internal/lib/password"
	"restapi/internal/lib/sl"
	"restapi/internal/models/data"
	"restapi/internal/models/response"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Password string `json:"password" binding:"required,min=8"`
}

type UserUpdater interface {
	UpdateUserPassword(userId int64, password string) error
}

func UpdateUserPasswordHandler(log *slog.Logger, uu UserUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		// load logger with necessary data
		const op = "handlers.user.update.UpdateUserHandler"
		helper.LoadLogger(&log, c, op)

		// fetch ID param
		userId := helper.GetIDFromParams(c, helper.UserIDKey)
		if userId == -1 {
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

		// password validation
		if !password.IsValidPassword(req.Password) {
			log.Error(errorset.ErrInvalidPassword.Error())
			response.Error(c, http.StatusBadRequest, errorset.ErrInvalidPassword.Error())
			return
		}

		// action with db
		err := uu.UpdateUserPassword(userId, req.Password)
		if err != nil {
			handleUpdatingUserError(c, log, err)
			return
		}

		var data data.Data = data.NewData()
		data[helper.UserIDKey] = userId

		response.Ok(c, http.StatusOK, data)
	}
}

func handleUpdatingUserError(c *gin.Context, log *slog.Logger, err error) {
	if errors.Is(err, errorset.ErrUserNotFound) {
		log.Error(errorset.ErrUserNotFound.Error(), sl.Err(err))
		response.Error(c, http.StatusNotFound, errorset.ErrUserNotFound.Error())
		return
	}

	log.Error("failed to update user password", sl.Err(err))
	response.Error(c, http.StatusInternalServerError, "failed to update user password")
	return
}