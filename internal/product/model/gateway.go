package product_model

type FetchProductsGatewayInput struct {
	Category      *int64
	Manufacturers []string
	MinPrice      *float64
	MaxPrice      *float64
	Show          *bool
	Sort          *string
	UserId        *int64
	Likes         *bool
	Sexes         []string
	Countries     []string
	Limit         *int64
	Offset        *int64
}
