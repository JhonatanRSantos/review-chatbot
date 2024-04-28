package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/JhonatanRSantos/review-chatbot/internal/datatypes"
)

type repository interface {
	Create(ctx context.Context, firstName string, lastName string, email string) (datatypes.User, error)
	FindByEmail(ctx context.Context, email string) (datatypes.User, error)
}

type UserService struct {
	repository repository
}

// NewUserService create a new user service
func NewUserService(repository repository) *UserService {
	return &UserService{repository}
}

// Create create a new user
func (us *UserService) Create(
	ctx context.Context,
	firstName string,
	lastName string,
	email string,
) (datatypes.User, error) {
	if sliceHasEmptyStrings(firstName, lastName, email) {
		return datatypes.User{}, fmt.Errorf("failed to create new user. Cause: %w", ErrMissingRequiredFields)
	}
	return us.repository.Create(ctx, firstName, lastName, email)
}

// FindByEmail finds a user by email
func (us *UserService) FindByEmail(ctx context.Context, email string) (datatypes.User, error) {
	if sliceHasEmptyStrings(email) {
		return datatypes.User{}, fmt.Errorf("failed to find user. Cause: %w", ErrMissingRequiredFields)
	}
	return us.repository.FindByEmail(ctx, email)
}

// sliceHasEmptyStrings
func sliceHasEmptyStrings(values ...string) bool {
	for _, value := range values {
		if strings.ReplaceAll(value, " ", "") == "" {
			return true
		}
	}
	return false
}
