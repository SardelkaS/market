package basket_model

type AddProductBody struct {
	ProductId string `json:"product_id"`
	Count     int64  `json:"count"`
}

type IncrementCountBody struct {
	ProductId string `json:"product_id"`
}

type DecrementCountBody struct {
	ProductId string `json:"product_id"`
}
