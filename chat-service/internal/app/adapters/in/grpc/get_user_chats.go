package grpc

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	chat "github.com/SamEkb/messenger-app/pkg/api/chat_service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ChatServer) GetUserChats(ctx context.Context, req *chat.GetUserChatsRequest) (*chat.GetUserChatsResponse, error) {
	s.logger.Info("getting user chats")

	chats, err := s.useCase.GetUserChats(ctx, req.UserId)
	if err != nil {
		s.logger.Error("failed to get user chats", "error", err)
		return nil, err
	}

	s.logger.Info("user chats retrieved successfully")

	return &chat.GetUserChatsResponse{Chats: mapDtoToProto(chats)}, nil
}

func mapDtoToProto(chats []*ports.ChatDto) []*chat.Chat {
	var chatsProto []*chat.Chat

	for _, v := range chats {
		var lastMsg *chat.Message
		msgs := v.Messages()
		if len(msgs) > 0 {
			m := msgs[len(msgs)-1]
			lastMsg = &chat.Message{
				MessageId: m.ID(),
				AuthorId:  m.AuthorID(),
				Content:   m.Content(),
				Timestamp: timestamppb.New(m.Timestamp()),
			}
		}
		ch := &chat.Chat{
			ChatId:       v.ID(),
			Participants: v.Participants(),
			LastMessage:  lastMsg,
		}
		chatsProto = append(chatsProto, ch)
	}
	return chatsProto
}
