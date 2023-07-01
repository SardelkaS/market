package feedback_http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
	"market_auth/internal/feedback"
)

func MapRoutes(r fiber.Router, mw auth.HttpHandler, h feedback.HttpHandler) {
	r.Post("/", mw.ValidateAccessToken(), h.CreateFeedback())
	r.Delete("/:internal_id", mw.ValidateAccessToken(), h.RemoveFeedback())
	r.Get("/:internal_id", mw.ValidateAccessToken(), h.GetFeedback())
	r.Get("/", mw.ValidateAccessToken(), h.FetchFeedback())
	r.Put("/:internal_id/like", mw.ValidateAccessToken(), h.LikeFeedback())
	r.Put("/:internal_id/unlike", mw.ValidateAccessToken(), h.UnlikeFeedback())
}
