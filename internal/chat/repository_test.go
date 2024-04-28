package chat

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/JhonatanRSantos/gocore/pkg/godb"
	"github.com/JhonatanRSantos/review-chatbot/internal/datatypes"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryCreateChat(t *testing.T) {
	t.Run("should create a new chat", func(t *testing.T) {
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return &godb.NamedStmtMock{
					CallbackExecContext: func(ctx context.Context, arg interface{}) (sql.Result, error) {
						return &godb.ResultMock{
							CallbackRowsAffected: func() (int64, error) {
								return 1, nil
							},
						}, nil
					},
				}, nil
			},
		})

		id, err := repository.CreateChat(context.Background(), datatypes.User{})
		assert.NoError(t, err)
		assert.True(t, len(id) > 0)
	})

	t.Run("should fail when prapare named context", func(t *testing.T) {
		errPrepareNamedContext := errors.New("error when preparing named context for tests")
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return nil, errPrepareNamedContext
			},
		})

		id, err := repository.CreateChat(context.Background(), datatypes.User{})
		assert.Error(t, err)
		assert.ErrorIs(t, err, errPrepareNamedContext)
		assert.Equal(t, "", id)
	})

	t.Run("should fail to run exec context", func(t *testing.T) {
		errExecContext := errors.New("error to run exec context for tests")
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return &godb.NamedStmtMock{
					CallbackExecContext: func(ctx context.Context, arg interface{}) (sql.Result, error) {
						return nil, errExecContext
					},
				}, nil
			},
		})

		id, err := repository.CreateChat(context.Background(), datatypes.User{})
		assert.Error(t, err)
		assert.ErrorIs(t, err, errExecContext)
		assert.Equal(t, "", id)
	})

	t.Run("should fail when to run affected rows", func(t *testing.T) {
		errRowsAffected := errors.New("error to run affected rows for tests")
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return &godb.NamedStmtMock{
					CallbackExecContext: func(ctx context.Context, arg interface{}) (sql.Result, error) {
						return &godb.ResultMock{
							CallbackRowsAffected: func() (int64, error) {
								return 0, errRowsAffected
							},
						}, nil
					},
				}, nil
			},
		})

		id, err := repository.CreateChat(context.Background(), datatypes.User{})
		assert.Error(t, err)
		assert.ErrorIs(t, err, errRowsAffected)
		assert.Equal(t, "", id)
	})

	t.Run("should fail when there ara no affected rows", func(t *testing.T) {
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return &godb.NamedStmtMock{
					CallbackExecContext: func(ctx context.Context, arg interface{}) (sql.Result, error) {
						return &godb.ResultMock{
							CallbackRowsAffected: func() (int64, error) {
								return 0, nil
							},
						}, nil
					},
				}, nil
			},
		})

		id, err := repository.CreateChat(context.Background(), datatypes.User{})
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrCantCreateChat)
		assert.Equal(t, "", id)
	})
}

func TestRepositoryCreateMessage(t *testing.T) {
	t.Run("should create a new message", func(t *testing.T) {
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return &godb.NamedStmtMock{
					CallbackExecContext: func(ctx context.Context, arg interface{}) (sql.Result, error) {
						return &godb.ResultMock{
							CallbackRowsAffected: func() (int64, error) {
								return 1, nil
							},
						}, nil
					},
				}, nil
			},
		})

		mockedMessage := datatypes.Message{
			ChatID:  "qwert",
			Message: "test",
			Author:  "user",
		}
		message, err := repository.CreateMessage(context.Background(), mockedMessage.ChatID, mockedMessage.Author, mockedMessage.Message)
		assert.NoError(t, err)
		assert.Equal(t, mockedMessage.ChatID, message.ChatID)
		assert.Equal(t, mockedMessage.Author, message.Author)
		assert.Equal(t, mockedMessage.Message, message.Message)
	})

	t.Run("should fail when prapare named context", func(t *testing.T) {
		errPrepareNamedContext := errors.New("error when preparing named context for tests")
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return nil, errPrepareNamedContext
			},
		})

		message, err := repository.CreateMessage(context.Background(), "", "", "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errPrepareNamedContext)
		assert.Equal(t, datatypes.Message{}, message)
	})

	t.Run("should fail to run exec context", func(t *testing.T) {
		errExecContext := errors.New("error to run exec context for tests")
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return &godb.NamedStmtMock{
					CallbackExecContext: func(ctx context.Context, arg interface{}) (sql.Result, error) {
						return nil, errExecContext
					},
				}, nil
			},
		})

		message, err := repository.CreateMessage(context.Background(), "", "", "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errExecContext)
		assert.Equal(t, datatypes.Message{}, message)
	})

	t.Run("should fail when to run affected rows", func(t *testing.T) {
		errRowsAffected := errors.New("error to run affected rows for tests")
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return &godb.NamedStmtMock{
					CallbackExecContext: func(ctx context.Context, arg interface{}) (sql.Result, error) {
						return &godb.ResultMock{
							CallbackRowsAffected: func() (int64, error) {
								return 0, errRowsAffected
							},
						}, nil
					},
				}, nil
			},
		})

		message, err := repository.CreateMessage(context.Background(), "", "", "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errRowsAffected)
		assert.Equal(t, datatypes.Message{}, message)
	})

	t.Run("should fail when there ara no affected rows", func(t *testing.T) {
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return &godb.NamedStmtMock{
					CallbackExecContext: func(ctx context.Context, arg interface{}) (sql.Result, error) {
						return &godb.ResultMock{
							CallbackRowsAffected: func() (int64, error) {
								return 0, nil
							},
						}, nil
					},
				}, nil
			},
		})

		message, err := repository.CreateMessage(context.Background(), "", "", "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrCantCreateMessage)
		assert.Equal(t, datatypes.Message{}, message)
	})
}
