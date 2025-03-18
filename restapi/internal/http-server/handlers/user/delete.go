package user

import (
	"errors"
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	helper "restapi/internal/lib/helperfunctions"
	"restapi/internal/lib/sl"
	"restapi/internal/models/response"

	"github.com/gin-gonic/gin"
)

// DeleteUser implements UserHandlers.
func (u UserHandler) DeleteUser(c *gin.Context) {
	// load logger with necessary data
	const op = "handlers.user.delete.DeleteUserHandler"
	logger := helper.LoadLogger(u.log, c, op)

	// fetch ID param
	userId := helper.GetIDFromParams(c, helper.UserIDKey)
	if userId == -1 {
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
		return
	}

	logger.Info("decoded request", slog.Any(helper.UserIDKey, userId))

	// action with db
	err := u.db.DeleteUser(userId)
	if err != nil {
		handleDeletingUserError(c, logger, err)
		return
	}

	response.Ok(c, http.StatusOK, nil)
}

func handleDeletingUserError(c *gin.Context, log *slog.Logger, err error) {
	if errors.Is(err, errorset.ErrUserNotFound) {
		log.Error(errorset.ErrUserNotFound.Error(), sl.Err(err))
		response.Error(c, http.StatusNotFound, errorset.ErrUserNotFound.Error())
		return
	}

	log.Error("failed to delete user", sl.Err(err))
	response.Error(c, http.StatusInternalServerError, "failed to delete user")
}
