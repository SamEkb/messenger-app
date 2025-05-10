package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"sync"

	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
	"github.com/SamEkb/messenger-app/users-service/config/env"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
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
	logger      *slog.Logger
}

func NewServer(cfg *env.ServerConfig, userUseCase ports.UserUseCase, logger *slog.Logger) (*UsersServiceServer, error) {
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
		grpcServer := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				protovalidatemw.UnaryServerInterceptor(s.validator),
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
