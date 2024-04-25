package chatbot

import (
	"context"
	"io"

	"github.com/google/generative-ai-go/genai"
)

type aiClient interface {
	Close() error
	DeleteFile(ctx context.Context, name string) error
	EmbeddingModel(name string) *genai.EmbeddingModel
	GenerativeModel(name string) *genai.GenerativeModel
	GetFile(ctx context.Context, name string) (*genai.File, error)
	ListFiles(ctx context.Context) *genai.FileIterator
	ListModels(ctx context.Context) *genai.ModelInfoIterator
	UploadFile(ctx context.Context, name string, r io.Reader, opts *genai.UploadFileOptions) (*genai.File, error)
}
