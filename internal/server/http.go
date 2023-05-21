package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggerMDW "github.com/gofiber/fiber/v2/middleware/logger"
	recoverMDW "github.com/gofiber/fiber/v2/middleware/recover"
	"market_auth/config"
	"market_auth/internal"
	"market_auth/internal/auth"
	auth_http "market_auth/internal/auth/delivery/http"
	"market_auth/pkg/logger"
)

type httpServer struct {
	fiber *fiber.App
	cfg   *config.Config
}

func NewHttpServer(cfg *config.Config) Server {
	return &httpServer{
		cfg: cfg,
	}
}

func (h *httpServer) Init() error {
	h.fiber = fiber.New(fiber.Config{
		Immutable:               true,
		AppName:                 "RSPFilter",
		EnableTrustedProxyCheck: true,
	})
	return nil
}

func (h *httpServer) MapHandlers(app *internal.App) error {
	h.fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))
	h.fiber.Use(recoverMDW.New())
	h.fiber.Use(loggerMDW.New())

	authHandler := auth_http.NewHttpHandler(app.UC["auth"].(auth.UC), app.UC["logger"].(logger.UC))

	auth_http.MapRoutes(h.fiber, authHandler)

	return nil
}

func (h *httpServer) Run() error {
	fmt.Printf("LISTENING %s:%s\n", h.cfg.Service.Host, h.cfg.Service.Port)
	err := h.fiber.Listen(h.cfg.Service.Host + ":" + h.cfg.Service.Port)
	return err
}
