package main

import (
	"log/slog"
	"net/http"
	"os"

	"restapi/internal/config"
	taskDeleter "restapi/internal/http-server/handlers/task/delete"
	taskGetter "restapi/internal/http-server/handlers/task/get/taskid"
	tasksGetter "restapi/internal/http-server/handlers/task/get/userid"
	taskSaver "restapi/internal/http-server/handlers/task/save"
	taskUpdater "restapi/internal/http-server/handlers/task/update"
	userDeleter "restapi/internal/http-server/handlers/user/delete"
	userGetter "restapi/internal/http-server/handlers/user/get"
	userSaver "restapi/internal/http-server/handlers/user/save"
	userUpdater "restapi/internal/http-server/handlers/user/update"
	"restapi/internal/http-server/middleware"
	logger "restapi/internal/http-server/middleware/logger"
	slogpretty "restapi/internal/lib/logger/handlers/prettyslog"
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

	middleware.MiddlewareAdd(router, requestid.New(), gin.Logger(), logger.New(log), gin.Recovery(), logger.URLFormat(), middleware.CorsDefault)

	userRouter := router.Group("/user")
	{
		userRouter.POST("/", userSaver.SaveUserHandler(log, db))
		userRouter.GET("/:userId", userGetter.GetUserHandler(log, db))
		userRouter.PUT("/:userId/password", userUpdater.UpdateUserPasswordHandler(log, db))
		userRouter.DELETE("/:userId", userDeleter.DeleteUserHandler(log, db))
	}

	taskRouter := router.Group("/task")
	{
		taskRouter.POST("/:userId", taskSaver.SaveTaskHandler(log, db))
		taskRouter.GET("/user/:userId", tasksGetter.GetUserTasksHandler(log, db))
		taskRouter.GET("/:taskId", taskGetter.GetTaskHandler(log, db))
		taskRouter.PUT("/:taskId", taskUpdater.UpdateTaskHandler(log, db))
		taskRouter.DELETE("/:taskId", taskDeleter.DeleteTaskHandler(log, db))
	}

	log.Info("starting HTTP server", slog.String("port", cfg.HTTPServer.Address))

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Timeout,
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
