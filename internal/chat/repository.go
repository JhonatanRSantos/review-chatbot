package chat

import (
	"context"
	"fmt"

	"github.com/JhonatanRSantos/review-chatbot/internal/datatypes"
	"github.com/gofrs/uuid/v5"

	"github.com/JhonatanRSantos/gocore/pkg/godb"
)

type Repository struct {
	db godb.DB
}

func NewRepository(
	db godb.DB,
) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateChat(ctx context.Context, user datatypes.User) (string, error) {
	stm, err := r.db.PrepareNamedContext(ctx, createChat)
	if err != nil {
		return "", fmt.Errorf("failed to create new chat. Cause: %w", err)
	}
	defer stm.Close()

	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to create new chat. Cause: %w", err)
	}

	params := map[string]interface{}{
		"id":      id.String(),
		"user_id": user.ID,
	}

	result, err := stm.ExecContext(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to create new chat. Cause: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to create new chat. Cause: %w", err)
	}

	if rows == 0 {
		return "", ErrCantCreateChat
	}

	return id.String(), nil
}

func (r *Repository) CreateMessage(ctx context.Context, chatID string, author string, message string) (datatypes.Message, error) {
	stm, err := r.db.PrepareNamedContext(ctx, createMessage)
	if err != nil {
		return datatypes.Message{}, fmt.Errorf("failed to create new message. Cause: %w", err)
	}
	defer stm.Close()

	id, err := uuid.NewV4()
	if err != nil {
		return datatypes.Message{}, fmt.Errorf("failed to create new message. Cause: %w", err)
	}

	params := map[string]interface{}{
		"id":      id.String(),
		"chat_id": chatID,
		"author":  author,
		"message": message,
	}

	result, err := stm.ExecContext(ctx, params)
	if err != nil {
		return datatypes.Message{}, fmt.Errorf("failed to create new message. Cause: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return datatypes.Message{}, fmt.Errorf("failed to create new message. Cause: %w", err)
	}

	if rows == 0 {
		return datatypes.Message{}, ErrCantCreateMessage
	}

	return datatypes.Message{
		ID:      chatID,
		ChatID:  chatID,
		Author:  author,
		Message: message,
	}, nil
}
