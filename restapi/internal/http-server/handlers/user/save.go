package user

import (
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	helper "restapi/internal/lib/helperfunctions"
	"restapi/internal/lib/password"
	"restapi/internal/lib/sl"
	"restapi/internal/storage"
	"restapi/internal/models/data"
	"restapi/internal/models/response"

	"github.com/gin-gonic/gin"
)

// SaveUser implements UserHandlers.
func (u UserHandler) SaveUser(c *gin.Context) {
	// load logger with necessary data
	const op = "handlers.user.save.SaveUserHandler"
	logger := helper.LoadLogger(u.log, c, op)

	// bind request
	var req saveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(errorset.ErrBindRequest, sl.Err(err))
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
		return
	}

	logger.Info("decoded request", slog.Any(helper.ReqKey, req))

	// validate request
	if err := validatingRequest(c, logger, req, u.db); err != nil {
		return
	}

	// action with db
	userId, err := u.db.SaveUser(req.Username, req.Password)
	if err != nil || userId == 0 {
		handleSavingUserError(c, logger, err, userId)
		return
	}

	var data data.Data = data.NewData()
	data[helper.UserIDKey] = userId

	logger.Info("user saved successfully", slog.String(helper.UsernameKey, req.Username))
	response.Ok(c, http.StatusCreated, data)
}

func validatingRequest(c *gin.Context, log *slog.Logger, req saveRequest, us storage.Storage) error {
	if !password.IsValidPassword(req.Password) {
		log.Error(errorset.ErrInvalidPassword.Error())
		response.Error(c, http.StatusBadRequest, errorset.ErrInvalidPassword.Error())
		return errorset.ErrValidation
	}

	userexists, err := us.UsernameExists(req.Username)
	if err != nil {
		log.Error("failed to check if username exists", sl.Err(err))
		response.Error(c, http.StatusInternalServerError, "failed to check if username exists")
		return errorset.ErrValidation
	}

	if userexists {
		log.Warn("username already exists")
		response.Error(c, http.StatusConflict, "username already exists")
		return errorset.ErrValidation
	}

	return nil
}

func handleSavingUserError(c *gin.Context, log *slog.Logger, err error, userId int64) {
	if err != nil {
		log.Error("failed to save user", sl.Err(err))
		response.Error(c, http.StatusInternalServerError, "failed to save user")
		return
	} else if userId == 0 {
		log.Error("unexpected user ID = 0 after saving user")
		response.Error(c, http.StatusInternalServerError, "unexpected server error")
		return
	}
}
