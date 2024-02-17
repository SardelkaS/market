package order_model

type FetchOrdersGatewayInput struct {
	UserId *int64
	Status *string
	Limit  *int64
	Offset *int64
}

type FetchOrderProductsGatewayInput struct {
	OrderId *int64
	Limit   *int64
	Offset  *int64
}
