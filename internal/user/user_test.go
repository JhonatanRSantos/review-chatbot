package user

import (
	"context"
	"errors"
	"testing"

	"github.com/JhonatanRSantos/review-chatbot/internal/datatypes"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
)

type repositoryMock struct {
	Error               error
	CallbackCreate      func(ctx context.Context, firstName string, lastName string, email string) (datatypes.User, error)
	CallbackFindByEmail func(ctx context.Context, email string) (datatypes.User, error)
}

func (rm *repositoryMock) Create(ctx context.Context, firstName string, lastName string, email string) (datatypes.User, error) {
	if rm.CallbackCreate != nil {
		return rm.CallbackCreate(ctx, firstName, lastName, email)
	}
	return datatypes.User{}, rm.Error
}

func (rm *repositoryMock) FindByEmail(ctx context.Context, email string) (datatypes.User, error) {
	if rm.CallbackFindByEmail != nil {
		return rm.CallbackFindByEmail(ctx, email)
	}
	return datatypes.User{}, rm.Error
}

func TestServicCreate(t *testing.T) {
	t.Run("should create new user", func(t *testing.T) {
		mockedUser := getMockedUser(t)
		service := NewUserService(&repositoryMock{
			CallbackCreate: func(ctx context.Context, firstName, lastName, email string) (datatypes.User, error) {
				return mockedUser, nil
			},
		})

		user, err := service.Create(context.Background(), mockedUser.FirstName, mockedUser.LastName, mockedUser.Email)
		assert.NoError(t, err)
		assert.EqualValues(t, mockedUser, user)
	})

	t.Run("should fail when required fields are missing", func(t *testing.T) {
		service := NewUserService(&repositoryMock{
			Error: errors.New("invalid service call when testing"),
		})

		_, err := service.Create(context.Background(), "", "", "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingRequiredFields)
	})
}

func TestServiceFindByEmail(t *testing.T) {
	t.Run("should find user", func(t *testing.T) {
		mockedUser := getMockedUser(t)
		service := NewUserService(&repositoryMock{
			CallbackFindByEmail: func(ctx context.Context, email string) (datatypes.User, error) {
				return mockedUser, nil
			},
		})

		user, err := service.FindByEmail(context.Background(), "john.wick@continental.com")
		assert.NoError(t, err)
		assert.EqualValues(t, mockedUser, user)
	})

	t.Run("should fail when required fields are missing", func(t *testing.T) {
		service := NewUserService(&repositoryMock{
			Error: errors.New("invalid service call when testing"),
		})

		_, err := service.FindByEmail(context.Background(), "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingRequiredFields)
	})
}

func getMockedUser(t *testing.T) datatypes.User {
	id, err := uuid.NewV4()
	if err != nil {
		assert.FailNow(t, "failed to create user id. %s", err)
	}

	return datatypes.User{
		ID:        id.String(),
		FirstName: "john",
		LastName:  "wick",
		Email:     "john.wick@continental.com",
	}
}
