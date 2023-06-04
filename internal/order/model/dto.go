package order_model

import product_model "market_auth/internal/product/model"

type CreateOrderBody struct {
	UserId      *int64         `json:"-"`
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
	UserId   *int64  `json:"-"`
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

type FetchOrdersResponse struct {
	Orders []OrderInfo `json:"orders"`
	Count  *int64      `json:"count"`
}

type FetchOrderProductsParams struct {
	OrderId *string
	UserId  *int64
	Limit   *int64
	Offset  *int64
}

type FetchOrderProductsResult struct {
	Products []product_model.Product
	Count    *int64
}

type FetchOrderProductsResponse struct {
	Products []product_model.ProductInfo
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
