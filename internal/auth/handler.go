package auth

import "github.com/gofiber/fiber/v2"

type HttpHandler interface {
	SignIn() fiber.Handler
	SignUp() fiber.Handler
	SignOut() fiber.Handler

	GenerateRefresh() fiber.Handler
	ChangePassword() fiber.Handler
	ChangeTimezone() fiber.Handler

	ValidateAccessToken() fiber.Handler
	ValidateRefreshToken() fiber.Handler
	NoMW() fiber.Handler
}
