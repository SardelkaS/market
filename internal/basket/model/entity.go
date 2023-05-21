package basket_model

import product_model "market_auth/internal/product/model"

type Basket struct {
	Id        *int64 `db:"id"`
	UserId    *int64 `db:"user_id"`
	ProductId *int64 `db:"product_id"`
	Count     *int64 `db:"count"`
}

type BasketProductInfo struct {
	Count   *int64
	Product *product_model.ProductInfo
}

type BasketInfo struct {
	UserId   *int64              `json:"user_id"`
	Products []BasketProductInfo `json:"products"`
}
