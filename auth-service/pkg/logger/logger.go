package logger

import (
	"log/slog"
	"os"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func NewLogger(env string, serviceName string) *slog.Logger {
	var logger *slog.Logger

	opts := &slog.HandlerOptions{
		AddSource: true,
	}

	switch env {
	case EnvLocal:
		opts.Level = slog.LevelDebug
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	case EnvDev:
		opts.Level = slog.LevelDebug
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	case EnvProd:
		opts.Level = slog.LevelInfo
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	default:
		opts.Level = slog.LevelInfo
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	return logger.With("service", serviceName)
}
