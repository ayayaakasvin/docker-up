package user

import (
	"errors"
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	helper "restapi/internal/lib/helperfunctions"
	"restapi/internal/lib/sl"
	"restapi/internal/models/data"
	"restapi/internal/models/response"

	"github.com/gin-gonic/gin"
)

// GetUser implements UserHandlers.
func (u UserHandler) GetUser(c *gin.Context) {
	// load logger with necessary data
	const op = "handlers.user.get.GetUserHandler"
	logger := helper.LoadLogger(u.log, c, op)

	// fetch ID param
	userId := helper.GetIDFromParams(c, helper.UserIDKey)
	if userId == -1 {
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
		return
	}

	logger.Info("decoded request", slog.Any(helper.UserIDKey, userId))

	// action with db
	userObject, err := u.db.GetUserByID(int64(userId))
	if err != nil {
		handleGettingUserError(c, logger, err)
		return
	}

	var data data.Data = data.NewData()
	data[helper.UserKey] = userObject

	response.Ok(c, http.StatusOK, data)
}

func handleGettingUserError(c *gin.Context, log *slog.Logger, err error) {
	if errors.Is(err, errorset.ErrUserNotFound) {
		log.Error(err.Error(), sl.Err(err))
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	log.Error("failed to get user", sl.Err(err))
	response.Error(c, http.StatusInternalServerError, "failed to get user")
}
