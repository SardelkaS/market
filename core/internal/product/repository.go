package product

import product_model "core/internal/product/model"

type Repository interface {
	GetProductByInternalId(internalId string) (*product_model.Product, error)
	GetManufacturerIdByName(name string) (*int64, error)
	GetCategoryIdByName(name string) (*int64, error)
	GetSubcategoryIdByName(name string) (*int64, error)
	GetSexIdByName(name string) (*int64, error)
	GetCountryIdByName(name string) (*int64, error)
	FetchCategories() ([]product_model.CategoryInfo, error)
	FetchSubcategories(categoryId int64) ([]string, error)
	FetchManufacturers() ([]string, error)
	FetchSexes() ([]string, error)
	FetchCountries() ([]string, error)
	FetchProducts(input product_model.FetchProductsGatewayInput) ([]product_model.Product, error)
	FindProducts(nameTail *string, limit *int64, offset *int64) ([]product_model.Product, error)
	FindProductsCount(nameTail *string) (*int64, error)
	GetProductsCount(input product_model.FetchProductsGatewayInput) (*int64, error)

	UpdateProductCount(internalId string, count int64) error

	CheckLiked(productId int64, userId int64) (*bool, error)
	LikeProduct(productId int64, userId int64) error
	UnlikeProduct(productId int64, userId int64) error

	FetchProductStars(productId int64) ([]product_model.ProductStars, error)

	ViewProduct(userId int64, productId int64) error
	FetchRecentlyViewedIds(userId int64, limit int64) ([]int64, error)

	FetchBoughtIds(userId int64, limit *int64, offset *int64) ([]int64, error)
	GetBoughtCount(userId int64) (*int64, error)

	GetProductsInfo(ids []int64, userId *int64) ([]product_model.ProductInfo, error)
}
