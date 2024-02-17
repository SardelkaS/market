package product

import product_model "market_auth/internal/product/model"

type UC interface {
	FetchCategories() ([]product_model.CategoryInfo, error)
	FetchManufacturers() ([]string, error)
	FetchSexes() ([]string, error)
	FetchCountries() ([]string, error)
	FetchProducts(input product_model.FetchProductsInput) ([]product_model.Product, *int64, error)
	FindProducts(input product_model.FindProductsInput) ([]product_model.Product, *int64, error)
	GetProduct(internalId string) (*product_model.Product, error)

	LikeProduct(internalId string, userId int64) error
	UnlikeProduct(internalId string, userId int64) error

	FetchProductStars(internalId string) ([]product_model.ProductStars, error)

	UpdateProductCount(internalId string, count int64) error

	ViewProduct(userId int64, productInternalId string) error
	FetchRecentlyViewedProductsInfo(userId int64, limit int64) ([]product_model.ProductInfo, error)

	FetchBoughtProductsInfo(userId int64, limit *int64, offset *int64) ([]product_model.ProductInfo, *int64, error)

	GetProductsInfo(input []product_model.Product, userId *int64) ([]product_model.ProductInfo, error)
}
