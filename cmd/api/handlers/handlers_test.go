package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/JhonatanRSantos/review-chatbot/internal/chatbot"
	"github.com/JhonatanRSantos/review-chatbot/internal/datatypes"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type userServiceMock struct {
	Error               error
	CallbackCreate      func(ctx context.Context, firstName string, lastName string, email string) (datatypes.User, error)
	CallbackFindByEmail func(ctx context.Context, email string) (datatypes.User, error)
}

func (usm *userServiceMock) Create(ctx context.Context, firstName string, lastName string, email string) (datatypes.User, error) {
	if usm.CallbackCreate != nil {
		return usm.CallbackCreate(ctx, firstName, lastName, email)
	}
	return datatypes.User{}, usm.Error
}

func (usm *userServiceMock) FindByEmail(ctx context.Context, email string) (datatypes.User, error) {
	if usm.CallbackFindByEmail != nil {
		return usm.CallbackFindByEmail(ctx, email)
	}
	return datatypes.User{}, usm.Error
}

type chatServiceMock struct {
	Error                 error
	CallbackCreateChat    func(ctx context.Context, user datatypes.User) (string, error)
	CallbackCreateMessage func(ctx context.Context, chatID string, author string, message string) (datatypes.Message, error)
}

func (csm *chatServiceMock) CreateChat(ctx context.Context, user datatypes.User) (string, error) {
	if csm.CallbackCreateChat != nil {
		return csm.CallbackCreateChat(ctx, user)
	}
	return "", csm.Error
}

func (csm *chatServiceMock) CreateMessage(ctx context.Context, chatID string, author string, message string) (datatypes.Message, error) {
	if csm.CallbackCreateMessage != nil {
		return csm.CallbackCreateMessage(ctx, chatID, author, message)
	}
	return datatypes.Message{}, csm.Error
}

type chatbotServiceMock struct {
	CallbackStartChat func() *chatbot.ChatbotServiceSession
}

func (cbsm *chatbotServiceMock) StartChat() *chatbot.ChatbotServiceSession {
	if cbsm.CallbackStartChat != nil {
		return cbsm.CallbackStartChat()
	}
	return &chatbot.ChatbotServiceSession{}
}

func TestHandlerCreateUser(t *testing.T) {
	t.Run("should create a new user", func(t *testing.T) {
		mockedUser := datatypes.User{
			ID:        "qwerty",
			FirstName: "john",
			LastName:  "wick",
			Email:     "john.wick@continental.com",
		}

		handlers := NewHandlers(
			&userServiceMock{
				CallbackCreate: func(ctx context.Context, firstName, lastName, email string) (datatypes.User, error) {
					return mockedUser, nil
				},
			},
			&chatServiceMock{},
			&chatbotServiceMock{},
		)

		mockedUserBs, err := json.Marshal(mockedUser)
		require.NoError(t, err)

		app := fiber.New()
		path := "/api/user"

		app.Post(path, handlers.CreateUser)
		req, err := http.NewRequest("POST", path, strings.NewReader(string(mockedUserBs)))
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")

		result, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusOK, result.StatusCode)

		rawBody, err := io.ReadAll(result.Body)
		require.NoError(t, err)

		var user datatypes.User
		require.NoError(t, json.Unmarshal(rawBody, &user))
		require.EqualValues(t, mockedUser, user)
	})
}
