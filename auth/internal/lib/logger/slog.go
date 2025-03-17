package logger

import (
	"log/slog"
	"os"

	"github.com/ayayaakasvin/auth/internal/lib/logger/handlers/prettyslog"
)

const (
	envProd  = "prod"
	envDev   = "dev"
	envLocal = "local"
)

func SetupLogger(env string) *slog.Logger {
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
