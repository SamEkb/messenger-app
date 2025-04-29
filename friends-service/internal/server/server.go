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

	"github.com/SamEkb/messenger-app/friends-service/internal/clients"
	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	protovalidatemw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	PortHttp = "8003"
	PortGrpc = "9003"
)

type FriendsServiceServer struct {
	friends.UnimplementedFriendsServiceServer
	validator   protovalidate.Validator
	usersClient *clients.UsersClient
}

func NewServer() (*FriendsServiceServer, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	usersClient, err := clients.NewUsersClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Users Service client: %w", err)
	}

	return &FriendsServiceServer{
		validator:   validator,
		usersClient: usersClient,
	}, nil
}

func (s *FriendsServiceServer) Close() error {
	if s.usersClient != nil {
		return s.usersClient.Close()
	}
	return nil
}

func (s *FriendsServiceServer) getUserIDByNickname(ctx context.Context, nickname string) (string, error) {
	resp, err := s.usersClient.GetUserProfileByNickname(ctx, nickname)
	if err != nil {
		return "", fmt.Errorf("failed to get user by nickname: %w", err)
	}

	return resp.GetNickname(), nil
}

func (s *FriendsServiceServer) GetFriendsList(ctx context.Context, req *friends.GetFriendsListRequest) (*friends.GetFriendsListResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	userID := req.GetUserId()
	log.Printf("GetFriendsList request for user %s", userID)

	friendResp, err := s.usersClient.GetUserProfile(ctx, uuid.New().String())
	if err != nil {
		log.Printf("Error getting user profile: %v", err)
		return &friends.GetFriendsListResponse{
			Friends: []*friends.FriendInfo{},
		}, nil
	}

	friendInfo := &friends.FriendInfo{
		UserId:    uuid.New().String(),
		Nickname:  friendResp.GetNickname(),
		AvatarUrl: friendResp.GetAvatarUrl(),
		Status:    friends.FriendshipStatus_FRIENDSHIP_STATUS_ACCEPTED,
		CreatedAt: timestamppb.New(time.Now().Add(-24 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	return &friends.GetFriendsListResponse{
		Friends: []*friends.FriendInfo{friendInfo},
	}, nil
}

func (s *FriendsServiceServer) SendFriendRequest(ctx context.Context, req *friends.SendFriendRequestRequest) (*friends.SendFriendRequestResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	userID := req.GetUserId()
	friendNickname := req.GetFriendNickname()

	log.Printf("SendFriendRequest from user %s to %s", userID, friendNickname)

	return &friends.SendFriendRequestResponse{
		Message: "Friend request sent successfully",
		Success: true,
	}, nil
}

func (s *FriendsServiceServer) AcceptFriendRequest(ctx context.Context, req *friends.AcceptFriendRequestRequest) (*friends.AcceptFriendRequestResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	return &friends.AcceptFriendRequestResponse{
		Message: "Friend request accepted successfully",
		Success: true,
	}, nil
}

func (s *FriendsServiceServer) RejectFriendRequest(ctx context.Context, req *friends.RejectFriendRequestRequest) (*friends.RejectFriendRequestResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	userID := req.GetUserId()
	friendNickname := req.GetFriendNickname()

	friendID, err := s.getUserIDByNickname(ctx, friendNickname)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "friend not found: %v", err)
	}

	log.Printf("RejectFriendRequest from user %s to %s", userID, friendID)

	return &friends.RejectFriendRequestResponse{
		Message: "Friend request rejected successfully",
		Success: true,
	}, nil
}

func (s *FriendsServiceServer) RemoveFriend(ctx context.Context, req *friends.RemoveFriendRequest) (*friends.RemoveFriendResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	userID := req.GetUserId()
	friendNickname := req.GetFriendNickname()

	friendID, err := s.getUserIDByNickname(ctx, friendNickname)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "friend not found: %v", err)
	}

	log.Printf("RemoveFriend from user %s to %s", userID, friendID)

	return &friends.RemoveFriendResponse{
		Message: "Friend removed successfully",
		Success: true,
	}, nil
}

func (s *FriendsServiceServer) CheckFriendshipStatus(ctx context.Context, req *friends.CheckFriendshipStatusRequest) (*friends.CheckFriendshipStatusResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	now := time.Now()
	return &friends.CheckFriendshipStatusResponse{
		Status:    friends.FriendshipStatus_FRIENDSHIP_STATUS_UNSPECIFIED,
		CreatedAt: timestamppb.New(now.Add(-24 * time.Hour)),
		UpdatedAt: timestamppb.New(now),
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

// RunServers starts both gRPC and HTTP servers
func (s *FriendsServiceServer) RunServers(ctx context.Context) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				protovalidatemw.UnaryServerInterceptor(s.validator),
			),
		)
		friends.RegisterFriendsServiceServer(grpcServer, s)
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
		if err := friends.RegisterFriendsServiceHandlerServer(ctx, mux, s); err != nil {
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
