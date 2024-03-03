package order_http

import (
	"core/internal/common"
	"core/internal/failure"
	"core/internal/order"
	order_model "core/internal/order/model"
	"core/internal/product"
	"core/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type httpHandler struct {
	uc        order.UC
	productUC product.UC
	reqUtil   *utils.Reader
}

func NewHttpHandler(uc order.UC, productUC product.UC, reqUtil *utils.Reader) order.HttpHandler {
	return &httpHandler{
		uc:        uc,
		productUC: productUC,
		reqUtil:   reqUtil,
	}
}

func (h httpHandler) CreateOrder() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body order_model.CreateOrderBody
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrGetUser
		}

		body.UserId = &userId
		rawResult, err := h.uc.CreateOrder(body)
		if err != nil {
			return err
		}
		if rawResult == nil {
			return fmt.Errorf("error to create order")
		}

		realResult, err := h.uc.GetOrdersInfo([]order_model.Order{*rawResult}, nil)
		if err != nil {
			return err
		}
		if len(realResult) == 0 {
			return fmt.Errorf("error to get order")
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: realResult[0],
		})
	}
}

func (h httpHandler) AttachProductToOrder() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body order_model.AttachProductBody
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrGetUser
		}

		orderId := ctx.Params("internal_id")
		if orderId == "" {
			return failure.ErrInput
		}

		body.UserId = &userId
		body.OrderId = &orderId
		_, err = h.uc.AttachProductToOrder(body)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) RemoveProductFromOrder() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body order_model.RemoveProductFromOrderBody
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrGetUser
		}

		orderId := ctx.Params("internal_id")
		if orderId == "" {
			return failure.ErrInput
		}

		body.UserId = &userId
		body.OrderId = &orderId
		err = h.uc.RemoveProductFromOrder(body)
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
		var body order_model.UpdateProductsCountBody
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrGetUser
		}

		orderId := ctx.Params("internal_id")
		if orderId == "" {
			return failure.ErrInput
		}

		body.UserId = &userId
		body.OrderId = &orderId
		err = h.uc.UpdateProductCount(body)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) PendingOrder() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrGetUser
		}

		orderId := ctx.Params("internal_id")
		if orderId == "" {
			return failure.ErrInput
		}

		err = h.uc.PendingOrder(orderId, userId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) FetchOrders() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params order_model.FetchOrdersParams
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrGetUser
		}
		params.UserId = &userId

		status := ctx.Query("status", "")
		if status != "" {
			params.Status = &status
		}
		limit, err := strconv.ParseInt(ctx.Query("limit", ""), 10, 64)
		if err == nil {
			params.Limit = &limit
		}

		offset, err := strconv.ParseInt(ctx.Query("offset", ""), 10, 64)
		if err == nil {
			params.Offset = &offset
		}

		var productsLimit *int64
		pl, err := strconv.ParseInt(ctx.Query("products_limit", ""), 10, 64)
		if err == nil {
			productsLimit = &pl
		}

		rawResult, err := h.uc.FetchOrders(params)
		if err != nil {
			return err
		}

		realResult, err := h.uc.GetOrdersInfo(rawResult.Orders, productsLimit)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: order_model.FetchOrdersResponse{
				Orders: realResult,
				Count:  rawResult.Count,
			},
		})
	}
}

func (h httpHandler) GetOrder() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrGetUser
		}

		orderId := ctx.Params("internal_id")
		if orderId == "" {
			return failure.ErrInput
		}

		rawResult, err := h.uc.GetOrder(orderId, userId)
		if err != nil {
			return err
		}
		if rawResult == nil {
			return fmt.Errorf("error to get order")
		}

		realResult, err := h.uc.GetOrdersInfo([]order_model.Order{*rawResult}, nil)
		if err != nil {
			return err
		}
		if len(realResult) == 0 {
			return fmt.Errorf("error to get order")
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: realResult[0],
		})
	}
}

func (h httpHandler) FetchOrderProducts() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params order_model.FetchOrderProductsParams
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrGetUser
		}
		params.UserId = &userId

		orderId := ctx.Params("internal_id")
		if orderId == "" {
			return failure.ErrInput
		}
		params.OrderId = &orderId

		limit, err := strconv.ParseInt(ctx.Query("limit", ""), 10, 64)
		if err == nil {
			params.Limit = &limit
		}
		offset, err := strconv.ParseInt(ctx.Query("offset", ""), 10, 64)
		if err == nil {
			params.Offset = &offset
		}

		rawResult, err := h.uc.FetchOrderProducts(params)
		if err != nil {
			return err
		}

		realResult, err := h.productUC.GetProductsInfo(rawResult.Products, &userId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: order_model.FetchOrderProductsResponse{
				Products: realResult,
				Count:    rawResult.Count,
			},
		})
	}
}
