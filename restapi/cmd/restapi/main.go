package main

import (
	"log/slog"
	"net/http"
	"os"

	"restapi/internal/config"
	userDeleter "restapi/internal/http-server/handlers/user/delete"
	userGetter "restapi/internal/http-server/handlers/user/get"
	userSaver "restapi/internal/http-server/handlers/user/save"
	userUpdater "restapi/internal/http-server/handlers/user/update"
	taskSaver "restapi/internal/http-server/handlers/task/save"
	taskGetter "restapi/internal/http-server/handlers/task/get"
	taskUpdater "restapi/internal/http-server/handlers/task/update"
	taskDeleter "restapi/internal/http-server/handlers/task/delete"
	logger "restapi/internal/http-server/middleware"
	"restapi/internal/lib/logger/handlers/prettyslog"
	"restapi/internal/lib/sl"
	"restapi/internal/storage"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

const (
	envProd  = "prod"
	envDev   = "dev"
	envLocal = "local"
)

func main() {
	cfg := config.MustLoadConfig()

	log := setupLogger(cfg.Env)

	log.Info("starting RESTful API")
	log.Debug("debug message", slog.String("env", cfg.Env))

	db, err := storage.NewPostgresStorage(cfg)
	if err != nil {
		log.Error("failed to setup storage", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	router := gin.Default()

	router.Use(requestid.New())
	router.Use(gin.Logger())
	router.Use(logger.New(log))
	router.Use(gin.Recovery())
	router.Use(logger.URLFormat())

	router.POST("/user",userSaver.SaveUserHandler(log, db))
	router.GET("/user", userGetter.GetUserHandler(log, db))
	router.PATCH("/user", userUpdater.UpdateUserPasswordHandler(log, db))
	router.DELETE("/user", userDeleter.DeleteUserHandler(log, db))
	
	router.POST("/task", taskSaver.SaveTaskHandler(log, db))
	router.GET("/task", taskGetter.GetTasksHandler(log, db))
	router.PATCH("/task", taskUpdater.UpdateTaskHandler(log, db))
	router.DELETE("/task", taskDeleter.DeleteTaskHandler(log, db))

	log.Info("starting HTTP server", slog.String("port", cfg.HTTPServer.Address))
	
	server := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.Timeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start HTTP server", sl.Err(err))
	}

	log.Error("Server is shutting down")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envDev, envLocal:
		logger = setupPrettySlog()
	default:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return logger
}

// setupPrettySlog returns a logger that outputs pretty logs
func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}