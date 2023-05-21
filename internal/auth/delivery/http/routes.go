package auth_http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
)

func MapRoutes(r fiber.Router, h auth.HttpHandler) {
	r.Post("/sign_in", h.SignIn())
	r.Post("/sign_up", h.SignUp())
	r.Post("/sign_out", h.ValidateAccessToken(), h.SignOut())

	r.Post("/refresh", h.ValidateRefreshToken(), h.GenerateRefresh())

	r.Patch("/password", h.ValidateAccessToken(), h.ChangePassword())
	r.Patch("/timezone", h.ValidateAccessToken(), h.ChangeTimezone())
}
