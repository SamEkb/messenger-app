package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/SamEkb/messenger-app/pkg/platform/middleware/metrics"
	"github.com/SamEkb/messenger-app/pkg/platform/middleware/resilience"
	"github.com/SamEkb/messenger-app/pkg/platform/middleware/tracing"

	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/SamEkb/messenger-app/users-service/config/env"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
	middlewaregrpc "github.com/SamEkb/messenger-app/users-service/internal/middleware/grpc"
	"github.com/bufbuild/protovalidate-go"
	protovalidatemw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var _ users.UsersServiceServer = (*UsersServiceServer)(nil)

type UsersServiceServer struct {
	users.UnimplementedUsersServiceServer
	userUseCase ports.UserUseCase
	validator   protovalidate.Validator
	cfg         *env.ServerConfig
	logger      logger.Logger
}

func NewServer(cfg *env.ServerConfig, userUseCase ports.UserUseCase, logger logger.Logger) (*UsersServiceServer, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	server := &UsersServiceServer{
		validator:   validator,
		cfg:         cfg,
		logger:      logger,
		userUseCase: userUseCase,
	}

	return server, nil
}

func liveHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}

func (s *UsersServiceServer) RunServers(ctx context.Context) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		mux := http.NewServeMux()
		mux.Handle("/metrics", metrics.Handler())
		mux.Handle("/debug/pprof/", metrics.PprofHandler())

		metricsServer := &http.Server{
			Addr:    ":9090",
			Handler: mux,
		}

		s.logger.InfoContext(ctx, "metrics and pprof server started", "address", ":9090")

		go func() {
			if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.logger.ErrorContext(ctx, "metrics server error", "error", err)
			}
		}()

		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := metricsServer.Shutdown(shutdownCtx); err != nil {
			s.logger.ErrorContext(ctx, "metrics server shutdown error", "error", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		recoverer := resilience.RecoveryInterceptor(s.logger)
		rls := resilience.NewServerInterceptor(s.logger, s.cfg.RateLimiter.DefaultLimit, s.cfg.RateLimiter.DefaultBurst)
		if s.cfg.RateLimiter.GlobalLimit > 0 && s.cfg.RateLimiter.GlobalBurst > 0 {
			rls = rls.WithGlobalLimit(s.cfg.RateLimiter.GlobalLimit, s.cfg.RateLimiter.GlobalBurst)
		}
		for method, lim := range s.cfg.RateLimiter.MethodLimits {
			rls = rls.WithMethodLimit(method, lim.Limit, lim.Burst)
		}

		grpcServer := grpc.NewServer(
			grpc.StatsHandler(tracing.GRPCServerHandler()),
			grpc.ChainUnaryInterceptor(
				recoverer,
				metrics.GRPCMetricsInterceptor("users-service"),
				rls.Interceptor(),
				protovalidatemw.UnaryServerInterceptor(s.validator),
				middlewaregrpc.ErrorsUnaryServerInterceptor(),
			),
		)
		users.RegisterUsersServiceServer(grpcServer, s)
		reflection.Register(grpcServer)

		addr := s.cfg.GrpcAddr()
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("gRPC listen error: %v", err)
		}
		log.Printf("gRPC listening on %s", addr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC serve error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		mux := runtime.NewServeMux()
		if err := users.RegisterUsersServiceHandlerServer(ctx, mux, s); err != nil {
			log.Fatalf("gateway registration error: %v", err)
		}

		root := http.NewServeMux()
		root.Handle("/", mux)
		root.HandleFunc("/live", liveHandler)
		root.HandleFunc("/ready", readyHandler)

		addr := s.cfg.HttpAddr()
		httpServer := &http.Server{
			Addr:    addr,
			Handler: root,
		}
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("HTTP listen error: %v", err)
		}
		log.Printf("HTTP listening on %s", addr)
		if err := httpServer.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP serve error: %v", err)
		}
	}()

	wg.Wait()
	return nil
}
