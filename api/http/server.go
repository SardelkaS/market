package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggerMDW "github.com/gofiber/fiber/v2/middleware/logger"
	recoverMDW "github.com/gofiber/fiber/v2/middleware/recover"
	"market_auth/internal/basket"
	basket_http "market_auth/internal/basket/delivery/http"
	"market_auth/internal/order"
	order_http "market_auth/internal/order/delivery/http"
	"market_auth/internal/product"
	product_http "market_auth/internal/product/delivery/http"
	"market_auth/pkg/logger"

	"market_auth/api"
	http_error "market_auth/api/http/error"
	api_http_model "market_auth/api/http/model"
	"market_auth/config"
	"market_auth/internal"
	"market_auth/internal/auth"
	authHttp "market_auth/internal/auth/delivery/http"
	"market_auth/internal/common"
	util_http "market_auth/pkg/utils"

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
	authHandler := authHttp.NewHttpHandler(app.UC["auth"].(auth.UC), app.UC["logger"].(logger.UC))
	basketHandler := basket_http.NewHttpHandler(app.UC["basket"].(basket.UC), reqReader)
	orderHandler := order_http.NewHttpHandler(app.UC["order"].(order.UC), app.UC["product"].(product.UC), reqReader)
	productHandler := product_http.NewHttpHandler(app.UC["product"].(product.UC), reqReader)

	basketGroup := h.fiber.Group("/basket")
	basket_http.MapRoutes(basketGroup, authHandler, basketHandler)

	orderGroup := h.fiber.Group("/order")
	order_http.MapRoutes(orderGroup, authHandler, orderHandler)

	productGroup := h.fiber.Group("/product")
	product_http.MapRoutes(productGroup, authHandler, productHandler)

	return nil
}

func (h *httpServer) Run() error {
	fmt.Printf("LISTENING %s:%s\n", h.cfg.Service.Host, h.cfg.Service.Port)
	err := h.fiber.Listen(h.cfg.Service.Host + ":" + h.cfg.Service.Port)
	return err
}
