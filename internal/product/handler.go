package product

import "github.com/gofiber/fiber/v2"

type HttpHandler interface {
	FetchCategories() fiber.Handler
	FetchManufacturers() fiber.Handler
	FetchSexes() fiber.Handler
	FetchCountries() fiber.Handler
	FetchProducts() fiber.Handler
	GetProduct() fiber.Handler
	LikeProduct() fiber.Handler
	UnlikeProduct() fiber.Handler
	FetchRecentlyViewedProducts() fiber.Handler
}
