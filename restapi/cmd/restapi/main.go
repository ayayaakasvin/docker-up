package main

import (
	"log/slog"
	"os"
	
	"restapi/internal/app"
	"restapi/internal/config"
	"restapi/internal/lib/logger"
	"restapi/internal/lib/sl"
	"restapi/internal/models/postgresql"
)


func main() {
	cfg := config.MustLoadConfig()
	log := logger.SetupLogger(cfg.Env)

	log.Info("starting RESTful API")
	log.Debug("debug message", slog.String("env", cfg.Env))

	db := postgresql.NewPostgreSQL(cfg)
	defer db.Close()

	err := app.App(db, log, cfg)
	if err != nil {
		log.Error("failed to run server", sl.Err(err))
		os.Exit(1)
	}

	os.Exit(0)
}