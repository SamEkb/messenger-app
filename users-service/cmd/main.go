package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
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
}

func NewServer() (*UsersServiceServer, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}
	return &UsersServiceServer{validator: validator}, nil
}

func (s *UsersServiceServer) GetUserProfile(_ context.Context, req *users.GetUserProfileRequest) (*users.GetUserProfileResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	log.Printf("GetUserProfile request: user_id=%s", req.GetUserId())

	return &users.GetUserProfileResponse{
		Nickname:    "user_" + req.GetUserId()[:8],
		Email:       "user" + req.GetUserId()[:4] + "@example.com",
		Description: "This is a test user profile",
		AvatarUrl:   "https://example.com/avatars/default.png",
	}, nil
}

func (s *UsersServiceServer) GetUserProfileByNickname(_ context.Context, req *users.GetUserProfileByNicknameRequest) (*users.GetUserProfileByNicknameResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	log.Printf("GetUserProfileByNickname request: nickname=%s", req.GetNickname())

	userId := uuid.NewString()
	return &users.GetUserProfileByNicknameResponse{
		Nickname:    req.GetNickname(),
		Email:       req.GetNickname() + "@example.com",
		Description: "User profile found by nickname",
		AvatarUrl:   "https://example.com/avatars/" + userId[:8] + ".png",
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

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("received shutdown signal")
		cancel()
	}()

	server, err := NewServer()
	if err != nil {
		log.Fatalf("server init error: %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				protovalidatemw.UnaryServerInterceptor(server.validator),
			),
		)
		users.RegisterUsersServiceServer(grpcServer, server)
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
		if err := users.RegisterUsersServiceHandlerServer(ctx, mux, server); err != nil {
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
}
