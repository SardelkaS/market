package feedback_http

import (
	"github.com/gofiber/fiber/v2"
	auth "market_auth/internal/auth"
	"market_auth/internal/feedback"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h feedback.HttpHandler) {
	r.Post("/", mw.VerifySignatureMiddleware(), h.CreateFeedback())
	r.Delete("/:internal_id", mw.VerifySignatureMiddleware(), h.RemoveFeedback())
	r.Get("/:internal_id", mw.VerifySignatureMiddleware(), h.GetFeedback())
	r.Get("/", mw.VerifySignatureMiddleware(), h.FetchFeedback())
	r.Put("/:internal_id/like", mw.VerifySignatureMiddleware(), h.LikeFeedback())
	r.Put("/:internal_id/unlike", mw.VerifySignatureMiddleware(), h.UnlikeFeedback())
}
