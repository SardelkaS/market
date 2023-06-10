package auth_http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/auth"
	"market_auth/internal/auth/model"
	"market_auth/internal/common"
	"market_auth/internal/failure"
	"market_auth/pkg/logger"
	"market_auth/pkg/utils"
)

const _adminRole = "admin"

type httpHandler struct {
	uc     auth.UC
	logger logger.UC
}

func NewHttpHandler(uc auth.UC, logger logger.UC) auth.HttpHandler {
	return httpHandler{
		uc:     uc,
		logger: logger,
	}
}

func (h httpHandler) SignIn() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body auth_model.SignInBody
		if err := utils.ReadRequest(ctx, &body); err != nil {
			h.logger.Log(logger.Error, fmt.Sprintf("Read sign in request error: %s", err.Error()))
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		response, err := h.uc.SignIn(auth_model.SignInLogicInput{
			Login:    &body.Login,
			Password: &body.Password,
		})
		if err != nil {
			return err
		}
		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: response,
		})
	}
}

func (h httpHandler) SignUp() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body auth_model.SignUpBody
		if err := utils.ReadRequest(ctx, &body); err != nil {
			h.logger.Log(logger.Error, fmt.Sprintf("Read sign up request error: %s", err.Error()))
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		response, err := h.uc.SignUp(auth_model.SignUpLogicInput{
			Login:    body.Login,
			Password: body.Password,
		})
		if err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: response,
		})
	}
}

func (h httpHandler) SignOut() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body auth_model.SignOutBody
		body.FingerKey = ctx.Get("FingerKey")
		err := h.uc.SignOut(auth_model.SignOutLogicInput{FingerKey: body.FingerKey})
		if err != nil {
			return err
		}
		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) GenerateRefresh() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body auth_model.RefreshBody
		body.FingerKey = ctx.Get("FingerKey")
		body.RefreshToken = ctx.Get("RefreshToken")
		response, err := h.uc.GenerateRefresh(auth_model.GenerateRefreshLogicInput{
			FingerKey:    body.FingerKey,
			RefreshToken: body.RefreshToken,
		})
		if err != nil {
			return err
		}
		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: response,
		})
	}
}

func (h httpHandler) ChangePassword() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body auth_model.ChangePasswordBody
		if err := utils.ReadRequest(ctx, &body); err != nil {
			h.logger.Log(logger.Error, fmt.Sprintf("change password request error: %v", err))
			return failure.ErrInput
		}

		err := h.uc.ChangePassword(auth_model.ChangePasswordLogicInput{
			FingerKey:   ctx.Get("FingerKey"),
			OldPassword: body.OldPassword,
			NewPassword: body.NewPassword,
		})
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) ChangeTimezone() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body auth_model.ChangeTimezoneBody
		if err := utils.ReadRequest(ctx, &body); err != nil {
			h.logger.Log(logger.Error, fmt.Sprintf("change timezone request error: %v", err))
			return failure.ErrInput
		}

		err := h.uc.ChangeTimezone(auth_model.ChangeTimezoneLogicInput{
			FingerKey:   ctx.Get("FingerKey"),
			NewTimezone: body.NewTimezone,
		})
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) ValidateAdminRole() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		role := ctx.Get("role")
		if role != _adminRole {
			return ctx.Status(fiber.StatusForbidden).JSON(common.Response{
				Status:      common.FailedStatus,
				Description: "permission denied",
			})
		}

		return ctx.Next()
	}
}

func (h httpHandler) ValidateAccessToken() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params auth_model.ValidateBody
		params.Token = ctx.Get("AccessToken", "")
		params.FingerKey = ctx.Get("FingerKey", "")
		result, err := h.uc.ValidateJWT(auth_model.ValidateJWTLogicInput{
			Token:     params.Token,
			FingerKey: params.FingerKey,
		})
		if err != nil {
			return failure.ErrAuth
		}
		ctx.Request().Header.Set("role", *result.Role)
		ctx.Request().Header.Set("user_id", fmt.Sprint(*result.UserId))
		return ctx.Next()
	}
}

func (h httpHandler) ValidateRefreshToken() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params auth_model.ValidateBody
		params.Token = ctx.Get("RefreshToken", "")
		params.FingerKey = ctx.Get("FingerKey", "")
		result, err := h.uc.ValidateJWT(auth_model.ValidateJWTLogicInput{
			Token:     params.Token,
			FingerKey: params.FingerKey,
		})
		if err != nil {
			return failure.ErrAuth
		}
		ctx.Request().Header.Set("role", *result.Role)
		ctx.Request().Header.Set("user_id", fmt.Sprint(*result.UserId))
		return ctx.Next()
	}
}

func (h httpHandler) NoMW() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
