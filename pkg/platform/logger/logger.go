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

const (
	FieldService     = "service"
	FieldEnvironment = "environment"
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

	return &slogLogger{
		l: logger.With(
			FieldService, serviceName,
			FieldEnvironment, env,
		),
	}
}

type MockLogger struct {
	Entries []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{Entries: make([]string, 0)}
}

func (m *MockLogger) Debug(msg string, args ...any) { m.Entries = append(m.Entries, "DEBUG: "+msg) }
func (m *MockLogger) Info(msg string, args ...any)  { m.Entries = append(m.Entries, "INFO: "+msg) }
func (m *MockLogger) Warn(msg string, args ...any)  { m.Entries = append(m.Entries, "WARN: "+msg) }
func (m *MockLogger) Error(msg string, args ...any) { m.Entries = append(m.Entries, "ERROR: "+msg) }
func (m *MockLogger) With(args ...any) Logger       { return m }
func (m *MockLogger) Fatal(msg string, args ...any) { m.Entries = append(m.Entries, "FATAL: "+msg) }
