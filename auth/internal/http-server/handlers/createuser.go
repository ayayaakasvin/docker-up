package handlers

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/ayayaakasvin/auth/internal/errorset"
	"github.com/ayayaakasvin/auth/internal/lib/proxy"
	"github.com/ayayaakasvin/auth/internal/lib/sl"
	"github.com/ayayaakasvin/auth/internal/models/response"
	"github.com/gin-gonic/gin"
)

// CreateUser implements AppHandlers.
func (a *AppHandler) CreateUser(c *gin.Context) {
	const op = "auth.internal.http-server.handlers.AppHandler.CreateUser"
	logger := a.log.With(
		slog.String("op", op),
	)

	var request Request
	if err := c.BindJSON(&request); err != nil {
		logger.Error(errorset.ErrBindRequest.Error(), sl.Err(err))
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest.Error())
		return
	}

	if err := validateSignInRequest(c, logger, request); err != nil {
		logger.Error(op, sl.Err(err))
		return
	}
}

func validateSignInRequest(c *gin.Context, log *slog.Logger, req Request) error {
	targetURL := "http://backend:8069/user"

	log.Info("Redirecting to another service", slog.String("url", targetURL))

	proxyRequest, err := proxy.NewProxyRequest(targetURL, req, c.Request.Method)
	if err != nil {
		log.Error(errorset.ErrRequestCreate.Error(), sl.Err(err))
		response.Error(c, http.StatusInternalServerError, errorset.ErrRequestCreate.Error())
		return err
	}

	proxyRequest.Header = c.Request.Header

	client := &http.Client{}
	resp, err := client.Do(proxyRequest)
	if err != nil {
		log.Error(errorset.ErrRequestSend.Error(), sl.Err(err))
		response.Error(c, http.StatusInternalServerError, errorset.ErrRequestSend.Error())
		return err
	}
	defer resp.Body.Close()

	c.Writer.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)

	return nil
}