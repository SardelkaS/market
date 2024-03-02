package basket

import basket_model "core/internal/basket/model"

type Repository interface {
	AddProduct(input basket_model.AddProductGatewayInput) (*int64, error)
	CheckRecordExists(userId int64, productId int64) (*bool, error)
	IncrementCount(userId int64, productId int64) error
	DecrementCount(userId int64, productId int64) error
	DeleteProduct(userId int64, productId int64) error
	ClearBasket(userId int64) error
	GetBasket(userId int64) ([]basket_model.Basket, error)
}
