package basket

import "github.com/gofiber/fiber/v2"

type HttpHandler interface {
	AddProduct() fiber.Handler
	IncrementCount() fiber.Handler
	DecrementCount() fiber.Handler
	Clear() fiber.Handler
	GetBasket() fiber.Handler
}
