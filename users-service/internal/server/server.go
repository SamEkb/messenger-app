package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/SamEkb/messenger-app/pkg/api/events/v1"
	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
	"github.com/SamEkb/messenger-app/users-service/internal/kafka"
	"github.com/bufbuild/protovalidate-go"
	protovalidatemw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	PortHttp = "8004"
	PortGrpc = "9004"
)

type UsersServiceServer struct {
	users.UnimplementedUsersServiceServer
	validator protovalidate.Validator
	//storage   *storage.Storage
	consumer *kafka.Consumer
}

func NewServer() (*UsersServiceServer, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	server := &UsersServiceServer{
		validator: validator,
	}

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer(server)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka consumer: %w", err)
	}
	server.consumer = consumer

	return server, nil
}

// Start starts the Kafka consumer
func (s *UsersServiceServer) Start(ctx context.Context) error {
	return s.consumer.Start(ctx)
}

// Close closes connections
func (s *UsersServiceServer) Close() error {
	if s.consumer != nil {
		return s.consumer.Close()
	}
	return nil
}

// Consumer returns the Kafka consumer
func (s *UsersServiceServer) Consumer() *kafka.Consumer {
	return s.consumer
}

// HandleUserRegistered implements the Kafka EventHandler for user registration events
func (s *UsersServiceServer) HandleUserRegistered(_ context.Context, event *events.UserRegisteredEvent) error {
	log.Printf("Handling UserRegisteredEvent for user %s (%s)", event.GetUsername(), event.GetUserId())
	return nil
}

func (s *UsersServiceServer) GetUserProfile(_ context.Context, req *users.GetUserProfileRequest) (*users.GetUserProfileResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	log.Printf("GetUserProfile request: user_id=%s", req.GetUserId())

	return &users.GetUserProfileResponse{
		Nickname:    "test-nickname",
		Email:       "test-email",
		Description: "test-description",
		AvatarUrl:   "test-avatar-url",
	}, nil
}

func (s *UsersServiceServer) GetUserProfileByNickname(_ context.Context, req *users.GetUserProfileByNicknameRequest) (*users.GetUserProfileByNicknameResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	log.Printf("GetUserProfileByNickname request: nickname=%s", req.GetNickname())

	return &users.GetUserProfileByNicknameResponse{
		Nickname:    "test-nickname",
		Email:       "test-email",
		Description: "test-description",
		AvatarUrl:   "test-avatar-url",
	}, nil
}

func (s *UsersServiceServer) UpdateUserProfile(_ context.Context, req *users.UpdateUserProfileRequest) (*users.UpdateUserProfileResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	log.Printf("UpdateUserProfile request: user_id=%s, nickname=%s", req.GetUserId(), req.GetNickname())

	return &users.UpdateUserProfileResponse{
		Success: true,
		Message: "User profile updated successfully",
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

// RunServers starts both gRPC and HTTP servers
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
		if err := users.RegisterUsersServiceHandlerServer(ctx, mux, s); err != nil {
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
