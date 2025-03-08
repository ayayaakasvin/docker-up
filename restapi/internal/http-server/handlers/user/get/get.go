package get

import (
	"errors"
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	helper "restapi/internal/lib/helperfunctions"
	"restapi/internal/lib/sl"
	"restapi/internal/models"
	"restapi/internal/models/data"
	"restapi/internal/models/response"

	"github.com/gin-gonic/gin"
)

type UserGetter interface {
	GetUserByID(userId int64) (*models.User, error)
}

func GetUserHandler(log *slog.Logger, ug UserGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.user.get.GetUserHandler"
		helper.LoadLogger(&log, c, op)

		userId := helper.GetIDFromParams(c, helper.UserIDKey)
		if userId == -1 {
			response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
			return
		}

		log.Info("decoded request", slog.Any(helper.UserIDKey, userId))

		userObject, err := ug.GetUserByID(int64(userId))
		if err != nil {
			handleGettingUserError(c, log, err)
			return
		}

		var data data.Data = data.NewData()
		data[helper.UserKey] = userObject

		response.Ok(c, http.StatusOK, data)
	}
}

func handleGettingUserError(c *gin.Context, log *slog.Logger, err error)  {
	if errors.Is(err, errorset.ErrUserNotFound) {
		log.Error(err.Error(), sl.Err(err))
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	log.Error("failed to get user", sl.Err(err))
	response.Error(c, http.StatusInternalServerError, "failed to get user")
	return
}