package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/JhonatanRSantos/gocore/pkg/godb"
	"github.com/JhonatanRSantos/review-chatbot/internal/datatypes"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryCreate(t *testing.T) {
	t.Run("should create new user", func(t *testing.T) {
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

		mockedUser := getMockedUser(t)
		user, err := repository.Create(
			context.Background(), mockedUser.FirstName, mockedUser.LastName, mockedUser.Email,
		)
		assert.NoError(t, err)
		assert.EqualValues(t, mockedUser.FirstName, user.FirstName)
		assert.EqualValues(t, mockedUser.LastName, user.LastName)
		assert.EqualValues(t, mockedUser.Email, user.Email)
	})

	t.Run("should fail to prepare named context", func(t *testing.T) {
		errPrepareNamedContext := errors.New("error when preparing named context for tests")
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return &godb.NamedStmtMock{}, errPrepareNamedContext
			},
		})

		_, err := repository.Create(context.Background(), "", "", "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errPrepareNamedContext)
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

		_, err := repository.Create(context.Background(), "", "", "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errExecContext)
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

		_, err := repository.Create(context.Background(), "", "", "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errRowsAffected)
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

		_, err := repository.Create(context.Background(), "", "", "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrCantSaveUser)
	})
}

func TestRepositoryFindByEmail(t *testing.T) {
	t.Run("should find user", func(t *testing.T) {
		mockedUser := getMockedUser(t)
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return &godb.NamedStmtMock{
					CallbackGetContext: func(ctx context.Context, dest, arg interface{}) error {
						if user, ok := dest.(*datatypes.User); ok {
							*user = mockedUser
						}
						return nil
					},
				}, nil
			},
		})

		user, err := repository.FindByEmail(context.Background(), mockedUser.Email)
		assert.NoError(t, err)
		assert.EqualValues(t, mockedUser, user)
	})

	t.Run("should fail to prepare named context", func(t *testing.T) {
		errPrepareNamedContext := errors.New("error when preparing named context for tests")
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return nil, errPrepareNamedContext
			},
		})

		_, err := repository.FindByEmail(context.Background(), "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errPrepareNamedContext)
	})

	t.Run("should fail run get context", func(t *testing.T) {
		errGetContext := errors.New("error to run get context for tests")
		repository := NewRepository(&godb.DBMock{
			CallbackPrepareNamedContext: func(ctx context.Context, query string) (godb.NamedStmt, error) {
				return &godb.NamedStmtMock{
					CallbackGetContext: func(ctx context.Context, dest, arg interface{}) error {
						return errGetContext
					},
				}, nil
			},
		})

		_, err := repository.FindByEmail(context.Background(), "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errGetContext)
	})
}
