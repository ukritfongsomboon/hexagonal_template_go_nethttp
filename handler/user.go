package handler

import "github.com/gofiber/fiber/v2"

type UserHandler interface {
	GetUser(*fiber.Ctx) error
	GetUsers(*fiber.Ctx) error
	EditUser(*fiber.Ctx) error
	DeleteUser(*fiber.Ctx) error

	SignIn(*fiber.Ctx) error
	SignUp(*fiber.Ctx) error
}
