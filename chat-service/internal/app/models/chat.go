package models

import (
	"time"

	"github.com/SamEkb/messenger-app/chat-service/pkg/errors"
	"github.com/google/uuid"
)

type MessageID uuid.UUID
type ChatID uuid.UUID

type Message struct {
	id        MessageID
	authorID  string
	content   string
	timestamp time.Time
}

type Chat struct {
	id           ChatID
	participants []string
	messages     []Message
	createdAt    time.Time
	updatedAt    time.Time
}

func NewChat(participants []string) (*Chat, error) {
	if len(participants) == 0 {
		return nil, errors.NewInvalidInputError("participants are required")
	}

	chatID := ChatID(uuid.New())
	return &Chat{
		id:           chatID,
		participants: participants,
		messages:     []Message{},
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
	}, nil
}

func NewMessage(authorID, content string) (*Message, error) {
	if authorID == "" {
		return nil, errors.NewInvalidInputError("author ID is required")
	}
	if content == "" {
		return nil, errors.NewInvalidInputError("content is required")
	}

	msgID := MessageID(uuid.New())
	return &Message{
		id:        msgID,
		authorID:  authorID,
		content:   content,
		timestamp: time.Now(),
	}, nil
}

func (u MessageID) IsEmpty() bool {
	return u == MessageID(uuid.Nil)
}

func (u MessageID) String() string {
	return uuid.UUID(u).String()
}

func (u ChatID) IsEmpty() bool {
	return u == ChatID(uuid.Nil)
}

func (u ChatID) String() string {
	return uuid.UUID(u).String()
}

func (m *Message) ID() MessageID {
	return m.id
}

func (m *Message) AuthorID() string {
	return m.authorID
}

func (m *Message) Content() string {
	return m.content
}

func (m *Message) Timestamp() time.Time {
	return m.timestamp
}

func (c *Chat) ID() ChatID {
	return c.id
}

func (c *Chat) Participants() []string {
	return c.participants
}

func (c *Chat) Messages() []Message {
	return c.messages
}

func (c *Chat) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Chat) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Chat) AddMessage(msg *Message) error {
	if msg == nil {
		return errors.NewInvalidInputError("message is required")
	}
	c.messages = append(c.messages, *msg)
	c.updatedAt = time.Now()
	return nil
}
