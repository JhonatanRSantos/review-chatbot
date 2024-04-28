package chat

import "errors"

var (
	ErrCantCreateChat    = errors.New("failed to create new chat. Cause: error saving chat")
	ErrCantCreateMessage = errors.New("failed to create new message. Cause: error saving message")
)
