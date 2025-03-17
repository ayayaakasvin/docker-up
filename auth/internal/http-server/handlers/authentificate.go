package handlers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ayayaakasvin/auth/internal/errorset"
	"github.com/ayayaakasvin/auth/internal/lib/sl"
	"github.com/ayayaakasvin/auth/internal/models/data"
	"github.com/ayayaakasvin/auth/internal/models/response"
	"github.com/ayayaakasvin/auth/internal/storage"
	"github.com/gin-gonic/gin"
)

// Authentificate implements AppHandlers.
func (a *AppHandler) Authentificate(c *gin.Context) {
	const op = "auth.internal.http-server.handlers.AppHandler.Authentificate"
	a.log = a.log.With(
		slog.String("op", op),
	)

	var request Request
	if err := c.BindJSON(&request); err != nil {
		a.log.Error(errorset.ErrBindRequest.Error(), sl.Err(err))
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest.Error())
		return
	}

	a.log.Info("decoded request", sl.Any("req", request))

	if err := validateAuthRequest(c, a.log, request, a.db); err != nil {
		a.log.Error(op, sl.Err(err))
		return
	}
}

func validateAuthRequest (c *gin.Context, log *slog.Logger, req Request, db storage.Storage) error {
	token, err := db.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error(errorset.ErrUserNotFound.Error(), sl.Err(sql.ErrNoRows))
			response.Error(c, http.StatusInternalServerError, errorset.ErrUserNotFound.Error())
			return err
		}
		log.Error(errorset.ErrAuthentificateUser.Error(), sl.Err(err))
		response.Error(c, http.StatusInternalServerError, errorset.ErrAuthentificateUser.Error())
		return err
	}

	data := data.NewDate()
	data["token"] = token

	response.Ok(c, http.StatusOK, data)
	return nil
}