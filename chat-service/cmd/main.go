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

	"net/http"

	chat "github.com/SamEkb/messenger-app/pkg/api/chat_service/v1"
	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	protovalidatemw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	PortHttp = "8002"
	PortGrpc = "9002"
)

type ChatServiceServer struct {
	chat.UnimplementedChatServiceServer
	validator protovalidate.Validator
}

func NewServer() (*ChatServiceServer, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}
	return &ChatServiceServer{validator: validator}, nil
}

func (c ChatServiceServer) CreateChat(_ context.Context, req *chat.CreateChatRequest) (*chat.CreateChatResponse, error) {
	if err := c.validator.Validate(req); err != nil {
		return nil, err
	}
	return &chat.CreateChatResponse{
		ChatId:  "chat-id-example",
		Success: true,
	}, nil
}

func (c *ChatServiceServer) GetUserChats(_ context.Context, req *chat.GetUserChatsRequest) (*chat.GetUserChatsResponse, error) {
	if err := c.validator.Validate(req); err != nil {
		return nil, err
	}

	chats := make([]*chat.Chat, 2)
	chats[0] = &chat.Chat{}
	chats[1] = &chat.Chat{}
	return &chat.GetUserChatsResponse{Chats: chats}, nil
}

func (c *ChatServiceServer) SendMessage(_ context.Context, req *chat.SendMessageRequest) (*chat.SendMessageResponse, error) {
	if err := c.validator.Validate(req); err != nil {
		return nil, err
	}

	msg := &chat.Message{MessageId: uuid.NewString()}
	return &chat.SendMessageResponse{Message: msg, Success: true, MessageInfo: "success"}, nil
}

func (c *ChatServiceServer) GetChatHistory(_ context.Context, req *chat.GetChatHistoryRequest) (*chat.GetChatHistoryResponse, error) {
	if err := c.validator.Validate(req); err != nil {
		return nil, err
	}

	return &chat.GetChatHistoryResponse{
		Messages:      nil,
		TotalMessages: 100,
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
		chat.RegisterChatServiceServer(grpcServer, server)
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
		if err := chat.RegisterChatServiceHandlerServer(ctx, mux, server); err != nil {
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
