package product_http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
	"market_auth/internal/product"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h product.HttpHandler) {
	r.Get("/category", mw.VerifySignatureMiddleware(), h.FetchCategories())
	r.Get("/manufacturer", mw.VerifySignatureMiddleware(), h.FetchManufacturers())
	r.Get("/sex", mw.VerifySignatureMiddleware(), h.FetchSexes())
	r.Get("/country", mw.VerifySignatureMiddleware(), h.FetchCountries())
	r.Get("/", mw.VerifySignatureMiddleware(), h.FetchProducts())
	r.Get("/find", mw.VerifySignatureMiddleware(), h.FindProducts())
	r.Get("/recently", mw.VerifySignatureMiddleware(), h.FetchRecentlyViewedProducts())
	r.Get("/bought", mw.VerifySignatureMiddleware(), h.FetchBoughtProducts())
	r.Get("/:internal_id", mw.VerifySignatureMiddleware(), h.GetProduct())
	r.Get("/:internal_id/stars", mw.VerifySignatureMiddleware(), h.FetchProductStars())

	r.Put("/:internal_id/like", mw.VerifySignatureMiddleware(), h.LikeProduct())
	r.Put("/:internal_id/unlike", mw.VerifySignatureMiddleware(), h.UnlikeProduct())
}
