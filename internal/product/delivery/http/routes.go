package product_http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
	"market_auth/internal/product"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h product.HttpHandler) {
	r.Get("/category", mw.SetUser(), h.FetchCategories())
	r.Get("/manufacturer", mw.SetUser(), h.FetchManufacturers())
	r.Get("/sex", mw.SetUser(), h.FetchSexes())
	r.Get("/country", mw.SetUser(), h.FetchCountries())
	r.Get("/", mw.SetUser(), h.FetchProducts())
	r.Get("/recently", mw.SetUser(), h.FetchRecentlyViewedProducts())
	r.Get("/:internal_id", mw.SetUser(), h.GetProduct())

	r.Put("/:internal_id/like", mw.ValidateAccessToken(), h.LikeProduct())
	r.Put("/:internal_id/unlike", mw.ValidateAccessToken(), h.UnlikeProduct())
}
