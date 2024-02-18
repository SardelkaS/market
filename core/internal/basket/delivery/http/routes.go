package basket_http

import (
	"github.com/gofiber/fiber/v2"
	auth "market_auth/internal/auth"
	"market_auth/internal/basket"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h basket.HttpHandler) {
	r.Post("/product", mw.VerifySignatureMiddleware(), h.AddProduct())
	r.Put("/incr", mw.VerifySignatureMiddleware(), h.IncrementCount())
	r.Put("/decr", mw.VerifySignatureMiddleware(), h.DecrementCount())
	r.Delete("/clear", mw.VerifySignatureMiddleware(), h.Clear())
	r.Get("/", mw.VerifySignatureMiddleware(), h.GetBasket())
}
