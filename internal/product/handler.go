package product

import "github.com/gofiber/fiber/v2"

type HttpHandler interface {
	InsertProduct() fiber.Handler
	InsertManufacturer() fiber.Handler
	InsertCategory() fiber.Handler
	FetchCategories() fiber.Handler
	FetchManufacturers() fiber.Handler
	FetchSexes() fiber.Handler
	FetchCountries() fiber.Handler
	FetchProducts() fiber.Handler
	GetProduct() fiber.Handler
	LikeProduct() fiber.Handler
	UnlikeProduct() fiber.Handler
	ShowProduct() fiber.Handler
	HideProduct() fiber.Handler
	UpdateProductCount() fiber.Handler
}
