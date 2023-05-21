package basket_http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
	"market_auth/internal/basket"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h basket.HttpHandler) {
	r.Post("/product", mw.ValidateAccessToken())
	r.Put("/incr", mw.ValidateAccessToken())
	r.Put("/decr", mw.ValidateAccessToken())
	r.Delete("/clear", mw.ValidateAccessToken())
	r.Get("/", mw.ValidateAccessToken())
}
