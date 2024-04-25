package handlers

import (
	"context"
	"fmt"
	"sync"

	"github.com/JhonatanRSantos/gocore/pkg/gocontext"
	"github.com/JhonatanRSantos/gocore/pkg/golog"
	"github.com/JhonatanRSantos/review-chatbot/internal/chatbot"
	"github.com/JhonatanRSantos/review-chatbot/internal/datatypes"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type connection struct {
	conn        *websocket.Conn
	chatID      string
	chatSession *chatbot.ChatbotServiceSession
}

type userService interface {
	Create(ctx context.Context, firstName string, lastName string, email string) (datatypes.User, error)
	FindByEmail(ctx context.Context, email string) (datatypes.User, error)
}

type chatService interface {
	CreateChat(ctx context.Context, user datatypes.User) (string, error)
	CreateMessage(ctx context.Context, chatID string, author string, message string) (datatypes.Message, error)
}

type chatbotService interface {
	StartChat() *chatbot.ChatbotServiceSession
}

type Handlers struct {
	sessions       map[string]connection
	sessionMutex   *sync.RWMutex
	userService    userService
	chatService    chatService
	chatbotService chatbotService
}

// NewHandlers
func NewHandlers(
	userService userService,
	chatService chatService,
	chatbotService chatbotService,
) *Handlers {
	return &Handlers{
		sessions:       make(map[string]connection),
		sessionMutex:   &sync.RWMutex{},
		userService:    userService,
		chatService:    chatService,
		chatbotService: chatbotService,
	}
}

// CreateUser
func (h *Handlers) CreateUser(fc *fiber.Ctx) error {
	ctx := gocontext.FromContext(fc.Context())

	var req datatypes.CreateUserRequest

	if err := fc.BodyParser(&req); err != nil {
		golog.Log().Error(ctx, err.Error())
		return fc.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := req.Validate(); err != nil {
		golog.Log().Error(ctx, err.Error())
		return fc.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	user, err := h.userService.Create(ctx, req.FirstName, req.LastName, req.Email)
	if err != nil {
		golog.Log().Error(ctx, err.Error())
		return fc.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return fc.JSON(user)
}

// CreateReview
func (h *Handlers) CreateReview(fc *fiber.Ctx) error {
	ctx := gocontext.FromContext(fc.Context())
	var req datatypes.CreateReviewRequest

	if err := fc.BodyParser(&req); err != nil {
		golog.Log().Error(ctx, err.Error())
		return fc.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := req.Validate(); err != nil {
		golog.Log().Error(ctx, err.Error())
		return fc.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	h.sessionMutex.Lock()
	defer h.sessionMutex.Unlock()

	session, ok := h.sessions[req.User.Email]
	if !ok {
		return fc.SendStatus(fiber.StatusNotFound)
	}

	message := fmt.Sprintf("Start a new review with %s. He just bought a new %s", req.User.Name, req.Product)
	messageResponse := session.chatSession.SendTextMessage(fc.Context(), message)

	if _, err := h.chatService.CreateMessage(fc.Context(), session.chatID, "chatbot", messageResponse); err != nil {
		golog.Log().Error(ctx, err.Error())
		return fc.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if err := session.conn.WriteMessage(websocket.TextMessage, []byte(messageResponse)); err != nil {
		delete(h.sessions, req.User.Email)
		golog.Log().Error(ctx, err.Error())
		return fc.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return fc.SendStatus(fiber.StatusOK)
}

// HandleWebsocketConnection
func (h *Handlers) HandleWebsocketConnection() func(*fiber.Ctx) error {
	return websocket.New(func(conn *websocket.Conn) {
		ctx := gocontext.FromContext(context.Background())

		user, err := h.userService.FindByEmail(ctx, conn.Params("email"))
		if err != nil {
			golog.Log().Error(ctx, err.Error())
			conn.Close()
			return
		}

		chatID, err := h.chatService.CreateChat(ctx, user)
		if err != nil {
			golog.Log().Error(ctx, err.Error())
			conn.Close()
			return
		}

		chatSession := h.chatbotService.StartChat()
		h.sessionMutex.Lock()
		h.sessions[user.Email] = connection{
			conn:        conn,
			chatID:      chatID,
			chatSession: chatSession,
		}
		h.sessionMutex.Unlock()

		removeConnection := func() {
			h.sessionMutex.Lock()
			delete(h.sessions, user.Email)
			h.sessionMutex.Unlock()
		}

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				removeConnection()
				golog.Log().Error(ctx, fmt.Sprintf("failed to read message. Cause: %s", err))
				break
			}

			if _, err = h.chatService.CreateMessage(ctx, chatID, "user", string(message)); err != nil {
				removeConnection()
				golog.Log().Error(ctx, err.Error())
				break
			}

			messageResponse := chatSession.SendTextMessage(ctx, string(message))

			if _, err = h.chatService.CreateMessage(ctx, chatID, "chatbot", messageResponse); err != nil {
				removeConnection()
				golog.Log().Error(ctx, err.Error())
				break
			}

			if err = conn.WriteMessage(messageType, []byte(messageResponse)); err != nil {
				removeConnection()
				golog.Log().Error(ctx, fmt.Sprintf("failed to write message. Cause: %s", err))
				break
			}
		}
	})
}
