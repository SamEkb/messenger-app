package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"net/http"

	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
	"github.com/bufbuild/protovalidate-go"
	protovalidatemw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	PortHttp = "8003"
	PortGrpc = "9003"
)

type FriendsServiceServer struct {
	friends.UnimplementedFriendsServiceServer
	validator protovalidate.Validator
}

func NewServer() (*FriendsServiceServer, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}
	return &FriendsServiceServer{validator: validator}, nil
}

func (f *FriendsServiceServer) GetFriendsList(_ context.Context, req *friends.GetFriendsListRequest) (*friends.GetFriendsListResponse, error) {
	if err := f.validator.Validate(req); err != nil {
		return nil, err
	}

	return &friends.GetFriendsListResponse{
		Friends: make([]*friends.FriendInfo, 1),
	}, nil
}
func (f *FriendsServiceServer) SendFriendRequest(_ context.Context, req *friends.SendFriendRequestRequest) (*friends.SendFriendRequestResponse, error) {
	if err := f.validator.Validate(req); err != nil {
		return nil, err
	}

	return &friends.SendFriendRequestResponse{
		Message: "Friend request sent successfully",
		Success: true,
	}, nil
}
func (f *FriendsServiceServer) AcceptFriendRequest(_ context.Context, req *friends.AcceptFriendRequestRequest) (*friends.AcceptFriendRequestResponse, error) {
	if err := f.validator.Validate(req); err != nil {
		return nil, err
	}

	return &friends.AcceptFriendRequestResponse{
		Message: "Friend request accepted successfully",
		Success: true,
	}, nil
}
func (f *FriendsServiceServer) RejectFriendRequest(_ context.Context, req *friends.RejectFriendRequestRequest) (*friends.RejectFriendRequestResponse, error) {
	if err := f.validator.Validate(req); err != nil {
		return nil, err
	}

	return &friends.RejectFriendRequestResponse{
		Message: "Friend request rejected successfully",
		Success: true,
	}, nil
}
func (f *FriendsServiceServer) RemoveFriend(_ context.Context, req *friends.RemoveFriendRequest) (*friends.RemoveFriendResponse, error) {
	if err := f.validator.Validate(req); err != nil {
		return nil, err
	}

	return &friends.RemoveFriendResponse{
		Message: "Friend removed successfully",
		Success: true,
	}, nil
}
func (f *FriendsServiceServer) CheckFriendshipStatus(_ context.Context, req *friends.CheckFriendshipStatusRequest) (*friends.CheckFriendshipStatusResponse, error) {
	if err := f.validator.Validate(req); err != nil {
		return nil, err
	}

	return &friends.CheckFriendshipStatusResponse{
		Status:    friends.FriendshipStatus_FRIENDSHIP_STATUS_ACCEPTED,
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}, nil
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
		friends.RegisterFriendsServiceServer(grpcServer, server)
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
		if err := friends.RegisterFriendsServiceHandlerServer(ctx, mux, server); err != nil {
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

func liveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}
