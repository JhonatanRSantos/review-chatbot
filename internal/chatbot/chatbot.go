package chatbot

import (
	"context"
	"fmt"

	"github.com/JhonatanRSantos/gocore/pkg/golog"
	"github.com/google/generative-ai-go/genai"
)

type ChatbotServiceSession struct {
	session *genai.ChatSession
}

// SendTextMessage
func (rcss *ChatbotServiceSession) SendTextMessage(ctx context.Context, message string) string {
	resp, err := rcss.session.SendMessage(ctx, genai.Text(message))

	messageResponse := "I'm sorry but can't help you right now. Can you please try later."

	if err != nil {
		golog.Log().Error(ctx, err.Error())
		return messageResponse
	}

	for _, candidate := range resp.Candidates {
		if len(candidate.Content.Parts) > 0 {
			if value, ok := candidate.Content.Parts[0].(genai.Text); ok {
				messageResponse = string(value)
			}
			break
		}
	}

	return messageResponse
}

type ChatbotServiceConfig struct {
	InitInstruction string
	AIClient        aiClient
}

// validate check if configs are valid
func (rcc ChatbotServiceConfig) validate() error {
	if rcc.InitInstruction == "" || rcc.AIClient == nil {
		return ErrMissingChatbotConfigs
	}
	return nil
}

type ChatbotService struct {
	model  *genai.GenerativeModel
	client aiClient
}

// NewChatbotService create a new review chatbot
func NewChatbotService(
	ctx context.Context,
	config ChatbotServiceConfig,
) (*ChatbotService, error) {
	baseError := "failed to create review chatbot. Cause: %w"

	if err := config.validate(); err != nil {
		return &ChatbotService{}, fmt.Errorf(baseError, err)
	}

	model := config.AIClient.GenerativeModel("gemini-1.5-pro-latest")
	model.GenerationConfig.SetTopK(0)
	model.GenerationConfig.SetTopP(0.95)
	model.GenerationConfig.SetMaxOutputTokens(8192)
	model.GenerationConfig.SetTemperature(1)

	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(config.InitInstruction),
		},
	}

	return &ChatbotService{
		model:  model,
		client: config.AIClient,
	}, nil
}

// Close closes the client connection
func (rc *ChatbotService) Close() error {
	if rc.model != nil && rc.client != nil {
		if err := rc.client.Close(); err != nil {
			return err
		}
	}
	rc.model = nil
	rc.client = nil
	return nil
}

// StartChat starts a chat session.
func (rc *ChatbotService) StartChat() *ChatbotServiceSession {
	return &ChatbotServiceSession{
		session: rc.model.StartChat(),
	}
}
