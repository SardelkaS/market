package order

import "github.com/gofiber/fiber/v2"

type HttpHandler interface {
	CreateOrder() fiber.Handler
	AttachProductToOrder() fiber.Handler
	RemoveProductFromOrder() fiber.Handler
	UpdateProductCount() fiber.Handler
	PendingOrder() fiber.Handler
	SendOrder() fiber.Handler
	DeliveryOrder() fiber.Handler
	CompleteOrder() fiber.Handler
	CancelOrder() fiber.Handler
	FetchOrders() fiber.Handler
	GetOrder() fiber.Handler
	FetchOrderProducts() fiber.Handler
}
