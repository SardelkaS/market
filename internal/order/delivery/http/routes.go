package order_http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
	"market_auth/internal/order"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h order.HttpHandler) {
	r.Post("/", mw.ValidateAccessToken(), h.CreateOrder())
	r.Put("/:internal_id/attach", mw.ValidateAccessToken(), h.AttachProductToOrder())
	r.Put("/:internal_id/detach", mw.ValidateAccessToken(), h.RemoveProductFromOrder())
	r.Put("/:internal_id/count", mw.ValidateAccessToken(), h.UpdateProductCount())

	r.Put("/:internal_id/pending", mw.ValidateAccessToken(), h.PendingOrder())
	r.Put("/:internal_id/send", mw.ValidateAccessToken(), mw.ValidateAdminRole(), h.SendOrder())
	r.Put("/:internal_id/delivery", mw.ValidateAccessToken(), mw.ValidateAdminRole(), h.DeliveryOrder())
	r.Put("/:internal_id/complete", mw.ValidateAccessToken(), mw.ValidateAdminRole(), h.CompleteOrder())
	r.Put("/:internal_id/cancel", mw.ValidateAccessToken(), mw.ValidateAdminRole(), h.CancelOrder())

	r.Get("/", mw.ValidateAccessToken(), h.FetchOrders())
	r.Get("/:internal_id/", mw.ValidateAccessToken(), h.GetOrder())
	r.Get("/:internal_id/products", mw.ValidateAccessToken(), h.FetchOrderProducts())
}
