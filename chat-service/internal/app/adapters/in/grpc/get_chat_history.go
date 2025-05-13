package grpc

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	chat "github.com/SamEkb/messenger-app/pkg/api/chat_service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ChatServer) GetChatHistory(ctx context.Context, req *chat.GetChatHistoryRequest) (*chat.GetChatHistoryResponse, error) {
	s.logger.Info("getting chat history")

	history, err := s.useCase.GetChatHistory(ctx, req.ChatId)
	if err != nil {
		s.logger.Error("failed to get chat history", "error", err)
		return nil, err
	}

	s.logger.Info("chat history retrieved successfully")

	totalMsg := int32(len(history))
	protoMsgs := dtoToProto(history)
	return &chat.GetChatHistoryResponse{
		Messages:      protoMsgs,
		TotalMessages: totalMsg,
	}, nil
}

func dtoToProto(msgs []*ports.MessageDto) []*chat.Message {
	var msgsProto []*chat.Message
	for _, m := range msgs {
		msgsProto = append(msgsProto, &chat.Message{
			MessageId: m.ID(),
			AuthorId:  m.AuthorID(),
			Content:   m.Content(),
			Timestamp: timestamppb.New(m.Timestamp()),
		})
	}
	return msgsProto
}
