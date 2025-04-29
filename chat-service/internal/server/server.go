package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/SamEkb/messenger-app/chat-service/internal/clients"
	chat "github.com/SamEkb/messenger-app/pkg/api/chat_service/v1"
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
	PortHttp = "8002"
	PortGrpc = "9002"
)

// Simple in-memory chat storage
// In a real application, this would be a MongoDB database
type ChatStorage struct {
	mu       sync.RWMutex
	chats    map[string]*chat.Chat
	messages map[string][]*chat.Message
}

func NewChatStorage() *ChatStorage {
	return &ChatStorage{
		chats:    make(map[string]*chat.Chat),
		messages: make(map[string][]*chat.Message),
	}
}

func (s *ChatStorage) CreateChat(chatID string, participants []string) *chat.Chat {
	s.mu.Lock()
	defer s.mu.Unlock()

	newChat := &chat.Chat{
		ChatId:       chatID,
		Participants: participants,
	}

	s.chats[chatID] = newChat
	s.messages[chatID] = []*chat.Message{}

	return newChat
}

func (s *ChatStorage) GetChat(chatID string) (*chat.Chat, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chat, exists := s.chats[chatID]
	return chat, exists
}

func (s *ChatStorage) GetUserChats(userID string) []*chat.Chat {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userChats []*chat.Chat

	for _, c := range s.chats {
		for _, participant := range c.GetParticipants() {
			if participant == userID {
				// If there are messages, set the last message
				if msgs, exists := s.messages[c.GetChatId()]; exists && len(msgs) > 0 {
					// Create a copy with the last message
					chatWithLastMsg := &chat.Chat{
						ChatId:       c.GetChatId(),
						Participants: c.GetParticipants(),
						LastMessage:  msgs[len(msgs)-1],
					}
					userChats = append(userChats, chatWithLastMsg)
				} else {
					userChats = append(userChats, c)
				}
				break
			}
		}
	}

	return userChats
}

func (s *ChatStorage) AddMessage(chatID string, msg *chat.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.messages[chatID]; !exists {
		s.messages[chatID] = []*chat.Message{}
	}

	s.messages[chatID] = append(s.messages[chatID], msg)

	// Update the last message in the chat
	if chat, exists := s.chats[chatID]; exists {
		chat.LastMessage = msg
	}
}

func (s *ChatStorage) GetMessages(chatID string, limit, offset int32) ([]*chat.Message, int32) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msgs, exists := s.messages[chatID]
	if !exists {
		return []*chat.Message{}, 0
	}

	totalMessages := int32(len(msgs))

	// Handle pagination
	if limit <= 0 {
		limit = 20 // Default limit
	}

	if offset < 0 {
		offset = 0
	}

	// Calculate start and end indices
	start := int(offset)
	end := int(offset + limit)

	if start >= len(msgs) {
		return []*chat.Message{}, totalMessages
	}

	if end > len(msgs) {
		end = len(msgs)
	}

	return msgs[start:end], totalMessages
}

type ChatServiceServer struct {
	chat.UnimplementedChatServiceServer
	validator     protovalidate.Validator
	storage       *ChatStorage
	usersClient   *clients.UsersClient
	friendsClient *clients.FriendsClient
}

func NewServer() (*ChatServiceServer, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	usersClient, err := clients.NewUsersClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Users Service client: %w", err)
	}

	friendsClient, err := clients.NewFriendsClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Friends Service client: %w", err)
	}

	return &ChatServiceServer{
		validator:     validator,
		storage:       NewChatStorage(),
		usersClient:   usersClient,
		friendsClient: friendsClient,
	}, nil
}

func (s *ChatServiceServer) Close() error {
	var err1, err2 error

	if s.usersClient != nil {
		err1 = s.usersClient.Close()
	}

	if s.friendsClient != nil {
		err2 = s.friendsClient.Close()
	}

	if err1 != nil {
		return err1
	}

	return err2
}

func (s *ChatServiceServer) CreateChat(ctx context.Context, req *chat.CreateChatRequest) (*chat.CreateChatResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	participants := req.GetParticipants()
	if len(participants) < 2 {
		return nil, status.Error(codes.InvalidArgument, "chat must have at least 2 participants")
	}

	log.Printf("CreateChat request with participants: %v", participants)

	for _, userID := range participants {
		_, err := s.usersClient.GetUserProfile(ctx, userID)
		if err != nil {
			log.Printf("Error verifying user %s: %v", userID, err)
			return nil, status.Errorf(codes.NotFound, "user %s not found", userID)
		}
	}

	for i := 0; i < len(participants)-1; i++ {
		for j := i + 1; j < len(participants); j++ {
			resp, err := s.friendsClient.CheckFriendshipStatus(ctx, participants[i], participants[j])
			if err != nil {
				log.Printf("Error checking friendship status: %v", err)
				return nil, status.Errorf(codes.Internal, "failed to check friendship status: %v", err)
			}

			if resp.GetStatus() != friends.FriendshipStatus_FRIENDSHIP_STATUS_ACCEPTED {
				return nil, status.Errorf(codes.PermissionDenied,
					"users %s and %s are not friends", participants[i], participants[j])
			}
		}
	}

	chatID := uuid.New().String()
	s.storage.CreateChat(chatID, participants)

	return &chat.CreateChatResponse{
		ChatId:  chatID,
		Success: true,
		Message: "Chat created successfully",
	}, nil
}

func (s *ChatServiceServer) GetUserChats(ctx context.Context, req *chat.GetUserChatsRequest) (*chat.GetUserChatsResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	userID := req.GetUserId()
	log.Printf("GetUserChats request for user %s", userID)

	_, err := s.usersClient.GetUserProfile(ctx, userID)
	if err != nil {
		log.Printf("Error verifying user %s: %v", userID, err)
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	userChats := s.storage.GetUserChats(userID)

	return &chat.GetUserChatsResponse{
		Chats: userChats,
	}, nil
}

func (s *ChatServiceServer) SendMessage(ctx context.Context, req *chat.SendMessageRequest) (*chat.SendMessageResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	chatID := req.GetChatId()
	authorID := req.GetAuthorId()
	content := req.GetContent()

	log.Printf("SendMessage request from user %s to chat %s", authorID, chatID)

	// Verify chat exists
	chatObj, exists := s.storage.GetChat(chatID)
	if !exists {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	// Verify author is a participant
	isParticipant := false
	for _, participant := range chatObj.GetParticipants() {
		if participant == authorID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, status.Errorf(codes.PermissionDenied, "user is not a participant in this chat")
	}

	// Create message
	msg := &chat.Message{
		MessageId: uuid.New().String(),
		ChatId:    chatID,
		AuthorId:  authorID,
		Content:   content,
		Timestamp: timestamppb.Now(),
	}

	// Store message
	s.storage.AddMessage(chatID, msg)

	return &chat.SendMessageResponse{
		Message:     msg,
		Success:     true,
		MessageInfo: "Message sent successfully",
	}, nil
}

func (s *ChatServiceServer) GetChatHistory(ctx context.Context, req *chat.GetChatHistoryRequest) (*chat.GetChatHistoryResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	chatID := req.GetChatId()
	limit := req.GetLimit()
	offset := req.GetOffset()

	log.Printf("GetChatHistory request for chat %s (limit: %d, offset: %d)", chatID, limit, offset)

	// Verify chat exists
	_, exists := s.storage.GetChat(chatID)
	if !exists {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	messages, totalCount := s.storage.GetMessages(chatID, limit, offset)

	return &chat.GetChatHistoryResponse{
		Messages:      messages,
		TotalMessages: totalCount,
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
func (s *ChatServiceServer) RunServers(ctx context.Context) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				protovalidatemw.UnaryServerInterceptor(s.validator),
			),
		)
		chat.RegisterChatServiceServer(grpcServer, s)
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
		if err := chat.RegisterChatServiceHandlerServer(ctx, mux, s); err != nil {
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
