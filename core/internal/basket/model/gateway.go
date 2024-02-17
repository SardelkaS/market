package basket_model

type AddProductGatewayInput struct {
	UserId    *int64
	ProductId *int64
	Count     *int64
}
