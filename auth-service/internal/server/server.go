package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/SamEkb/messenger-app/auth-service/internal/kafka"
	auth "github.com/SamEkb/messenger-app/pkg/api/auth_service/v1"
	"github.com/SamEkb/messenger-app/pkg/api/events/v1"
	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	protovalidatemw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	PortHttp = "8001"
	PortGrpc = "9001"
)

type AuthServiceServer struct {
	auth.UnimplementedAuthServiceServer
	validator protovalidate.Validator
	producer  *kafka.Producer
}

func NewServer() (*AuthServiceServer, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	producer, err := kafka.NewProducer()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka producer: %w", err)
	}

	return &AuthServiceServer{
		validator: validator,
		producer:  producer,
	}, nil
}

func (s *AuthServiceServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	userID := uuid.NewString()
	log.Printf("Register request: username=%s → userID=%s", req.GetUsername(), userID)

	event := &events.UserRegisteredEvent{
		UserId:       userID,
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		RegisteredAt: timestamppb.Now(),
	}

	if err := s.producer.PublishUserRegistered(ctx, event); err != nil {
		log.Printf("Failed to publish user registered event: %v", err)
	}

	return &auth.RegisterResponse{
		UserId:  userID,
		Message: "User registered successfully",
		Success: true,
	}, nil
}

func (s *AuthServiceServer) Login(_ context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	log.Printf("Login request: email=%s", req.GetEmail())
	return &auth.LoginResponse{
		Token:     "jwt-token-example",
		UserId:    uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Success:   true,
		Message:   "Login successful",
	}, nil
}

func (s *AuthServiceServer) Logout(_ context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	log.Println("Logout request")
	return &auth.LogoutResponse{
		Success: true,
		Message: "Logged out successfully",
	}, nil
}

func (s *AuthServiceServer) Close() error {
	if s.producer != nil {
		return s.producer.Close()
	}
	return nil
}

func liveHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}

// RunServers starts both gRPC and HTTP servers
func (s *AuthServiceServer) RunServers(ctx context.Context) error {
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

		addr := ":" + PortGrpc
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
		if err := auth.RegisterAuthServiceHandlerServer(ctx, mux, s); err != nil {
			log.Fatalf("gateway registration error: %v", err)
		}

		root := http.NewServeMux()
		root.Handle("/", mux)
		root.HandleFunc("/live", liveHandler)
		root.HandleFunc("/ready", readyHandler)

		addr := ":" + PortHttp
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
