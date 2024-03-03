package product_http

import (
	"core/internal/common"
	"core/internal/failure"
	"core/internal/product"
	product_model "core/internal/product/model"
	"core/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

const (
	_defaultLimit = int64(5)
)

type httpHandler struct {
	uc      product.UC
	reqUtil *utils.Reader
}

func NewHttpHandler(uc product.UC, reqUtil *utils.Reader) product.HttpHandler {
	return httpHandler{
		uc:      uc,
		reqUtil: reqUtil,
	}
}

func (h httpHandler) FetchCategories() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		result, err := h.uc.FetchCategories()
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: result,
		})
	}
}

func (h httpHandler) FetchManufacturers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		result, err := h.uc.FetchManufacturers()
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: result,
		})
	}
}

func (h httpHandler) FetchSexes() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		result, err := h.uc.FetchSexes()
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: result,
		})
	}
}

func (h httpHandler) FetchCountries() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		result, err := h.uc.FetchCountries()
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: result,
		})
	}
}

func (h httpHandler) FetchProducts() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params product_model.FetchProductsInput
		err := h.reqUtil.Read(ctx.Context(), ctx.QueryParser, &params)
		if err != nil {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err == nil {
			params.UserId = &userId
		}

		rawResult, count, err := h.uc.FetchProducts(params)
		if err != nil {
			return err
		}

		realResult, err := h.uc.GetProductsInfo(rawResult, &userId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: product_model.FetchProductsResponse{
				Products: realResult,
				Count:    count,
			},
		})
	}
}

func (h httpHandler) FindProducts() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params product_model.FindProductsInput
		err := h.reqUtil.Read(ctx.Context(), ctx.QueryParser, &params)
		if err != nil {
			return failure.ErrInput
		}

		var userId *int64
		ui, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err == nil {
			userId = &ui
		}

		rawResult, count, err := h.uc.FindProducts(params)
		if err != nil {
			return err
		}

		realResult, err := h.uc.GetProductsInfo(rawResult, userId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: product_model.FetchProductsResponse{
				Products: realResult,
				Count:    count,
			},
		})
	}
}

func (h httpHandler) GetProduct() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		productId := ctx.Params("internal_id", "")
		if productId == "" {
			return failure.ErrInput
		}

		var paramsUserId *int64
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err == nil {
			paramsUserId = &userId
		}

		rawResult, err := h.uc.GetProduct(productId)
		if err != nil {
			return err
		}
		if rawResult == nil {
			return fmt.Errorf("error to get product")
		}

		realResult, err := h.uc.GetProductsInfo([]product_model.Product{*rawResult}, paramsUserId)
		if err != nil {
			return err
		}
		if len(realResult) == 0 {
			return fmt.Errorf("error to get product info")
		}

		if paramsUserId != nil {
			_ = h.uc.ViewProduct(*paramsUserId, productId)
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: realResult,
		})
	}
}

func (h httpHandler) FetchProductStars() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		productId := ctx.Params("internal_id", "")
		if productId == "" {
			return failure.ErrInput
		}

		result, err := h.uc.FetchProductStars(productId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: result,
		})
	}
}

func (h httpHandler) LikeProduct() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		productId := ctx.Params("internal_id", "")
		if productId == "" {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrGetUser
		}

		err = h.uc.LikeProduct(productId, userId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) UnlikeProduct() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		productId := ctx.Params("internal_id", "")
		if productId == "" {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrGetUser
		}

		err = h.uc.UnlikeProduct(productId, userId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) FetchRecentlyViewedProducts() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrGetUser
		}

		limit, err := strconv.ParseInt(ctx.Query("limit", ""), 10, 64)
		if err != nil {
			limit = _defaultLimit
		}

		result, err := h.uc.FetchRecentlyViewedProductsInfo(userId, limit)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: result,
		})
	}
}

func (h httpHandler) FetchBoughtProducts() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrGetUser
		}

		var paramLimit *int64
		limit, err := strconv.ParseInt(ctx.Query("limit", ""), 10, 64)
		if err == nil {
			paramLimit = &limit
		}

		var paramOffset *int64
		offset, err := strconv.ParseInt(ctx.Query("offset", ""), 10, 64)
		if err == nil {
			paramOffset = &offset
		}

		result, count, err := h.uc.FetchBoughtProductsInfo(userId, paramLimit, paramOffset)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: product_model.FetchProductsResponse{
				Products: result,
				Count:    count,
			},
		})
	}
}
