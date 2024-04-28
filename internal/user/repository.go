package user

import (
	"context"
	"fmt"

	"github.com/JhonatanRSantos/gocore/pkg/godb"
	"github.com/JhonatanRSantos/review-chatbot/internal/datatypes"
	"github.com/gofrs/uuid/v5"
)

type Repository struct {
	db godb.DB
}

// NewRepository create a new repository
func NewRepository(db godb.DB) *Repository {
	return &Repository{db}
}

// Create create a new user
func (r *Repository) Create(ctx context.Context, firstName string, lastName string, email string) (datatypes.User, error) {
	stm, err := r.db.PrepareNamedContext(ctx, createUser)
	if err != nil {
		return datatypes.User{}, fmt.Errorf("failed to create new user. Cause: %w", err)
	}
	defer stm.Close()

	id, err := uuid.NewV4()
	if err != nil {
		return datatypes.User{}, fmt.Errorf("failed to create new user. Cause: %w", err)
	}

	params := map[string]interface{}{
		"id":         id.String(),
		"first_name": firstName,
		"last_name":  lastName,
		"email":      email,
	}

	result, err := stm.ExecContext(ctx, params)
	if err != nil {
		return datatypes.User{}, fmt.Errorf("failed to create new user. Cause: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return datatypes.User{}, fmt.Errorf("failed to create new user. Cause: %w", err)
	}

	if rows == 0 {
		return datatypes.User{}, ErrCantSaveUser
	}

	return datatypes.User{
		ID:        id.String(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}, nil
}

// FindByEmail finds a user by email
func (r *Repository) FindByEmail(ctx context.Context, email string) (datatypes.User, error) {
	var user datatypes.User

	stm, err := r.db.PrepareNamedContext(ctx, findUserByEmail)
	if err != nil {
		return datatypes.User{}, fmt.Errorf("failed to find user. Cause: %w", err)
	}
	defer stm.Close()

	params := map[string]interface{}{
		"email": email,
	}

	if err = stm.GetContext(ctx, &user, params); err != nil {
		return datatypes.User{}, fmt.Errorf("failed to find user. Cause: %w", err)
	}

	return user, nil
}
