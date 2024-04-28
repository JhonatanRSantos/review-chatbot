package chat

import (
	"context"
	"testing"

	"github.com/JhonatanRSantos/review-chatbot/internal/datatypes"
	"github.com/stretchr/testify/assert"
)

type repositoryMock struct {
	Error                 error
	CallbackCreateChat    func(ctx context.Context, user datatypes.User) (string, error)
	CallbackCreateMessage func(ctx context.Context, chatID string, author string, message string) (datatypes.Message, error)
}

func (rm *repositoryMock) CreateChat(ctx context.Context, user datatypes.User) (string, error) {
	if rm.CallbackCreateChat != nil {
		return rm.CallbackCreateChat(ctx, user)
	}
	return "", rm.Error
}

func (rm *repositoryMock) CreateMessage(ctx context.Context, chatID string, author string, message string) (datatypes.Message, error) {
	if rm.CallbackCreateMessage != nil {
		return rm.CallbackCreateMessage(ctx, chatID, author, message)
	}
	return datatypes.Message{}, rm.Error
}

func TestServiceCreateChat(t *testing.T) {
	t.Run("should create a new chat", func(t *testing.T) {
		chatID := "qwerty"
		service := NewChatService(&repositoryMock{
			CallbackCreateChat: func(ctx context.Context, user datatypes.User) (string, error) {
				return chatID, nil
			},
		})

		id, err := service.CreateChat(context.Background(), datatypes.User{})
		assert.NoError(t, err)
		assert.Equal(t, chatID, id)
	})
}

func TestServiceCreateMessage(t *testing.T) {
	t.Run("should create a new message", func(t *testing.T) {
		mockedMessage := datatypes.Message{
			ID:      "qwerty",
			ChatID:  "qwerty",
			Message: "Hi",
			Author:  "user",
		}
		service := NewChatService(&repositoryMock{
			CallbackCreateMessage: func(ctx context.Context, chatID, author, message string) (datatypes.Message, error) {
				return mockedMessage, nil
			},
		})

		message, err := service.CreateMessage(context.Background(), "", "", "")
		assert.NoError(t, err)
		assert.EqualValues(t, mockedMessage, message)
	})
}
