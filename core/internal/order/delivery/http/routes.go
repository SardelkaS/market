package order_http

import (
	"auth/internal/auth"
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/order"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h order.HttpHandler) {
	r.Post("/", mw.ValidateAccessToken(), h.CreateOrder())
	r.Put("/:internal_id/attach", mw.ValidateAccessToken(), h.AttachProductToOrder())
	r.Put("/:internal_id/detach", mw.ValidateAccessToken(), h.RemoveProductFromOrder())
	r.Put("/:internal_id/count", mw.ValidateAccessToken(), h.UpdateProductCount())

	r.Put("/:internal_id/pending", mw.ValidateAccessToken(), h.PendingOrder())

	r.Get("/", mw.ValidateAccessToken(), h.FetchOrders())
	r.Get("/:internal_id/", mw.ValidateAccessToken(), h.GetOrder())
	r.Get("/:internal_id/products", mw.ValidateAccessToken(), h.FetchOrderProducts())
}
