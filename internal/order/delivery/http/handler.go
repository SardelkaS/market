package order_http

import "market_auth/internal/order"

type httpHandler struct {
	uc order.UC
}

func NewHttpHandler(uc order.UC) order.HttpHandler {
	return &httpHandler{
		uc: uc,
	}
}
