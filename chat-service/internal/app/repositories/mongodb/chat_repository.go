package mongodb

import (
	"context"
	"time"

	"github.com/SamEkb/messenger-app/chat-service/config/env"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/models"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ ports.ChatRepository = (*ChatRepository)(nil)

type ChatRepository struct {
	db     *mongo.Database
	logger logger.Logger
	config *env.MongoDBConfig
}

type chatDocument struct {
	ID           string        `bson:"_id"`
	Participants []string      `bson:"participants"`
	Messages     []msgDocument `bson:"messages"`
	CreatedAt    time.Time     `bson:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at"`
}

type msgDocument struct {
	ID        string    `bson:"_id"`
	AuthorID  string    `bson:"author_id"`
	Content   string    `bson:"content"`
	Timestamp time.Time `bson:"timestamp"`
}

func NewChatRepository(client *mongo.Client, dbName string, logger logger.Logger) *ChatRepository {
	db := client.Database(dbName)

	_, err := db.Collection("chats").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "participants", Value: 1}},
			Options: options.Index().SetBackground(true),
		},
	)

	if err != nil {
		logger.Error("failed to create index", "error", err)
	}

	return &ChatRepository{
		db:     db,
		logger: logger.With("component", "chat_repository"),
	}
}

func (r *ChatRepository) Create(ctx context.Context, participants []string) (*models.Chat, error) {
	r.logger.Debug("creating chat", "participants", participants)

	chat, err := models.NewChat(participants)
	if err != nil {
		r.logger.Error("failed to create chat model", "error", err)
		return nil, err
	}

	doc := chatDocument{
		ID:           chat.ID().String(),
		Participants: chat.Participants(),
		Messages:     []msgDocument{},
		CreatedAt:    chat.CreatedAt(),
		UpdatedAt:    chat.UpdatedAt(),
	}

	_, err = r.db.Collection("chats").InsertOne(ctx, doc)
	if err != nil {
		r.logger.Error("failed to insert chat", "error", err)
		return nil, errors.NewInternalError(err, "failed to create chat")
	}

	r.logger.Info("chat created", "chat_id", chat.ID())
	return chat, nil
}

func (r *ChatRepository) Get(ctx context.Context, userID string) ([]*models.Chat, error) {
	r.logger.Debug("getting chats", "user_id", userID)

	filter := bson.M{"participants": userID}
	cursor, err := r.db.Collection("chats").Find(ctx, filter)
	if err != nil {
		r.logger.Error("failed to find chats", "error", err)
		return nil, errors.NewInternalError(err, "failed to get chats")
	}
	defer cursor.Close(ctx)

	var docs []chatDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("failed to decode chats", "error", err)
		return nil, errors.NewInternalError(err, "failed to decode chats")
	}

	chats := make([]*models.Chat, 0, len(docs))
	for _, doc := range docs {
		chat, err := r.documentToModel(doc)
		if err != nil {
			r.logger.Error("failed to convert document to model", "error", err)
			continue
		}
		chats = append(chats, chat)
	}

	r.logger.Debug("got chats", "user_id", userID, "count", len(chats))
	return chats, nil
}

func (r *ChatRepository) SendMessage(ctx context.Context, chatID models.ChatID, authorID, content string) (*models.Message, error) {
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > r.config.Timeout {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.config.Timeout)
		defer cancel()
	}

	r.logger.Debug("sending message", "chat_id", chatID, "author_id", authorID)

	message, err := models.NewMessage(authorID, content)
	if err != nil {
		r.logger.Error("failed to create message model", "error", err)
		return nil, err
	}

	msgDoc := msgDocument{
		ID:        message.ID().String(),
		AuthorID:  message.AuthorID(),
		Content:   message.Content(),
		Timestamp: message.Timestamp(),
	}

	filter := bson.M{"_id": chatID.String()}
	update := bson.M{
		"$push": bson.M{"messages": msgDoc},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.db.Collection("chats").UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("failed to update chat", "error", err)
		return nil, errors.NewInternalError(err, "failed to send message")
	}

	if result.MatchedCount == 0 {
		r.logger.Error("chat not found", "chat_id", chatID)
		return nil, errors.NewNotFoundError("chat not found")
	}

	r.logger.Info("message sent", "chat_id", chatID, "message_id", message.ID())
	return message, nil
}

func (r *ChatRepository) GetMessages(ctx context.Context, chatID models.ChatID) ([]*models.Message, error) {
	r.logger.Debug("getting messages", "chat_id", chatID)

	filter := bson.M{"_id": chatID.String()}
	var doc chatDocument
	err := r.db.Collection("chats").FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Error("chat not found", "chat_id", chatID)
			return nil, errors.NewNotFoundError("chat not found")
		}
		r.logger.Error("failed to find chat", "error", err)
		return nil, errors.NewInternalError(err, "failed to get messages")
	}

	messages := make([]*models.Message, 0, len(doc.Messages))
	for _, msgDoc := range doc.Messages {
		_, err := models.ParseMessageID(msgDoc.ID)
		if err != nil {
			r.logger.Error("failed to parse message ID", "error", err)
			continue
		}

		msg, err := models.NewMessage(msgDoc.AuthorID, msgDoc.Content)
		if err != nil {
			r.logger.Error("failed to create message model", "error", err)
			continue
		}

		messages = append(messages, msg)
	}

	r.logger.Debug("got messages", "chat_id", chatID, "count", len(messages))
	return messages, nil
}

func (r *ChatRepository) documentToModel(doc chatDocument) (*models.Chat, error) {
	_, err := models.ParseChatID(doc.ID)
	if err != nil {
		return nil, err
	}

	chat, err := models.NewChat(doc.Participants)
	if err != nil {
		return nil, err
	}

	for _, msgDoc := range doc.Messages {
		_, err := uuid.Parse(msgDoc.ID)
		if err != nil {
			continue
		}

		msg, err := models.NewMessage(msgDoc.AuthorID, msgDoc.Content)
		if err != nil {
			continue
		}

		chat.AddMessage(msg)
	}

	return chat, nil
}
