package http

import (
	"auth/api"
	api_http_model "auth/api/http/model"
	"auth/config"
	"auth/internal"
	"auth/internal/auth"
	auth_http "auth/internal/auth/delivery/http"
	common2 "auth/internal/common"
	"auth/pkg/logger"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggerMDW "github.com/gofiber/fiber/v2/middleware/logger"
	recoverMDW "github.com/gofiber/fiber/v2/middleware/recover"

	http_error "auth/api/http/error"
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
			Response: common2.Response{
				Status: common2.SuccessStatus,
			},
		})
	})

	// HANDLERS
	authHandler := auth_http.NewHttpHandler(app.UC["auth"].(auth.UC), app.UC["logger"].(logger.UC))
	auth_http.MapRoutes(h.fiber, authHandler)

	return nil
}

func (h *httpServer) Run() error {
	fmt.Printf("LISTENING %s:%s\n", h.cfg.Service.Host, h.cfg.Service.Port)
	err := h.fiber.Listen(h.cfg.Service.Host + ":" + h.cfg.Service.Port)
	return err
}
