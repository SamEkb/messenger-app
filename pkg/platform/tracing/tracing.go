package tracing

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
)

type Config struct {
	ServiceName    string
	ServiceVersion string
	JaegerURL      string
	Environment    string
	Enabled        bool
}

func LoadConfig() Config {
	return Config{
		ServiceName:    getEnv("OTEL_SERVICE_NAME", "unknown-service"),
		ServiceVersion: getEnv("OTEL_SERVICE_VERSION", "0.1.0"),
		JaegerURL:      getEnv("OTEL_EXPORTER_JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		Environment:    getEnv("OTEL_ENVIRONMENT", "development"),
		Enabled:        getEnv("OTEL_TRACING_ENABLED", "true") == "true",
	}
}

func Initialize(cfg Config) (func(context.Context) error, error) {
	if !cfg.Enabled {
		return func(context.Context) error { return nil }, nil
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerURL)))
	if err != nil {
		return nil, err
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}

func GinMiddleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}

func GRPCServerInterceptor() grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor()
}

func GRPCClientInterceptor() grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor()
}

func HTTPClient() *http.Client {
	return &http.Client{
		Transport: &tracedTransport{
			base: http.DefaultTransport,
		},
	}
}

type tracedTransport struct {
	base http.RoundTripper
}

func (t *tracedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	tracer := otel.Tracer("http-client")

	ctx, span := tracer.Start(req.Context(), "HTTP "+req.Method)
	defer span.End()

	span.SetAttributes(
		semconv.HTTPMethod(req.Method),
		semconv.HTTPURL(req.URL.String()),
	)

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	req = req.WithContext(ctx)
	resp, err := t.base.RoundTrip(req)

	if err != nil {
		span.RecordError(err)
		return resp, err
	}

	span.SetAttributes(semconv.HTTPStatusCode(resp.StatusCode))

	return resp, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
