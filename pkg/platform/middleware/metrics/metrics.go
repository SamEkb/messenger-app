package metrics

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status_code", "service"},
	)

	grpcRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status_code", "service"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "service"},
	)

	grpcRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "gRPC request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "service"},
	)

	httpErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors (4xx, 5xx)",
		},
		[]string{"method", "path", "status_code", "service"},
	)

	grpcErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_errors_total",
			Help: "Total number of gRPC errors",
		},
		[]string{"method", "status_code", "service"},
	)

	httpActiveRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Current number of active HTTP requests",
		},
		[]string{"service"},
	)

	grpcActiveRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "grpc_active_requests",
			Help: "Current number of active gRPC requests",
		},
		[]string{"service"},
	)
)

func HTTPMetricsMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		httpActiveRequests.WithLabelValues(serviceName).Inc()
		defer httpActiveRequests.WithLabelValues(serviceName).Dec()

		c.Next()

		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			statusCode,
			serviceName,
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			serviceName,
		).Observe(duration)

		if c.Writer.Status() >= 400 {
			httpErrorsTotal.WithLabelValues(
				c.Request.Method,
				c.FullPath(),
				statusCode,
				serviceName,
			).Inc()
		}
	}
}

func GRPCMetricsInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		grpcActiveRequests.WithLabelValues(serviceName).Inc()
		defer grpcActiveRequests.WithLabelValues(serviceName).Dec()

		resp, err := handler(ctx, req)

		duration := time.Since(start).Seconds()

		statusCode := codes.OK.String()
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code().String()
			}
		}

		grpcRequestsTotal.WithLabelValues(
			info.FullMethod,
			statusCode,
			serviceName,
		).Inc()

		grpcRequestDuration.WithLabelValues(
			info.FullMethod,
			serviceName,
		).Observe(duration)

		if err != nil {
			grpcErrorsTotal.WithLabelValues(
				info.FullMethod,
				statusCode,
				serviceName,
			).Inc()
		}

		return resp, err
	}
}

func Handler() http.Handler {
	return promhttp.Handler()
}
