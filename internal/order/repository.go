package order

import (
	order_model "market_auth/internal/order/model"
	product_model "market_auth/internal/product/model"
)

type Repository interface {
	CreateOrder(input order_model.Order) (*int64, error)
	AttachProductToOrder(orderId int64, productId int64, count int64) (*int64, error)

	RemoveProductFromOrder(orderId int64, productId int64) error
	UpdateProductsCount(orderId int64, productId int64, count int64) error

	UpdateOrderStatus(orderId string, statusId int64) error
	CompleteOrder(orderId string) error
	CancelOrder(orderId string) error

	FetchOrders(input order_model.FetchOrdersGatewayInput) ([]order_model.Order, error)
	GetOrdersCount(input order_model.FetchOrdersGatewayInput) (*int64, error)
	GetOrderByInternalId(internalId string) (*order_model.Order, error)
	GetOrderById(id int64) (*order_model.Order, error)
	FetchOrderProducts(input order_model.FetchOrderProductsGatewayInput) ([]product_model.Product, error)
	GetOrderProductsCount(input order_model.FetchOrderProductsGatewayInput) (*int64, error)
	GetOrdersInfo(ids []int64) ([]order_model.OrderInfo, error)
}
