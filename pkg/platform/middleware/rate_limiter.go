package middleware

import (
	"context"
	"sync"

	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RateLimitConfig struct {
	Limit rate.Limit
	Burst int
}

func NewClientInterceptor(log logger.Logger, limit rate.Limit, burst int) grpc.UnaryClientInterceptor {
	limiter := rate.NewLimiter(limit, burst)

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if err := limiter.Wait(ctx); err != nil {
			log.Warn("Client rate limit exceeded for method %s: %v", method, err)
			return status.Error(codes.ResourceExhausted, "client-side rate limit exceeded")
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func NewServerInterceptor(log logger.Logger, limit float64, burst int) *ServerLimiter {
	rateLimit := convertToRateLimit(limit)
	return &ServerLimiter{
		logger:       log,
		defaultLimit: rateLimit,
		defaultBurst: burst,
		methodLimits: make(map[string]*rate.Limiter),
		mutex:        &sync.RWMutex{},
	}
}

func convertToRateLimit(limit float64) rate.Limit {
	return rate.Limit(limit)
}

type ServerLimiter struct {
	logger        logger.Logger
	defaultLimit  rate.Limit
	defaultBurst  int
	methodLimits  map[string]*rate.Limiter
	globalLimiter *rate.Limiter
	mutex         *sync.RWMutex
}

func (s *ServerLimiter) WithMethodLimit(fullMethodName string, limit float64, burst int) *ServerLimiter {
	rateLimit := convertToRateLimit(limit)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.methodLimits[fullMethodName] = rate.NewLimiter(rateLimit, burst)
	return s
}

func (s *ServerLimiter) WithGlobalLimit(limit float64, burst int) *ServerLimiter {
	rateLimit := convertToRateLimit(limit)
	s.globalLimiter = rate.NewLimiter(rateLimit, burst)
	return s
}

func (s *ServerLimiter) Interceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if s.globalLimiter != nil && !s.globalLimiter.Allow() {
			s.logger.Warn("Global rate limit exceeded for method %s", info.FullMethod)
			return nil, status.Error(codes.ResourceExhausted, "global rate limit exceeded")
		}

		s.mutex.RLock()
		limiter, ok := s.methodLimits[info.FullMethod]
		s.mutex.RUnlock()

		if !ok {
			limiter = rate.NewLimiter(s.defaultLimit, s.defaultBurst)

			s.mutex.Lock()
			s.methodLimits[info.FullMethod] = limiter
			s.mutex.Unlock()
		}

		if !limiter.Allow() {
			s.logger.Warn("Method rate limit exceeded for %s", info.FullMethod)
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(ctx, req)
	}
}
