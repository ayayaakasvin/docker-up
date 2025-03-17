package app

import (
	"log/slog"
	"net/http"

	"github.com/ayayaakasvin/auth/internal/config"
	"github.com/ayayaakasvin/auth/internal/http-server/handlers"
	"github.com/ayayaakasvin/auth/internal/http-server/middleware"
	"github.com/ayayaakasvin/auth/internal/http-server/middleware/logger"
	"github.com/ayayaakasvin/auth/internal/storage"
	"github.com/gin-gonic/gin"
)

func App (storage storage.Storage, log *slog.Logger, cfg *config.Config) error {
	server := setupServer(*cfg, log, storage)
	log.Info("Serving on address", slog.String("address", cfg.Address))
	return server.ListenAndServe()
}

func setupRouter (db storage.Storage, log *slog.Logger, cfg config.ServiceAddresses) *gin.Engine {
	router := gin.Default()

	middleware.LoadRouterWithMiddleware(router, 
		middleware.CorsWithConfig(cfg), 
		logger.URLFormat(),
		logger.New(log),
		logger.RequestIDLoggerMiddleware(log),
	)
	
	appHandlers := handlers.NewAppHandler(db, log)

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello World!")
	})
	router.POST("/log-in", appHandlers.Authentificate)
	router.POST("/sign-in", appHandlers.CreateUser)

	return router
}

func setupServer(cfg config.Config, log *slog.Logger, db storage.Storage) *http.Server {
	router := setupRouter(db, log, cfg.ServiceAddresses)
	log.Info("Router was set up")

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IddleTimeout,
	}

	return server
}