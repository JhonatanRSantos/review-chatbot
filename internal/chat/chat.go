package chat

import (
	"context"

	"github.com/JhonatanRSantos/review-chatbot/internal/datatypes"
)

type repository interface {
	CreateChat(ctx context.Context, user datatypes.User) (string, error)
	CreateMessage(ctx context.Context, chatID string, author string, message string) (datatypes.Message, error)
}

type ChatService struct {
	repository repository
}

func NewChatService(
	repository repository,
) *ChatService {
	return &ChatService{repository}
}

func (cs *ChatService) CreateChat(ctx context.Context, user datatypes.User) (string, error) {
	return cs.repository.CreateChat(ctx, user)
}

func (cs *ChatService) CreateMessage(ctx context.Context, chatID string, author string, message string) (datatypes.Message, error) {
	return cs.repository.CreateMessage(ctx, chatID, author, message)
}
