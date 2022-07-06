package middleware

import "github.com/gofiber/fiber/v2/middleware/cors"

var Cors = cors.Config{
	Next:             nil,
	AllowOrigins:     "*",
	AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
	AllowHeaders:     "",
	AllowCredentials: false,
	ExposeHeaders:    "",
	MaxAge:           0,
}
