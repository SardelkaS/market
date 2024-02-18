package auth

import "github.com/gofiber/fiber/v2"

type HttpHandler interface {
	VerifySignatureMiddleware() fiber.Handler
	NoMW() fiber.Handler
}
