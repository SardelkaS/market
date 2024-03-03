package basket_usecase

import (
	"core/internal/basket"
	basket_model "core/internal/basket/model"
	"core/internal/failure"
	"core/internal/product"
	"fmt"
)

type uc struct {
	repo        basket.Repository
	productRepo product.Repository
}

func New(repo basket.Repository, productRepo product.Repository) basket.UC {
	return &uc{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (u *uc) AddProduct(userId int64, productId string, count int64) error {
	productData, err := u.productRepo.GetProductByInternalId(productId)
	if err != nil {
		fmt.Printf("Error to get product %s: %s\n", productId, err.Error())
		return failure.ErrGetProduct.Wrap(err)
	}

	exists, err := u.repo.CheckRecordExists(userId, *productData.Id)
	if err != nil {
		fmt.Printf("Error to check record %d for user %d: %s", *productData.Id, userId, err.Error())
		return failure.ErrAddProduct.Wrap(err)
	}
	if exists != nil && *exists {
		err = u.repo.IncrementCount(userId, *productData.Id)
		if err != nil {
			fmt.Printf("Error to increment product %s count: %s\n", productId, err.Error())
			return failure.ErrAddProduct.Wrap(err)
		}
		return nil
	}

	_, err = u.repo.AddProduct(basket_model.AddProductGatewayInput{
		UserId:    &userId,
		ProductId: productData.Id,
		Count:     &count,
	})
	if err != nil {
		fmt.Printf("Error to add product %s to basket: %s\n", productId, err.Error())
		return failure.ErrAddProduct.Wrap(err)
	}
	return nil
}

func (u *uc) IncrementCount(userId int64, productId string) error {
	productData, err := u.productRepo.GetProductByInternalId(productId)
	if err != nil {
		fmt.Printf("Error to get product %s: %s\n", productId, err.Error())
		return failure.ErrGetProduct.Wrap(err)
	}

	err = u.repo.IncrementCount(userId, *productData.Id)
	if err != nil {
		fmt.Printf("Error to increment product %s count: %s\n", productId, err.Error())
		return failure.ErrIncrProduct.Wrap(err)
	}
	return nil
}

func (u *uc) DecrementCount(userId int64, productId string) error {
	productData, err := u.productRepo.GetProductByInternalId(productId)
	if err != nil {
		fmt.Printf("Error to get product %s: %s\n", productId, err.Error())
		return failure.ErrGetProduct.Wrap(err)
	}

	err = u.repo.DecrementCount(userId, *productData.Id)
	if err != nil {
		fmt.Printf("Error to decrement product %s count: %s\n", productId, err.Error())
		return failure.ErrDecrProduct.Wrap(err)
	}
	return nil
}

func (u *uc) ClearBasket(userId int64) error {
	err := u.repo.ClearBasket(userId)
	if err != nil {
		fmt.Printf("Error to clear user %d basket: %s\n", userId, err.Error())
		return failure.ErrClearBasket.Wrap(err)
	}
	return nil
}

func (u *uc) GetBasket(userId int64) ([]basket_model.Basket, error) {
	result, err := u.repo.GetBasket(userId)
	if err != nil {
		fmt.Printf("Error to get user %d basket: %s\n", userId, err.Error())
		return nil, failure.ErrGetBasket.Wrap(err)
	}
	return result, nil
}

func (u *uc) GetBasketInfo(rawData []basket_model.Basket, userId *int64) (*basket_model.BasketInfo, error) {
	if len(rawData) == 0 {
		return &basket_model.BasketInfo{
			Products: []basket_model.BasketProductInfo{},
		}, nil
	}

	var products []basket_model.BasketProductInfo
	for _, data := range rawData {
		productData, err := u.productRepo.GetProductsInfo([]int64{*data.ProductId}, userId)
		if err != nil {
			fmt.Printf("Error to get product %d data: %s\n", *data.ProductId, err.Error())
			return nil, failure.ErrGetBasket.Wrap(err)
		}
		products = append(products, basket_model.BasketProductInfo{
			Count:   data.Count,
			Product: &productData[0],
		})
	}

	return &basket_model.BasketInfo{
		Products: products,
	}, nil
}
