package error

import (
	"auth/internal/common"
	"auth/internal/failure"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Control(ctx *fiber.Ctx, err error) error {
	var lErr failure.LogicError
	ok := errors.As(err, &lErr)
	if ok {
		return ctx.Status(lErr.Code()).JSON(common.Response{
			Status:      common.FailedStatus,
			Description: lErr.Description(),
		})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(common.Response{
		Status: common.FailedStatus,
	})
}

func (h Handler) StackTrace(ctx *fiber.Ctx, e interface{}) {
	fmt.Printf("%s %s: %+v\n", ctx.Path(), ctx.Method(), e)
}
