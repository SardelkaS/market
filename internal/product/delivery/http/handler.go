package product_http

import (
	"market_auth/internal/product"
	"market_auth/pkg/logger"
)

type httpHandler struct {
	uc     product.UC
	logger logger.UC
}

func NewHttpHandler(uc product.UC, logger logger.UC) product.HttpHandler {
	return httpHandler{
		uc:     uc,
		logger: logger,
	}
}
