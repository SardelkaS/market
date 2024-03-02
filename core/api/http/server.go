package http

import (
	"core/api"
	http_error "core/api/http/error"
	api_http_model "core/api/http/model"
	"core/config"
	"core/internal"
	"core/internal/auth"
	auth_http "core/internal/auth/delivery/http"
	"core/internal/basket"
	basket_http "core/internal/basket/delivery/http"
	"core/internal/common"
	"core/internal/feedback"
	feedback_http "core/internal/feedback/delivery/http"
	"core/internal/order"
	order_http "core/internal/order/delivery/http"
	"core/internal/product"
	product_http "core/internal/product/delivery/http"
	util_http "core/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggerMDW "github.com/gofiber/fiber/v2/middleware/logger"
	recoverMDW "github.com/gofiber/fiber/v2/middleware/recover"

	goJson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type httpServer struct {
	fiber *fiber.App
	cfg   *config.Config
}

func NewHttpServer(cfg *config.Config) api.Server {
	return &httpServer{
		cfg: cfg,
	}
}

func (h *httpServer) Init() error {
	errorHandler := http_error.NewHandler()
	h.fiber = fiber.New(fiber.Config{
		Immutable:               true,
		AppName:                 "market",
		EnableTrustedProxyCheck: true,
		ErrorHandler:            errorHandler.Control,
		JSONEncoder:             goJson.Marshal,
		JSONDecoder:             goJson.Unmarshal,
	})

	h.fiber.Use(recoverMDW.New(recoverMDW.Config{
		EnableStackTrace:  true,
		StackTraceHandler: errorHandler.StackTrace,
	}))
	h.fiber.Use(loggerMDW.New(loggerMDW.Config{}))

	return nil
}

func (h *httpServer) MapHandlers(app *internal.App) error {
	// ENGINE
	h.fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	h.fiber.Get("/version", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(api_http_model.VersionResponse{
			Version: h.cfg.Service.Version,
			Response: common.Response{
				Status: common.SuccessStatus,
			},
		})
	})

	// UTILS
	reqReader := util_http.NewReader()

	// HANDLERS
	basketHandler := basket_http.NewHttpHandler(app.UC["basket"].(basket.UC), reqReader)
	orderHandler := order_http.NewHttpHandler(app.UC["order"].(order.UC), app.UC["product"].(product.UC), reqReader)
	productHandler := product_http.NewHttpHandler(app.UC["product"].(product.UC), reqReader)
	feedbackHandler := feedback_http.NewHttpHandler(app.UC["feedback"].(feedback.UC), reqReader)
	authHandler := auth_http.NewHttpHandler(app.UC["auth"].(auth.UC))

	basketGroup := h.fiber.Group("/basket")
	basket_http.MapRoutes(basketGroup, authHandler, basketHandler)

	orderGroup := h.fiber.Group("/order")
	order_http.MapRoutes(orderGroup, authHandler, orderHandler)

	productGroup := h.fiber.Group("/product")
	product_http.MapRoutes(productGroup, authHandler, productHandler)

	feedbackGroup := h.fiber.Group("/feedback")
	feedback_http.MapRoutes(feedbackGroup, authHandler, feedbackHandler)

	return nil
}

func (h *httpServer) Run() error {
	fmt.Printf("LISTENING %s:%s\n", h.cfg.Service.Host, h.cfg.Service.Port)
	err := h.fiber.Listen(h.cfg.Service.Host + ":" + h.cfg.Service.Port)
	return err
}
