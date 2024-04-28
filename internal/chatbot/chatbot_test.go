package chatbot

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/stretchr/testify/assert"
)

type aiClientMock struct {
	Error                   error
	CallbackClose           func() error
	CallbackDeleteFile      func(ctx context.Context, name string) error
	CallbackEmbeddingModel  func(name string) *genai.EmbeddingModel
	CallbackGenerativeModel func(name string) *genai.GenerativeModel
	CallbackGetFile         func(ctx context.Context, name string) (*genai.File, error)
	CallbackListFiles       func(ctx context.Context) *genai.FileIterator
	CallbackListModels      func(ctx context.Context) *genai.ModelInfoIterator
	CallbackUploadFile      func(ctx context.Context, name string, r io.Reader, opts *genai.UploadFileOptions) (*genai.File, error)
}

func (aicm *aiClientMock) Close() error {
	if aicm.CallbackClose != nil {
		return aicm.CallbackClose()
	}
	return aicm.Error
}

func (aicm *aiClientMock) DeleteFile(ctx context.Context, name string) error {
	if aicm.CallbackDeleteFile != nil {
		return aicm.CallbackDeleteFile(ctx, name)
	}
	return aicm.Error
}

func (aicm *aiClientMock) EmbeddingModel(name string) *genai.EmbeddingModel {
	if aicm.CallbackEmbeddingModel != nil {
		return aicm.CallbackEmbeddingModel(name)
	}
	return &genai.EmbeddingModel{}
}

func (aicm *aiClientMock) GenerativeModel(name string) *genai.GenerativeModel {
	if aicm.CallbackGenerativeModel != nil {
		return aicm.CallbackGenerativeModel(name)
	}
	return &genai.GenerativeModel{}
}

func (aicm *aiClientMock) GetFile(ctx context.Context, name string) (*genai.File, error) {
	if aicm.CallbackGetFile != nil {
		return aicm.CallbackGetFile(ctx, name)
	}
	return &genai.File{}, aicm.Error
}

func (aicm *aiClientMock) ListFiles(ctx context.Context) *genai.FileIterator {
	if aicm.CallbackListFiles != nil {
		return aicm.CallbackListFiles(ctx)
	}
	return &genai.FileIterator{}
}

func (aicm *aiClientMock) ListModels(ctx context.Context) *genai.ModelInfoIterator {
	if aicm.CallbackListModels != nil {
		return aicm.CallbackListModels(ctx)
	}
	return &genai.ModelInfoIterator{}
}

func (aicm *aiClientMock) UploadFile(ctx context.Context, name string, r io.Reader, opts *genai.UploadFileOptions) (*genai.File, error) {
	if aicm.CallbackUploadFile != nil {
		return aicm.CallbackUploadFile(ctx, name, r, opts)
	}
	return &genai.File{}, aicm.Error
}

func TestChatbotService(t *testing.T) {
	t.Run("should create a new chatbot service", func(t *testing.T) {
		service, err := NewChatbotService(context.Background(), ChatbotServiceConfig{
			InitInstruction: "abcde",
			AIClient:        &aiClientMock{},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, service)
	})

	t.Run("should fail to create a new chatbot service", func(t *testing.T) {
		service, err := NewChatbotService(context.Background(), ChatbotServiceConfig{
			InitInstruction: "",
			AIClient:        &aiClientMock{},
		})
		assert.Error(t, err)
		assert.Empty(t, service)
	})

	t.Run("should fail to close chatbot service", func(t *testing.T) {
		errClose := errors.New("failed to run close for tests")
		service, err := NewChatbotService(context.Background(), ChatbotServiceConfig{
			InitInstruction: "abcde",
			AIClient: &aiClientMock{
				CallbackClose: func() error {
					return errClose
				},
			},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, service)
		err = service.Close()
		assert.Error(t, err)
		assert.ErrorIs(t, err, errClose)
	})
}
