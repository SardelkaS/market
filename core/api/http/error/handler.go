package error

import (
	"core/internal/common"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Control(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(common.Response{
		Status:      common.FailedStatus,
		Description: err.Error(),
	})
}

func (h Handler) StackTrace(ctx *fiber.Ctx, e interface{}) {
	fmt.Printf("%s %s: %+v\n", ctx.Path(), ctx.Method(), e)
}
