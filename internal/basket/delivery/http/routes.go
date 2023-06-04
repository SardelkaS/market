package basket_http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
	"market_auth/internal/basket"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h basket.HttpHandler) {
	r.Post("/product", mw.ValidateAccessToken(), h.AddProduct())
	r.Put("/incr", mw.ValidateAccessToken(), h.IncrementCount())
	r.Put("/decr", mw.ValidateAccessToken(), h.DecrementCount())
	r.Delete("/clear", mw.ValidateAccessToken(), h.Clear())
	r.Get("/", mw.ValidateAccessToken(), h.GetBasket())
}
