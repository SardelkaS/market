package order

import order_model "market_auth/internal/order/model"

type UC interface {
	CreateOrder(input order_model.CreateOrderBody) (*order_model.Order, error)
	AttachProductToOrder(input order_model.AttachProductBody) error

	RemoveProductFromOrder(input order_model.RemoveProductFromOrderBody) error
	UpdateProductCount(input order_model.UpdateProductsCountBody) error

	PendingOrder(orderId string, userId int64) error
	FetchOrders(input order_model.FetchOrdersParams) (*order_model.FetchOrdersResult, error)
	GetOrder(orderId string, userId int64) (*order_model.Order, error)
	FetchOrderProducts(input order_model.FetchOrderProductsParams) (*order_model.FetchOrderProductsResult, error)

	GetOrdersInfo(orders []order_model.Order) ([]order_model.OrderInfo, error)
}
