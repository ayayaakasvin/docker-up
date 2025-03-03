package get

import (
	"errors"
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	"restapi/internal/lib/sl"
	"restapi/internal/models"
	"restapi/internal/models/status"

	"github.com/gin-gonic/gin"
)

type Request struct {
	ID int64 `json:"id" binding:"required,gt=0"`
}

type Response struct {
	Status status.Status `json:"status"`
	User   *models.User  `json:"user,omitempty"`
}

type UserGetter interface {
	GetUserByID(id int64) (*models.User, error) 
}

func GetUserHandler (log *slog.Logger, ug UserGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.user.get.GetUserHandler"
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

		userObject, err := ug.GetUserByID(req.ID)
		if err != nil {
			if errors.Is(err, errorset.ErrUserNotFound) {
				log.Error(err.Error(), sl.Err(err))
				responseError(c, http.StatusNotFound, err.Error())
				return
			}
			
			log.Error("failed to get user", sl.Err(err))
			responseError(c, http.StatusInternalServerError, "failed to get user")
			return
		}

		responseOk(c, userObject)
	}
}

func responseOk(c *gin.Context, userObj *models.User) {
	c.JSON(http.StatusOK, Response{
		Status: status.OK(),
		User: userObj,
	})
}

func responseError(c *gin.Context, code int, errormsg string) {
	c.JSON(code, Response{
		Status: status.Error(errormsg),
	})
}