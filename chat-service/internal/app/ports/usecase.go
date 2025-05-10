package ports

import (
	"context"
	"time"
)

type ChatUseCase interface {
	CreateChat(ctx context.Context, participants []string) (*ChatDto, error)
	GetUserChats(ctx context.Context, userID string) ([]*ChatDto, error)
	SendMessage(ctx context.Context, chatID string, authorID, content string) error
	GetChatHistory(ctx context.Context, chatID string) ([]*MessageDto, error)
}

type ChatDto struct {
	id           string
	participants []string
	messages     []*MessageDto
	createdAt    time.Time
	updatedAt    time.Time
}

type MessageDto struct {
	id        string
	authorID  string
	content   string
	timestamp time.Time
}

func NewChatDto(id string, participants []string, messages []*MessageDto, createdAt, updatedAt time.Time) *ChatDto {
	return &ChatDto{
		id:           id,
		participants: participants,
		messages:     messages,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func NewMessageDto(id, authorID, content string, timestamp time.Time) *MessageDto {
	return &MessageDto{
		id:        id,
		authorID:  authorID,
		content:   content,
		timestamp: timestamp,
	}
}
