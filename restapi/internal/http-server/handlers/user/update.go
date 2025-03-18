package user

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

// UpdateUserPassword implements UserHandlers.
func (u UserHandler) UpdateUserPassword(c *gin.Context) {
	// load logger with necessary data
	const op = "handlers.user.update.UpdateUserHandler"
	logger := helper.LoadLogger(u.log, c, op)

	// fetch ID param
	userId := helper.GetIDFromParams(c, helper.UserIDKey)
	if userId == -1 {
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
		return
	}

	// bind request
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(errorset.ErrBindRequest, sl.Err(err))
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
		return
	}

	logger.Info("decoded request", slog.Any(helper.ReqKey, req))

	// password validation
	if !password.IsValidPassword(req.Password) {
		logger.Error(errorset.ErrInvalidPassword.Error())
		response.Error(c, http.StatusBadRequest, errorset.ErrInvalidPassword.Error())
		return
	}

	// action with db
	err := u.db.UpdateUserPassword(userId, req.Password)
	if err != nil {
		handleUpdatingUserError(c, logger, err)
		return
	}

	var data data.Data = data.NewData()
	data[helper.UserIDKey] = userId

	response.Ok(c, http.StatusOK, data)
}

func handleUpdatingUserError(c *gin.Context, log *slog.Logger, err error) {
	if errors.Is(err, errorset.ErrUserNotFound) {
		log.Error(errorset.ErrUserNotFound.Error(), sl.Err(err))
		response.Error(c, http.StatusNotFound, errorset.ErrUserNotFound.Error())
		return
	}

	log.Error("failed to update user password", sl.Err(err))
	response.Error(c, http.StatusInternalServerError, "failed to update user password")
}
