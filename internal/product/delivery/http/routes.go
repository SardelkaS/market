package product_http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
	"market_auth/internal/product"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h product.HttpHandler) {
	r.Post("/", h.InsertProduct())
	r.Post("/manufacturer", h.InsertManufacturer())
	r.Post("/category", h.InsertCategory())

	r.Get("/category", mw.ValidateAccessToken(), h.FetchCategories())
	r.Get("/manufacturer", mw.ValidateAccessToken(), h.FetchManufacturers())
	r.Get("/", mw.ValidateAccessToken(), h.FetchProducts())
	r.Get("/:internal_id", mw.ValidateAccessToken(), h.GetProduct())

	r.Put("/:internal_id/like", mw.ValidateAccessToken(), h.LikeProduct())
	r.Put("/:internal_id/unlike", mw.ValidateAccessToken(), h.UnlikeProduct())

	r.Put("/:internal_id/show", h.ShowProduct())
	r.Put("/:internal_id/hide", h.HideProduct())
	r.Put("/:internal_id/count", h.UpdateProductCount())
}
