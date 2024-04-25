package router

import (
	"github.com/JhonatanRSantos/gocore/pkg/goweb"
	"github.com/gofiber/fiber/v2"
)

type handlers interface {
	CreateReview(*fiber.Ctx) error
	CreateUser(ctx *fiber.Ctx) error
	HandleWebsocketConnection() func(*fiber.Ctx) error
}

// NewWebRoutes
func NewWebRoutes(handlers handlers) []goweb.WebRoute {
	return []goweb.WebRoute{
		{
			Method:   "GET",
			Path:     "/api/ws/:email",
			Handlers: []func(c *fiber.Ctx) error{handlers.HandleWebsocketConnection()},
		},
		{
			Method:   "POST",
			Path:     "/api/user",
			Handlers: []func(c *fiber.Ctx) error{handlers.CreateUser},
		},
		{
			Method:   "POST",
			Path:     "/api/review",
			Handlers: []func(c *fiber.Ctx) error{handlers.CreateReview},
		},
	}
}
