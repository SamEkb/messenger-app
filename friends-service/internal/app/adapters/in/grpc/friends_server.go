package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/SamEkb/messenger-app/friends-service/config/env"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/ports"
	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/bufbuild/protovalidate-go"
	protovalidatemw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var _ friends.FriendsServiceServer = (*FriendshipServiceServer)(nil)

type FriendshipServiceServer struct {
	friends.UnimplementedFriendsServiceServer
	friendshipUseCase ports.FriendshipUseCase
	validator         protovalidate.Validator
	cfg               *env.ServerConfig
	logger            logger.Logger
}

func NewServer(cfg *env.ServerConfig, friendshipUseCase ports.FriendshipUseCase, logger logger.Logger) (*FriendshipServiceServer, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	return &FriendshipServiceServer{
		validator:         validator,
		cfg:               cfg,
		logger:            logger,
		friendshipUseCase: friendshipUseCase,
	}, nil
}

func liveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}

func (s *FriendshipServiceServer) RunServers(ctx context.Context) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer := grpclib.NewServer(
			grpclib.ChainUnaryInterceptor(
				protovalidatemw.UnaryServerInterceptor(s.validator),
			),
		)
		friends.RegisterFriendsServiceServer(grpcServer, s)
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
		if err := friends.RegisterFriendsServiceHandlerServer(ctx, mux, s); err != nil {
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
