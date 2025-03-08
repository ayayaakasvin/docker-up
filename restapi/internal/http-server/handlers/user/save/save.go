package save

import (
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

// Request
type Request struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserSaver interface {
	SaveUser(username, password string) (int64, error)
	UsernameExists(username string) (bool, error)
}

// SaveUserHandler saves a new user
func SaveUserHandler(log *slog.Logger, us UserSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		// load logger with necessary data
		const op = "handlers.user.save.SaveUserHandler"
		helper.LoadLogger(&log, c, op)

		// bind request
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(errorset.ErrBindRequest, sl.Err(err))
			response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
			return
		}

		log.Info("decoded request", slog.Any(helper.ReqKey, req))

		// validate request
		if err := validatingRequest(c, log, req, us); err != nil {
			return
		}

		// action with db
		userId, err := us.SaveUser(req.Username, req.Password)
		if err != nil || userId == 0 {
			handleSavingUserError(c, log, err, userId)
			return
		}

		var data data.Data = data.NewData()
		data[helper.UserIDKey] = userId

		log.Info("user saved successfully", slog.String(helper.UsernameKey, req.Username))
		response.Ok(c, http.StatusCreated, data)
	}
}

func validatingRequest(c *gin.Context, log *slog.Logger, req Request, us UserSaver) error {
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

func handleSavingUserError(c *gin.Context, log *slog.Logger, err error, userId int64)  {
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