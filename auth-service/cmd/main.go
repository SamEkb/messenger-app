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
	"time"

	auth "github.com/SamEkb/messenger-app/pkg/api/auth_service/v1"
	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	protovalidatemw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	PortHttp = "8001"
	PortGrpc = "9001"
)

type AuthServiceServer struct {
	auth.UnimplementedAuthServiceServer
	validator protovalidate.Validator
}

func NewServer() (*AuthServiceServer, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}
	return &AuthServiceServer{validator: validator}, nil
}

func (s *AuthServiceServer) Register(_ context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	userID := uuid.NewString()
	log.Printf("Register request: username=%s → userID=%s", req.GetUsername(), userID)
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
		auth.RegisterAuthServiceServer(grpcServer, server)
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
		if err := auth.RegisterAuthServiceHandlerServer(ctx, mux, server); err != nil {
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
