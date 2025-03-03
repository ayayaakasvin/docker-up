package save

import (
	"log/slog"
	"net/http"

	"restapi/internal/lib/password"
	"restapi/internal/lib/sl"
	"restapi/internal/models/status"

	"github.com/gin-gonic/gin"
)

// Request
type Request struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
}

// Response
type Response struct {
	Status status.Status `json:"status"`
	UserID int64         `json:"user_id,omitempty"`
}

type UserSaver interface {
	SaveUser(username, password string) (int64, error)
	UsernameExists(username string) (bool, error)
}

// SaveUserHandler saves a new user
func SaveUserHandler(log *slog.Logger, us UserSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.user.save.SaveUserHandler"
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

		userexists, err := us.UsernameExists(req.Username)
		if err != nil {
			log.Error("failed to check if username exists", sl.Err(err))
			responseError(c, http.StatusInternalServerError, "failed to check if username exists")
			return
		}

		if userexists {
			log.Warn("username already exists")
			responseError(c, http.StatusConflict, "username already exists")
			return
		}

		id, err := us.SaveUser(req.Username, req.Password)
		if err != nil {
			log.Error("failed to save user", sl.Err(err))
			responseError(c, http.StatusInternalServerError, "failed to save user")
			return
		}
		if id == 0 {
			log.Error("unexpected user ID = 0 after saving user")
			responseError(c, http.StatusInternalServerError, "unexpected server error")
		}

		log.Info("user saved successfully", slog.String("username", req.Username))
		responseOk(c, id)
	}
}

func responseOk(c *gin.Context, id int64) {
	c.JSON(http.StatusCreated, Response{
		Status: status.OK(),
		UserID: id,
	})
}

func responseError(c *gin.Context, code int, errormsg string) {
	c.JSON(code, Response{
		Status: status.Error(errormsg),
	})
}