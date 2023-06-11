package product_http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
	"market_auth/internal/product"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h product.HttpHandler) {
	r.Post("/", mw.ValidateAccessToken(), mw.ValidateAdminRole(), h.InsertProduct())
	r.Post("/manufacturer", mw.ValidateAccessToken(), mw.ValidateAdminRole(), h.InsertManufacturer())
	r.Post("/category", mw.ValidateAccessToken(), mw.ValidateAdminRole(), h.InsertCategory())

	r.Get("/category", h.FetchCategories())
	r.Get("/manufacturer", h.FetchManufacturers())
	r.Get("/", h.FetchProducts())
	r.Get("/:internal_id", h.GetProduct())

	r.Put("/:internal_id/like", mw.ValidateAccessToken(), h.LikeProduct())
	r.Put("/:internal_id/unlike", mw.ValidateAccessToken(), h.UnlikeProduct())

	r.Put("/:internal_id/show", mw.ValidateAccessToken(), mw.ValidateAdminRole(), h.ShowProduct())
	r.Put("/:internal_id/hide", mw.ValidateAccessToken(), mw.ValidateAdminRole(), h.HideProduct())
	r.Put("/:internal_id/count", mw.ValidateAccessToken(), mw.ValidateAdminRole(), h.UpdateProductCount())
}
