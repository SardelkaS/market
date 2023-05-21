package order_http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
	"market_auth/internal/order"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h order.HttpHandler) {
	r.Post("/", mw.ValidateAccessToken())
	r.Put("/:internal_id/attach", mw.ValidateAccessToken())
	r.Put("/:internal_id/detach", mw.ValidateAccessToken())
	r.Put("/:internal_id/count", mw.ValidateAccessToken())

	r.Put("/:internal_id/pending")
	r.Put("/:internal_id/send")
	r.Put("/:internal_id/delivery")
	r.Put("/:internal_id/complete")
	r.Put("/:internal_id/cancel")

	r.Get("/", mw.ValidateAccessToken())
	r.Get("/:internal_id/", mw.ValidateAccessToken())
	r.Get("/:internal_id/products", mw.ValidateAccessToken())
}
