package handlers

import (
	"log/slog"

	"github.com/ayayaakasvin/auth/internal/storage"
	"github.com/gin-gonic/gin"
)

type AppHandlers interface {
	Authentificate(*gin.Context)
	CreateUser(*gin.Context)
}

type AppHandler struct {
	db  storage.Storage
	log *slog.Logger
}

func NewAppHandler(db storage.Storage, log *slog.Logger) AppHandlers {
	return &AppHandler{
		db: db,
		log: log,
	}
}

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}