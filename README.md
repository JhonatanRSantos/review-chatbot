### Review chatbot

This repository implements a chatbot using the Gemini API.

Reduired: Valid Gemini API Key. https://ai.google.dev

#### Run

To run the API, a few steps are required:

1 - Config the env vars:
```
REVIEW_CHATBOT_DEBUG=true
REVIEW_CHATBOT_SERVER_PORT=9000
REVIEW_CHATBOT_GEN_AI_API_KEY=YOUR_APIY_KEY
REVIEW_CHATBOT_DB_HOST=127.0.0.1
REVIEW_CHATBOT_DB_PORT=3306
REVIEW_CHATBOT_DB_USER=admin
REVIEW_CHATBOT_DB_PASSWORD=qwerty
REVIEW_CHATBOT_DB_DATABASE=review-chatbot
```

2 - Run the DB:
```bash
docker-compose -f docker-compose.yaml up -d
```

3 - Run the API:
```bash
go run ./cmd/api/main.go
```