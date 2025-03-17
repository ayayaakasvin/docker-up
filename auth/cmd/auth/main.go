package main

import (
	"log/slog"
	"os"

	"github.com/ayayaakasvin/auth/internal/app"
	"github.com/ayayaakasvin/auth/internal/config"
	"github.com/ayayaakasvin/auth/internal/lib/logger"
	"github.com/ayayaakasvin/auth/internal/lib/sl"
	"github.com/ayayaakasvin/auth/internal/models/postgresql"
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.SetupLogger(cfg.Env)
	log.Info("Config", slog.String("env", cfg.Env))

	storage := postgresql.NewPostgresStorage(cfg.Database)
	log.Info("Storage was set up", sl.Any("database", cfg.Database))

	err := app.App(storage, log, cfg)
	if err != nil {
		log.Error("failed to run server", sl.Err(err))
		os.Exit(1)
	}

	os.Exit(0)
}