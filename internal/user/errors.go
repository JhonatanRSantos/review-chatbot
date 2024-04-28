package user

import "errors"

var (
	ErrCantSaveUser          = errors.New("failed to create new user. Cause: can't save the user")
	ErrMissingRequiredFields = errors.New("missing required fields")
)
