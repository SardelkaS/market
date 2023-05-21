package basket

import basket_model "market_auth/internal/basket/model"

type UC interface {
	AddProduct(userId int64, productId string, count int64) error

	IncrementCount(userId int64, productId string) error
	DecrementCount(userId int64, productId string) error
	ClearBasket(userId int64) error

	GetBasket(userId int64) ([]basket_model.Basket, error)
	GetBasketInfo(rawData []basket_model.Basket) (*basket_model.BasketInfo, error)
}
