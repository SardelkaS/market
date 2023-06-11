package feedback

import "github.com/gofiber/fiber/v2"

type HttpHandler interface {
	CreateFeedback() fiber.Handler
	RemoveFeedback() fiber.Handler
	GodRemoveFeedback() fiber.Handler
	GetFeedback() fiber.Handler
	FetchFeedback() fiber.Handler
	LikeFeedback() fiber.Handler
	UnlikeFeedback() fiber.Handler
}
