package handler

import "github.com/gofiber/fiber/v2"

type UserHandler interface {
	GetUser(*fiber.Ctx) error
	GetUsers(*fiber.Ctx) error
	CreateUser(*fiber.Ctx) error
	EditUser(*fiber.Ctx) error
	DeleteUser(*fiber.Ctx) error
}
