package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"buf.build/go/protovalidate"
	"github.com/SamEkb/messenger-app/auth-service/config/env"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	auth "github.com/SamEkb/messenger-app/pkg/api/auth_service/v1"
	protovalidatemw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	auth.UnimplementedAuthServiceServer
	validator   protovalidate.Validator
	authUseCase ports.AuthUseCase
	cfg         *env.ServerConfig
	logger      *slog.Logger
}

func NewServer(cfg *env.ServerConfig, authUseCase ports.AuthUseCase, logger *slog.Logger) (*Server, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	return &Server{
		validator:   validator,
		authUseCase: authUseCase,
		cfg:         cfg,
		logger:      logger,
	}, nil
}

func liveHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}

func (s *Server) RunServers(ctx context.Context) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				protovalidatemw.UnaryServerInterceptor(s.validator),
			),
		)
		auth.RegisterAuthServiceServer(grpcServer, s)
		reflection.Register(grpcServer)

		addr := ":" + strconv.Itoa(s.cfg.GRPCPort)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			s.logger.Error("gRPC listen error", "error", err)
			return
		}
		s.logger.Info("gRPC server started", "address", addr)

		grpcErrCh := make(chan error, 1)
		go func() {
			grpcErrCh <- grpcServer.Serve(lis)
		}()

		select {
		case <-ctx.Done():
			s.logger.Info("gRPC server shutdown initiated")
			grpcServer.GracefulStop()
			return
		case err := <-grpcErrCh:
			if err != nil {
				s.logger.Error("gRPC server error", "error", err)
			}
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		mux := runtime.NewServeMux()
		if err := auth.RegisterAuthServiceHandlerServer(ctx, mux, s); err != nil {
			s.logger.Error("gateway registration error", "error", err)
			return
		}

		root := http.NewServeMux()
		root.Handle("/", mux)
		root.HandleFunc("/live", liveHandler)
		root.HandleFunc("/ready", readyHandler)

		addr := ":" + strconv.Itoa(s.cfg.HTTPPort)
		httpServer := &http.Server{
			Addr:    addr,
			Handler: root,
		}
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			s.logger.Error("HTTP listen error", "error", err)
			return
		}
		s.logger.Info("HTTP server started", "address", addr)

		httpErrCh := make(chan error, 1)
		go func() {
			httpErrCh <- httpServer.Serve(lis)
		}()

		select {
		case <-ctx.Done():
			s.logger.Info("HTTP server shutdown initiated")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := httpServer.Shutdown(shutdownCtx); err != nil {
				s.logger.Error("HTTP server shutdown error", "error", err)
			}
			return
		case err := <-httpErrCh:
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.logger.Error("HTTP server error", "error", err)
			}
			return
		}
	}()

	s.logger.Info("all servers started")

	<-ctx.Done()
	s.logger.Info("server shutdown completed")

	wg.Wait()
	return nil
}
