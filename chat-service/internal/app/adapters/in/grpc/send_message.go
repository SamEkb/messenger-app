package grpc

import (
	"context"

	chat "github.com/SamEkb/messenger-app/pkg/api/chat_service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ChatServer) SendMessage(ctx context.Context, req *chat.SendMessageRequest) (*chat.SendMessageResponse, error) {
	s.logger.Info("sending message")

	msg, err := s.useCase.SendMessage(ctx, req.ChatId, req.AuthorId, req.Content)
	if err != nil {
		s.logger.Error("failed to send message", "error", err)
		return nil, err
	}

	protoMsg := &chat.Message{
		MessageId: msg.ID(),
		AuthorId:  msg.AuthorID(),
		Content:   msg.Content(),
		Timestamp: timestamppb.New(msg.Timestamp()),
	}

	s.logger.Info("message sent successfully")

	return &chat.SendMessageResponse{
		Message:     protoMsg,
		Success:     true,
		MessageInfo: "message sent successfully",
	}, nil
}
