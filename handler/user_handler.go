package handler

import (
	"hexagonal/service"
	"hexagonal/utils"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userSrv service.UserService
}

func NewUserHandler(userSrv service.UserService) UserHandler {
	return userHandler{userSrv: userSrv}
}

func (h userHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userSrv.GetUsers()
	if err != nil {
		appErr, ok := err.(utils.HandlerError)
		if ok {
			return c.Status(appErr.Code).JSON(fiber.Map{
				"code":    appErr.Code,
				"status":  false,
				"message": appErr.Message,
				"data":    "",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    200,
		"status":  true,
		"message": "get user success",
		"data":    users,
	})
}

func (h userHandler) GetUser(c *fiber.Ctx) error {
	userID := c.Params("userid")
	users, err := h.userSrv.GetUser(userID)
	if err != nil {
		appErr, ok := err.(utils.HandlerError)
		if ok {
			return c.Status(appErr.Code).JSON(fiber.Map{
				"code":    appErr.Code,
				"status":  false,
				"message": appErr.Message,
				"data":    "",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    200,
		"status":  true,
		"message": "get user success",
		"data":    users,
	})
}

func (h userHandler) CreateUser(c *fiber.Ctx) error {
	body := new(service.AddUserReq)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"status":  false,
			"message": "Failed to parse body",
			"data":    "",
		})
	}

	if err := h.userSrv.CreateUser(*body); err != nil {
		appErr, ok := err.(utils.HandlerError)
		if ok {
			return c.Status(appErr.Code).JSON(fiber.Map{
				"code":    appErr.Code,
				"status":  false,
				"message": appErr.Message,
				"data":    "",
			})
		}
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"status":  true,
		"message": "create user success",
		"data":    "",
	})

}

func (h userHandler) EditUser(c *fiber.Ctx) error {
	return nil
}

func (h userHandler) DeleteUser(c *fiber.Ctx) error {
	return nil
}

func (h userHandler) Login(c *fiber.Ctx) error {
	body := new(service.AuthenReq)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"status":  false,
			"message": "Failed to parse body",
			"data":    "",
		})
	}

	data, err := h.userSrv.Authentication(body)
	if err != nil {
		appErr, ok := err.(utils.HandlerError)
		if ok {
			return c.Status(appErr.Code).JSON(fiber.Map{
				"code":    appErr.Code,
				"status":  false,
				"message": appErr.Message,
				"data":    "",
			})
		}
	}

	// # Success Case
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  true,
		"message": "login success",
		"data":    data,
	})

}
