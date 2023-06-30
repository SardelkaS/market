package product

import product_model "market_auth/internal/product/model"

type UC interface {
	InsertProduct(input product_model.InsertProductBody) error
	InsertManufacturer(input product_model.InsertManufacturerBody) error
	InsertCategory(input product_model.InsertCategoryBody) error

	FetchCategories() ([]string, error)
	FetchManufacturers() ([]string, error)
	FetchSexes() ([]string, error)
	FetchCountries() ([]string, error)
	FetchProducts(input product_model.FetchProductsInput) ([]product_model.Product, *int64, error)
	GetProduct(internalId string) (*product_model.Product, error)

	LikeProduct(internalId string, userId int64) error
	UnlikeProduct(internalId string, userId int64) error

	ShowProduct(internalId string) error
	HideProduct(internalId string) error
	UpdateProductCount(internalId string, count int64) error

	GetProductsInfo(input []product_model.Product, userId *int64) ([]product_model.ProductInfo, error)
}
