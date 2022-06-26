package handler

import (
	"fmt"
	"hexagonal/service"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userSrv service.UserService
}

func NewUserHandler(userSrv service.UserService) userHandler {
	return userHandler{userSrv: userSrv}
}

func (h userHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userSrv.GetUsers()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"code":    200,
		"message": "",
		"data":    users,
	})

}

func (h userHandler) GetUser(c *fiber.Ctx) error {
	userID := c.Params("userid")
	users, err := h.userSrv.GetUser(userID)
	if err != nil {
		// handleError(w, err)
		fmt.Println(err)
		return err
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"code":    200,
		"message": "",
		"data":    users,
	})
}
