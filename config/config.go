package config

import (
	"os"
	"strings"

	"github.com/JhonatanRSantos/gocore/pkg/godb"
)

type Configuration struct {
	ServerPort                   string
	GenAIAPIKey                  string
	StaticFilesRelativePath      string
	ReviewChatbotInitInstruction string
	Database                     godb.DBConfig
}

func LoadConfiguration() Configuration {
	config := Configuration{
		ServerPort:                   os.Getenv("REVIEW_CHATBOT_SERVER_PORT"),
		GenAIAPIKey:                  os.Getenv("REVIEW_CHATBOT_GEN_AI_API_KEY"),
		StaticFilesRelativePath:      "./static",
		ReviewChatbotInitInstruction: reviewChatbotInitInstruction,
		Database: godb.DBConfig{
			Host:             os.Getenv("REVIEW_CHATBOT_DB_HOST"),
			Port:             os.Getenv("REVIEW_CHATBOT_DB_PORT"),
			User:             os.Getenv("REVIEW_CHATBOT_DB_USER"),
			Password:         os.Getenv("REVIEW_CHATBOT_DB_PASSWORD"),
			Database:         os.Getenv("REVIEW_CHATBOT_DB_DATABASE"),
			DatabaseType:     godb.MySQLDB,
			ConnectionParams: map[string]string{},
		},
	}

	if strings.ToLower(strings.TrimSpace(os.Getenv("REVIEW_CHATBOT_DEBUG"))) == "true" {
		config.StaticFilesRelativePath = "../../static"
	}

	return config
}
