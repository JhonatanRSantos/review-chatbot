package main

import (
	"context"
	"fmt"
	"os"

	"github.com/JhonatanRSantos/review-chatbot/cmd/api/handlers"
	"github.com/JhonatanRSantos/review-chatbot/cmd/api/router"
	"github.com/JhonatanRSantos/review-chatbot/config"
	"github.com/JhonatanRSantos/review-chatbot/internal/chat"
	"github.com/JhonatanRSantos/review-chatbot/internal/chatbot"
	"github.com/JhonatanRSantos/review-chatbot/internal/user"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/JhonatanRSantos/gocore/pkg/gocontext"
	"github.com/JhonatanRSantos/gocore/pkg/godb"
	"github.com/JhonatanRSantos/gocore/pkg/goenv"
	"github.com/JhonatanRSantos/gocore/pkg/golog"
	"github.com/JhonatanRSantos/gocore/pkg/goweb"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	golog.SetEnv(goenv.Local)

	ctx := gocontext.FromContext(context.Background())
	configs := config.LoadConfiguration()

	database := newDatabaseConnection(ctx, configs)
	defer database.Close()

	userService := user.NewUserService(user.NewRepository(database))
	chatService := chat.NewChatService(chat.NewRepository(database))

	chatbotService := newChatbotService(ctx, configs)
	defer chatbotService.Close()

	ws := newWebServer(configs)
	configureWebRoutes(ws, userService, chatService, chatbotService)

	if err := ws.Listen(fmt.Sprintf(":%s", configs.ServerPort)); err != nil {
		golog.Log().Error(ctx, fmt.Sprintf("failed to start server. Cause: %s", err))
	}

}

// newDatabaseConnection
func newDatabaseConnection(ctx context.Context, configs config.Configuration) godb.DB {
	db, err := godb.NewDB(configs.Database)
	if err != nil {
		fatal(ctx, fmt.Errorf("failed to open new database connection. Cause %w", err))
	}
	return db
}

// newWebServer
func newWebServer(configs config.Configuration) *goweb.WebServer {
	ws := goweb.NewWebServer(goweb.DefaultConfig(goweb.WebServerDefaultConfig{}))

	// websocket configs
	app := ws.GetApp()
	app.Static("/", configs.StaticFilesRelativePath)
	app.Use("/api/ws", func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			ctx.Locals("allowed", true)
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	return ws
}

// newChatbotService
func newChatbotService(ctx context.Context, configs config.Configuration) *chatbot.ChatbotService {
	var (
		err    error
		bot    *chatbot.ChatbotService
		client *genai.Client
	)

	if client, err = genai.NewClient(ctx, option.WithAPIKey(configs.GenAIAPIKey)); err != nil {
		fatal(ctx, fmt.Errorf("failed to create new genai client. Cause: %w", err))
	}

	if bot, err = chatbot.NewChatbotService(ctx, chatbot.ChatbotServiceConfig{
		InitInstruction: configs.ReviewChatbotInitInstruction,
		AIClient:        client,
	}); err != nil {
		fatal(ctx, err)
	}

	return bot
}

// configureWebRoutes
func configureWebRoutes(
	ws *goweb.WebServer,
	userService *user.UserService,
	chatService *chat.ChatService,
	chatbotService *chatbot.ChatbotService,
) {
	handlers := handlers.NewHandlers(userService, chatService, chatbotService)
	ws.AddRoutes(router.NewWebRoutes(handlers)...)
}

// fatal
func fatal(ctx context.Context, message error) {
	golog.Log().Error(ctx, message.Error())
	os.Exit(1)
}
