package feedback_http

import (
	"core/internal/common"
	"core/internal/failure"
	"core/internal/feedback"
	feedback_model "core/internal/feedback/model"
	"core/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type httpHandler struct {
	uc      feedback.UC
	reqUtil *utils.Reader
}

func NewHttpHandler(uc feedback.UC, reqUtil *utils.Reader) feedback.HttpHandler {
	return &httpHandler{
		uc:      uc,
		reqUtil: reqUtil,
	}
}

func (h httpHandler) CreateFeedback() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body feedback_model.CreateFeedbackBody
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return failure.ErrInput
		}

		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}

		body.UserId = &userId
		err = h.uc.CreateFeedback(body)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) RemoveFeedback() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}

		feedbackId := ctx.Params("internal_id", "")
		if feedbackId == "" {
			return failure.ErrInput
		}

		err = h.uc.RemoveFeedback(feedback_model.RemoveFeedbackBody{
			UserId:     &userId,
			FeedbackId: &feedbackId,
		})
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) GetFeedback() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}

		feedbackId := ctx.Params("internal_id", "")
		if feedbackId == "" {
			return failure.ErrInput
		}

		rawResult, err := h.uc.GetFeedbackByInternalId(feedbackId)
		if err != nil {
			return err
		}

		realResult, err := h.uc.GetFeedbackInfo([]feedback_model.Feedback{*rawResult}, userId)
		if err != nil {
			return err
		}
		if len(realResult) == 0 {
			return failure.ErrNotFound
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: realResult[0],
		})
	}
}

func (h httpHandler) FetchFeedback() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params feedback_model.FetchFeedbackParams
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}
		params.UserId = &userId

		productId := ctx.Query("product_id", "")
		if productId != "" {
			params.ProductId = &productId
		}

		onlyMy, err := strconv.ParseBool(ctx.Query("only_my", ""))
		if err == nil {
			params.OnlyMy = &onlyMy
		}

		limit, err := strconv.ParseInt(ctx.Query("limit", ""), 10, 64)
		if err == nil {
			params.Limit = &limit
		}

		offset, err := strconv.ParseInt(ctx.Query("offset", ""), 10, 64)
		if err == nil {
			params.Offset = &offset
		}

		rawResult, err := h.uc.FetchFeedback(params)
		if err != nil {
			return err
		}

		realResult, err := h.uc.GetFeedbackInfo(rawResult.Feedback, userId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: feedback_model.FetchFeedbackResponse{
				Feedback: realResult,
				Count:    rawResult.Count,
			},
		})
	}
}

func (h httpHandler) LikeFeedback() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}

		feedbackId := ctx.Params("internal_id", "")
		if feedbackId == "" {
			return failure.ErrInput
		}

		err = h.uc.LikeFeedback(feedbackId, userId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h httpHandler) UnlikeFeedback() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := strconv.ParseInt(ctx.Get("user_id", ""), 10, 64)
		if err != nil {
			return failure.ErrToGetUser
		}

		feedbackId := ctx.Params("internal_id", "")
		if feedbackId == "" {
			return failure.ErrInput
		}

		err = h.uc.UnlikeFeedback(feedbackId, userId)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}
