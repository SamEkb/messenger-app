package middleware

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CircuitBreakerOptions struct {
	Name             string
	MaxRequests      uint32
	Interval         time.Duration
	Timeout          time.Duration
	MinRequests      uint32
	FailureRatio     float64
	ServerErrorCodes map[codes.Code]bool
}

func DefaultCircuitBreakerOptions() *CircuitBreakerOptions {
	return &CircuitBreakerOptions{
		Name:         "grpc_circuit_breaker",
		MaxRequests:  10,
		Interval:     60 * time.Second,
		Timeout:      5 * time.Minute,
		MinRequests:  40,
		FailureRatio: 0.6,
		ServerErrorCodes: map[codes.Code]bool{
			codes.Internal:          true,
			codes.Unavailable:       true,
			codes.DeadlineExceeded:  true,
			codes.ResourceExhausted: true,
			codes.DataLoss:          true,
			codes.Unknown:           true,
			codes.Aborted:           true,
		},
	}
}

type CircuitBreakerOption func(*CircuitBreakerOptions)

func WithMaxRequests(max uint32) CircuitBreakerOption {
	return func(o *CircuitBreakerOptions) {
		o.MaxRequests = max
	}
}

func WithInterval(interval time.Duration) CircuitBreakerOption {
	return func(o *CircuitBreakerOptions) {
		o.Interval = interval
	}
}

func WithTimeout(timeout time.Duration) CircuitBreakerOption {
	return func(o *CircuitBreakerOptions) {
		o.Timeout = timeout
	}
}

func WithMinRequests(min uint32) CircuitBreakerOption {
	return func(o *CircuitBreakerOptions) {
		o.MinRequests = min
	}
}

func WithFailureRatio(ratio float64) CircuitBreakerOption {
	return func(o *CircuitBreakerOptions) {
		o.FailureRatio = ratio
	}
}

func WithServerErrorCodes(errorCodes []string) CircuitBreakerOption {
	return func(o *CircuitBreakerOptions) {
		o.ServerErrorCodes = make(map[codes.Code]bool)
		for _, code := range errorCodes {
			convertedCode := stringToCode(code)
			o.ServerErrorCodes[convertedCode] = true
		}
	}
}

func WithName(name string) CircuitBreakerOption {
	return func(o *CircuitBreakerOptions) {
		o.Name = name
	}
}

func stringToCode(codeStr string) codes.Code {
	codeStr = strings.ToUpper(strings.TrimSpace(codeStr))

	switch codeStr {
	case "OK":
		return codes.OK
	case "CANCELLED", "CANCELED":
		return codes.Canceled
	case "UNKNOWN":
		return codes.Unknown
	case "INVALID_ARGUMENT", "INVALIDARGUMENT":
		return codes.InvalidArgument
	case "DEADLINE_EXCEEDED", "DEADLINEEXCEEDED":
		return codes.DeadlineExceeded
	case "NOT_FOUND", "NOTFOUND":
		return codes.NotFound
	case "ALREADY_EXISTS", "ALREADYEXISTS":
		return codes.AlreadyExists
	case "PERMISSION_DENIED", "PERMISSIONDENIED":
		return codes.PermissionDenied
	case "RESOURCE_EXHAUSTED", "RESOURCEEXHAUSTED":
		return codes.ResourceExhausted
	case "FAILED_PRECONDITION", "FAILEDPRECONDITION":
		return codes.FailedPrecondition
	case "ABORTED":
		return codes.Aborted
	case "OUT_OF_RANGE", "OUTOFRANGE":
		return codes.OutOfRange
	case "UNIMPLEMENTED":
		return codes.Unimplemented
	case "INTERNAL":
		return codes.Internal
	case "UNAVAILABLE":
		return codes.Unavailable
	case "DATA_LOSS", "DATALOSS":
		return codes.DataLoss
	case "UNAUTHENTICATED":
		return codes.Unauthenticated
	default:
		return codes.Unknown
	}
}

func NewCircuitBreakerInterceptor(log logger.Logger, opts ...CircuitBreakerOption) grpc.UnaryClientInterceptor {
	options := DefaultCircuitBreakerOptions()

	for _, opt := range opts {
		opt(options)
	}

	settings := gobreaker.Settings{
		Name:        options.Name,
		MaxRequests: options.MaxRequests,
		Interval:    options.Interval,
		Timeout:     options.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= options.MinRequests && failureRatio >= options.FailureRatio
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Info("Circuit breaker %s state changed from %s to %s", name, from, to)
		},
	}

	cb := gobreaker.NewCircuitBreaker[any](settings)

	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if cb.State() == gobreaker.StateOpen {
			log.Warn("Circuit breaker is open for method %s, fast failing", method)
			return status.Error(codes.Unavailable, "circuit breaker is open")
		}

		result, err := cb.Execute(func() (any, error) {
			callErr := invoker(ctx, method, req, reply, cc, opts...)

			if callErr == nil {
				return reply, nil
			}

			st, ok := status.FromError(callErr)

			if !ok || options.ServerErrorCodes[st.Code()] {
				log.Warn("Server error occurred for method %s: %v", method, callErr)
				return nil, callErr
			}

			return callErr, nil
		})

		if err != nil {
			if errors.Is(err, gobreaker.ErrOpenState) {
				log.Warn("Circuit is open for method %s", method)
				return status.Error(codes.Unavailable, "circuit breaker is open")
			}

			if errors.Is(err, gobreaker.ErrTooManyRequests) {
				log.Warn("Too many requests in half-open state for method %s", method)
				return status.Error(codes.ResourceExhausted, "circuit breaker: too many requests")
			}

			return err
		}

		if errResult, ok := result.(error); ok {
			return errResult
		}

		return nil
	}
}
