package basket_http

import (
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/basket"
	basket_model "market_auth/internal/basket/model"
	"market_auth/internal/common"
	"market_auth/internal/failure"
	"market_auth/pkg/utils"
	"strconv"
)

type httpHandler struct {
	uc      basket.UC
	reqUtil *utils.Reader
}

func NewHttpHandler(uc basket.UC, reqUtil *utils.Reader) basket.HttpHandler {
	return httpHandler{
		uc:      uc,
		reqUtil: reqUtil,
	}
}

func (h httpHandler) AddProduct() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body basket_model.AddProductBody
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}

		err = h.uc.AddProduct(userId, body.ProductId, body.Count)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) IncrementCount() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body basket_model.IncrementCountBody
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}

		err = h.uc.IncrementCount(userId, body.ProductId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) DecrementCount() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body basket_model.DecrementCountBody
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}

		err = h.uc.DecrementCount(userId, body.ProductId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) Clear() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}

		err = h.uc.ClearBasket(userId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) GetBasket() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}

		rawResult, err := h.uc.GetBasket(userId)
		if err != nil {
			return err
		}

		realResult, err := h.uc.GetBasketInfo(rawResult, &userId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: realResult,
		})
	}
}
