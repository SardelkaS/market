package product_http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
	"market_auth/internal/product"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h product.HttpHandler) {
	r.Post("/")
	r.Post("/manufacturer")
	r.Post("/category")

	r.Get("/category", mw.ValidateAccessToken())
	r.Get("/manufacturer", mw.ValidateAccessToken())
	r.Get("/", mw.ValidateAccessToken())
	r.Get("/:internal_id", mw.ValidateAccessToken())

	r.Put("/:internal_id/like", mw.ValidateAccessToken())
	r.Put("/:internal_id/unlike", mw.ValidateAccessToken())

	r.Put("/:internal_id/show")
	r.Put("/:internal_id/hide")
	r.Put("/:internal_id/count")
}
