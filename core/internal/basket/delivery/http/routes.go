package basket_http

import (
	auth "core/internal/auth"
	"core/internal/basket"
	"github.com/gofiber/fiber/v2"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h basket.HttpHandler) {
	r.Post("/product", mw.VerifySignatureMiddleware(), h.AddProduct())
	r.Put("/incr", mw.VerifySignatureMiddleware(), h.IncrementCount())
	r.Put("/decr", mw.VerifySignatureMiddleware(), h.DecrementCount())
	r.Delete("/clear", mw.VerifySignatureMiddleware(), h.Clear())
	r.Get("/", mw.VerifySignatureMiddleware(), h.GetBasket())
}
