package handler

import (
	"fmt"
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
		fmt.Println(err)
		return err
	}

	var aa utils.ResFormat
	aa.Code = 200
	aa.Message = "test"
	aa.Data = users
	return utils.Send(c, aa)
}

func (h userHandler) GetUser(c *fiber.Ctx) error {
	userID := c.Params("userid")
	users, err := h.userSrv.GetUser(userID)
	if err != nil {
		e, ok := err.(utils.HandlerError)
		if ok {
			return e
		}

		// handleError(w, err)
		fmt.Println(err)
		return e
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"code":    200,
		"message": "",
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

	newBody := service.AddUserReq{
		Email:    body.Email,
		Password: body.Password,
		Name:     body.Name,
	}

	if err := h.userSrv.CreateUser(newBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"status":  false,
			"message": "Failed to parse body",
			"data":    "",
		})
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
