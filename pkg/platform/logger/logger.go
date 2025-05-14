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

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	With(args ...any) Logger
}

type slogLogger struct {
	l *slog.Logger
}

func (s *slogLogger) Debug(msg string, args ...any) {
	s.l.Debug(msg, args...)
}
func (s *slogLogger) Info(msg string, args ...any) {
	s.l.Info(msg, args...)
}
func (s *slogLogger) Warn(msg string, args ...any) {
	s.l.Warn(msg, args...)
}
func (s *slogLogger) Error(msg string, args ...any) {
	s.l.Error(msg, args...)
}
func (s *slogLogger) Fatal(msg string, args ...any) {
	s.l.Error(msg, args...)
	os.Exit(1)
}
func (s *slogLogger) With(args ...any) Logger {
	return &slogLogger{l: s.l.With(args...)}
}

func NewLogger(env string, serviceName string) Logger {
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

	return &slogLogger{l: logger.With("service", serviceName)}
}
