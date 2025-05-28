package logger

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

const (
	FieldService     = "service"
	FieldEnvironment = "environment"
	FieldTraceID     = "trace_id"
	FieldSpanID      = "span_id"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	With(args ...any) Logger

	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	FatalContext(ctx context.Context, msg string, args ...any)
	WithContext(ctx context.Context) Logger
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

func (s *slogLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	s.logWithTracing(ctx, slog.LevelDebug, msg, args...)
}

func (s *slogLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	s.logWithTracing(ctx, slog.LevelInfo, msg, args...)
}

func (s *slogLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	s.logWithTracing(ctx, slog.LevelWarn, msg, args...)
}

func (s *slogLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	s.logWithTracing(ctx, slog.LevelError, msg, args...)

	s.logErrorToSpan(ctx, msg, args...)
}

func (s *slogLogger) FatalContext(ctx context.Context, msg string, args ...any) {
	s.logWithTracing(ctx, slog.LevelError, msg, args...)

	s.logErrorToSpan(ctx, msg, args...)

	os.Exit(1)
}

func (s *slogLogger) WithContext(ctx context.Context) Logger {
	traceArgs := s.extractTracingFields(ctx)
	return &slogLogger{l: s.l.With(traceArgs...)}
}

func (s *slogLogger) logWithTracing(ctx context.Context, level slog.Level, msg string, args ...any) {
	tracingArgs := s.extractTracingFields(ctx)

	allArgs := append(tracingArgs, args...)

	switch level {
	case slog.LevelDebug:
		s.l.Debug(msg, allArgs...)
	case slog.LevelInfo:
		s.l.Info(msg, allArgs...)
	case slog.LevelWarn:
		s.l.Warn(msg, allArgs...)
	case slog.LevelError:
		s.l.Error(msg, allArgs...)
	}
}

func (s *slogLogger) extractTracingFields(ctx context.Context) []any {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return []any{}
	}

	spanContext := span.SpanContext()
	return []any{
		FieldTraceID, spanContext.TraceID().String(),
		FieldSpanID, spanContext.SpanID().String(),
	}
}

func (s *slogLogger) logErrorToSpan(ctx context.Context, msg string, args ...any) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	errorMsg := msg
	if len(args) > 0 {
		errorMsg = msg + " " + formatArgs(args...)
	}

	span.AddEvent("error", trace.WithAttributes(
		attribute.String("error.message", errorMsg),
		attribute.String("log.level", "error"),
	))

	span.SetStatus(codes.Error, errorMsg)
}

func formatArgs(args ...any) string {
	if len(args) == 0 {
		return ""
	}

	result := ""
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			if key, ok := args[i].(string); ok {
				if value, ok := args[i+1].(string); ok {
					result += key + "=" + value + " "
				}
			}
		}
	}
	return result
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

func (m *MockLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	m.Debug(msg, args...)
}
func (m *MockLogger) InfoContext(ctx context.Context, msg string, args ...any) { m.Info(msg, args...) }
func (m *MockLogger) WarnContext(ctx context.Context, msg string, args ...any) { m.Warn(msg, args...) }
func (m *MockLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	m.Error(msg, args...)
}
func (m *MockLogger) FatalContext(ctx context.Context, msg string, args ...any) {
	m.Fatal(msg, args...)
}
func (m *MockLogger) WithContext(ctx context.Context) Logger { return m }
