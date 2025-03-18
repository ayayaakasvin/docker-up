package user

import (
	"log/slog"
	"restapi/internal/storage"

	"github.com/gin-gonic/gin"
)

type UserHandlers interface {
	DeleteUser(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUserPassword(c *gin.Context)
	SaveUser(c *gin.Context)
}

type UserHandler struct {
	log *slog.Logger
	db  storage.Storage
}

func NewUserHandler(log *slog.Logger, db storage.Storage) UserHandlers {
	return UserHandler{
		log: log,
		db:  db,
	}
}

type saveRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
}

type updateRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}
