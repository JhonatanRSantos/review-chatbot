package datatypes

import (
	"errors"
	"strings"
)

type User struct {
	ID        string `db:"id"         json:"id"`
	FirstName string `db:"first_name" json:"firstName"`
	LastName  string `db:"last_name"  json:"lastName"`
	Email     string `db:"email"      json:"email"`
}

type Chat struct {
	ID     string `db:"id" json:"-"`
	UserID string `db:"id" json:"-"`
}

type Message struct {
	ID      string `db:"id"      json:"id"`
	ChatID  string `db:"id"      json:"chatId"`
	Message string `db:"message" json:"message"`
	Author  string `db:"author"  json:"author"`
}

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (cur *CreateUserRequest) Validate() error {
	cur.FirstName = strings.ReplaceAll(cur.FirstName, " ", "")
	cur.LastName = strings.ReplaceAll(cur.LastName, " ", "")
	cur.Email = strings.ReplaceAll(cur.Email, " ", "")

	if cur.FirstName == "" || cur.LastName == "" || cur.Email == "" {
		return errors.New("missing required fields")
	}
	return nil
}

type CreateReviewRequest struct {
	User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
	Product string `json:"product"`
}

func (crr *CreateReviewRequest) Validate() error {
	crr.User.Name = strings.ReplaceAll(crr.User.Name, " ", "")
	crr.User.Email = strings.ReplaceAll(crr.User.Email, " ", "")
	crr.Product = strings.ReplaceAll(crr.Product, " ", "")

	if crr.User.Name == "" || crr.User.Email == "" || crr.Product == "" {
		return errors.New("missing required fields")
	}
	return nil
}
