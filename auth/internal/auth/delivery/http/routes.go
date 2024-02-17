package auth_http

import (
	"auth/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func MapRoutes(r fiber.Router, h auth.HttpHandler) {
	r.Post("/sign_in", h.SignIn())
	r.Post("/sign_up", h.SignUp())
	r.Post("/sign_out", h.ValidateAccessToken(), h.SignOut())

	r.Post("/refresh", h.ValidateRefreshToken(), h.GenerateRefresh())

	r.Patch("/password", h.ValidateAccessToken(), h.ChangePassword())
	r.Patch("/timezone", h.ValidateAccessToken(), h.ChangeTimezone())
	r.Post("/user/info", h.ValidateAccessToken(), h.UpdateUserInfo())
	r.Get("/user/info", h.ValidateAccessToken(), h.GetUserInfo())

	r.Post("*", h.ValidateAccessToken(), h.Proxy())
	r.Get("*", h.ValidateAccessToken(), h.Proxy())
	r.Put("*", h.ValidateAccessToken(), h.Proxy())
	r.Patch("*", h.ValidateAccessToken(), h.Proxy())
	r.Delete("*", h.ValidateAccessToken(), h.Proxy())
}
