package http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
	"market_auth/internal/failure"
)

type httpHandler struct {
	uc auth.UC
}

func NewHttpHandler(uc auth.UC) auth.HttpHandler {
	return &httpHandler{
		uc: uc,
	}
}

func (h *httpHandler) VerifySignatureMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ok, err := h.uc.VerifySignature(ctx.Get("Service"), ctx.Get("Signature"), string(ctx.Body()), ctx.Get("Timestamp"), ctx.Get("RequestId"))
		if err != nil {
			return err
		}

		if !ok {
			return failure.ErrAuth
		}

		return ctx.Next()
	}
}

func (h *httpHandler) NoMW() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
