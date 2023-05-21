package basket_http

import (
	"market_auth/internal/basket"
	"market_auth/pkg/logger"
)

type httpHandler struct {
	uc     basket.UC
	logger logger.UC
}

func NewHttpHandler(uc basket.UC, logger logger.UC) basket.HttpHandler {
	return httpHandler{
		uc:     uc,
		logger: logger,
	}
}
