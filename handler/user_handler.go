package handler

import (
	"hexagonal/repository"
	"hexagonal/service"
	"hexagonal/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userSrv service.UserService
}

func NewUserHandler(userSrv service.UserService) UserHandler {
	return userHandler{userSrv: userSrv}
}

func (h userHandler) GetUsers(c *fiber.Ctx) error {
	var p repository.PaginationUser
	p.Page, _ = strconv.Atoi(c.Query("page", "1"))
	p.Row, _ = strconv.Atoi(c.Query("row", "10"))

	users, err := h.userSrv.GetUsers(p)
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

func (h userHandler) SignUp(c *fiber.Ctx) error {
	body := new(service.AddUserReq)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"status":  false,
			"message": "Failed to parse body",
			"data":    "",
		})
	}

	// if err := h.userSrv.CreateUser(*body); err != nil {
	// 	appErr, ok := err.(utils.HandlerError)
	// 	if ok {
	// 		return c.Status(appErr.Code).JSON(fiber.Map{
	// 			"code":    appErr.Code,
	// 			"status":  false,
	// 			"message": appErr.Message,
	// 			"data":    "",
	// 		})
	// 	}
	// }

	user, err := h.userSrv.CreateUser(*body)
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

	// DTO
	// userRes := service.UserResponse{

	// }

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"status":  true,
		"message": "create user success",
		"data":    user,
	})

}

func (h userHandler) EditUser(c *fiber.Ctx) error {
	return nil
}

func (h userHandler) DeleteUser(c *fiber.Ctx) error {
	return nil
}

func (h userHandler) SignIn(c *fiber.Ctx) error {
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

	// clear cookie client
	c.ClearCookie()

	// create cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "Accesstoken"
	cookie.Value = data.Accesstoken
	cookie.Secure = false
	cookie.SessionOnly = true
	cookie.MaxAge = 3000
	cookie.Expires = time.Now().Add(10 * time.Second)

	// set cookie
	c.Cookie(cookie)

	// # Success Case
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  true,
		"message": "login success",
		"data":    data,
	})

}
