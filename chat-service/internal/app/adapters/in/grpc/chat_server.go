package grpc

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/SamEkb/messenger-app/chat-service/config/env"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	middlewaregrpc "github.com/SamEkb/messenger-app/chat-service/internal/middleware/grpc"
	chat "github.com/SamEkb/messenger-app/pkg/api/chat_service/v1"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/bufbuild/protovalidate-go"
	protovalidatemw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ChatServer struct {
	chat.UnimplementedChatServiceServer
	validator protovalidate.Validator
	useCase   ports.ChatUseCase
	cfg       *env.ServerConfig
	logger    logger.Logger
}

func NewChatServer(useCase ports.ChatUseCase, cfg *env.ServerConfig, logger logger.Logger) (*ChatServer, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, errors.NewInternalError(err, "failed to initialize validator")
	}

	return &ChatServer{
		validator: validator,
		useCase:   useCase,
		cfg:       cfg,
		logger:    logger,
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

func (s *ChatServer) RunServers(ctx context.Context) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				protovalidatemw.UnaryServerInterceptor(s.validator),
				middlewaregrpc.ErrorsUnaryServerInterceptor(),
			),
		)
		chat.RegisterChatServiceServer(grpcServer, s)
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
		if err := chat.RegisterChatServiceHandlerServer(ctx, mux, s); err != nil {
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
