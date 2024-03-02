package order_http

import (
	auth "core/internal/auth"
	"core/internal/order"
	"github.com/gofiber/fiber/v2"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h order.HttpHandler) {
	r.Post("/", mw.VerifySignatureMiddleware(), h.CreateOrder())
	r.Put("/:internal_id/attach", mw.VerifySignatureMiddleware(), h.AttachProductToOrder())
	r.Put("/:internal_id/detach", mw.VerifySignatureMiddleware(), h.RemoveProductFromOrder())
	r.Put("/:internal_id/count", mw.VerifySignatureMiddleware(), h.UpdateProductCount())

	r.Put("/:internal_id/pending", mw.VerifySignatureMiddleware(), h.PendingOrder())

	r.Get("/", mw.VerifySignatureMiddleware(), h.FetchOrders())
	r.Get("/:internal_id/", mw.VerifySignatureMiddleware(), h.GetOrder())
	r.Get("/:internal_id/products", mw.VerifySignatureMiddleware(), h.FetchOrderProducts())
}
