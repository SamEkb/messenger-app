package ports

import (
	"context"
	"time"
)

type ChatUseCase interface {
	CreateChat(ctx context.Context, participants []string) (*ChatDto, error)
	GetUserChats(ctx context.Context, userID string) ([]*ChatDto, error)
	SendMessage(ctx context.Context, chatID string, authorID, content string) (*MessageDto, error)
	GetChatHistory(ctx context.Context, chatID string) ([]*MessageDto, error)
}

type ChatDto struct {
	id           string
	participants []string
	messages     []*MessageDto
	createdAt    time.Time
	updatedAt    time.Time
}

func (c *ChatDto) ID() string {
	return c.id
}

func (c *ChatDto) Participants() []string {
	return c.participants
}

func (c *ChatDto) Messages() []*MessageDto {
	return c.messages
}

func (c *ChatDto) CreatedAt() time.Time {
	return c.createdAt
}

func (c *ChatDto) UpdatedAt() time.Time {
	return c.updatedAt
}

type MessageDto struct {
	id        string
	authorID  string
	content   string
	timestamp time.Time
}

func (m *MessageDto) ID() string {
	return m.id
}

func (m *MessageDto) AuthorID() string {
	return m.authorID
}

func (m *MessageDto) Content() string {
	return m.content
}

func (m *MessageDto) Timestamp() time.Time {
	return m.timestamp
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
