package grpc

import (
	"context"

	chat "github.com/SamEkb/messenger-app/pkg/api/chat_service/v1"
)

func (s *ChatServer) CreateChat(ctx context.Context, req *chat.CreateChatRequest) (*chat.CreateChatResponse, error) {
	s.logger.Info("creating chat")

	chatDto, err := s.useCase.CreateChat(ctx, req.Participants)
	if err != nil {
		s.logger.Error("failed to create chat", "error", err)
		return nil, err
	}

	s.logger.Info("chat created successfully")

	return &chat.CreateChatResponse{
		ChatId:  chatDto.ID(),
		Success: true,
		Message: "chat created successfully",
	}, nil
}
