package product_http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"market_auth/internal/common"
	"market_auth/internal/failure"
	"market_auth/internal/product"
	product_model "market_auth/internal/product/model"
	"market_auth/pkg/utils"
	"strconv"
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

func (h httpHandler) InsertProduct() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body product_model.InsertProductBody
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return failure.ErrInput
		}

		err = h.uc.InsertProduct(body)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) InsertManufacturer() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body product_model.InsertManufacturerBody
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return failure.ErrInput
		}

		err = h.uc.InsertManufacturer(body)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) InsertCategory() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body product_model.InsertCategoryBody
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return failure.ErrInput
		}

		err = h.uc.InsertCategory(body)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
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

func (h httpHandler) FetchProducts() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params product_model.FetchProductsInput
		err := h.reqUtil.Read(ctx.Context(), ctx.QueryParser, &params)
		if err != nil {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}
		params.UserId = &userId

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

func (h httpHandler) GetProduct() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		productId := ctx.Params("internal_id", "")
		if productId == "" {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}

		rawResult, err := h.uc.GetProduct(productId)
		if err != nil {
			return err
		}
		if rawResult == nil {
			return fmt.Errorf("error to get product")
		}

		realResult, err := h.uc.GetProductsInfo([]product_model.Product{*rawResult}, &userId)
		if err != nil {
			return err
		}
		if len(realResult) == 0 {
			return fmt.Errorf("error to get product info")
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: realResult,
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
			return failure.ErrToGetUser
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
			return failure.ErrToGetUser
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

func (h httpHandler) ShowProduct() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		productId := ctx.Params("internal_id", "")
		if productId == "" {
			return failure.ErrInput
		}

		err := h.uc.ShowProduct(productId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) HideProduct() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		productId := ctx.Params("internal_id", "")
		if productId == "" {
			return failure.ErrInput
		}

		err := h.uc.HideProduct(productId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) UpdateProductCount() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body product_model.UpdateProductCountBody
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return failure.ErrInput
		}

		productId := ctx.Params("internal_id", "")
		if productId == "" {
			return failure.ErrInput
		}

		err = h.uc.UpdateProductCount(productId, body.Count)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}
