package middleware

import (
	"context"
	"time"

	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func RetryUnaryClientInterceptor(maxRetries int, delay time.Duration, log logger.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if !canRetryGrpcRequest(ctx) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		var lastErr error
		for attempts := 0; attempts <= maxRetries; attempts++ {
			if attempts > 0 {
				log.Info("Retry attempt %d for method %s", attempts, method)
			}

			err := invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}

			lastErr = err
			st, ok := status.FromError(err)
			if !ok {
				return err
			}

			if !isRetryableStatusCode(st.Code()) {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		return lastErr
	}
}

func canRetryGrpcRequest(ctx context.Context) bool {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return false
	}

	values := md.Get("x-idempotency-token")
	return len(values) > 0 && values[0] != ""
}

func isRetryableStatusCode(code codes.Code) bool {
	switch code {
	case codes.Unavailable,
		codes.DeadlineExceeded,
		codes.ResourceExhausted,
		codes.Aborted,
		codes.Internal:
		return true
	}
	return false
}
