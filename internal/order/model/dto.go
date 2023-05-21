package order_model

type CreateOrderBody struct {
	UserId      *int64         `json:"user_id"`
	Address     *string        `json:"address"`
	ContactData *string        `json:"contact_data"`
	Products    []ProductsBody `json:"products"`
}

type ProductsBody struct {
	ProductId *string `json:"product_id"`
	Count     *int64  `json:"count"`
}

type AttachProductBody struct {
	OrderId  *string `json:"-"`
	Products []ProductsBody
}

type FetchOrdersParams struct {
	UserId *int64
	Status *string
	Limit  *int64
	Offset *int64
}

type FetchOrdersResult struct {
	Orders []Order
	Count  *int64
}

type FetchOrderProductsParams struct {
	OrderId *string
	UserId  *int64
	Limit   *int64
	Offset  *int64
}

type FetchOrderProductsResult struct {
	Products []OrderProduct
	Count    *int64
}

type RemoveProductFromOrderBody struct {
	OrderId   *string `json:"-"`
	UserId    *int64  `json:"-"`
	ProductId *string `json:"product_id"`
}

type UpdateProductsCountBody struct {
	OrderId   *string `json:"-"`
	UserId    *int64  `json:"-"`
	ProductId *string `json:"product_id"`
	Count     *int64  `json:"count"`
}
