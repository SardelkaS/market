package product_usecase

import (
	"core/internal/product"
	product_model "core/internal/product/model"
	"fmt"
)

type uc struct {
	repo product.Repository
}

func New(repo product.Repository) product.UC {
	return &uc{
		repo: repo,
	}
}

func (u *uc) FetchCategories() ([]product_model.CategoryInfo, error) {
	categories, err := u.repo.FetchCategories()
	if err != nil {
		fmt.Printf("Error to fecth categories: %s\n", err.Error())
		return nil, fmt.Errorf("error to fecth categories")
	}

	for i, category := range categories {
		subcategories, err := u.repo.FetchSubcategories(*category.Id)
		if err != nil {
			fmt.Printf("Error to fecth subcategories: %s\n", err.Error())
			return nil, fmt.Errorf("error to fecth subcategories")
		}
		categories[i].Subcategories = subcategories
	}
	return categories, nil
}

func (u *uc) FetchManufacturers() ([]string, error) {
	manufacturers, err := u.repo.FetchManufacturers()
	if err != nil {
		fmt.Printf("Error to fecth manufacturers: %s\n", err.Error())
		return nil, fmt.Errorf("error to fecth manufacturers")
	}
	return manufacturers, nil
}

func (u *uc) FetchSexes() ([]string, error) {
	manufacturers, err := u.repo.FetchSexes()
	if err != nil {
		fmt.Printf("Error to fecth sexes: %s\n", err.Error())
		return nil, fmt.Errorf("error to fecth sexes")
	}
	return manufacturers, nil
}

func (u *uc) FetchCountries() ([]string, error) {
	manufacturers, err := u.repo.FetchCountries()
	if err != nil {
		fmt.Printf("Error to fecth countries: %s\n", err.Error())
		return nil, fmt.Errorf("error to fecth countries")
	}
	return manufacturers, nil
}

func (u *uc) FetchProducts(input product_model.FetchProductsInput) ([]product_model.Product, *int64, error) {
	var subcategoryId *int64
	var err error
	if input.Subcategory != nil {
		subcategoryId, err = u.repo.GetSubcategoryIdByName(*input.Subcategory)
		if err != nil {
			fmt.Printf("Error to get category %s: %s\n", *input.Subcategory, err.Error())
			return nil, nil, fmt.Errorf("error to fetch products")
		}
	}
	params := product_model.FetchProductsGatewayInput{
		SubcategoryId: subcategoryId,
		Manufacturers: input.Manufacturers,
		MinPrice:      input.MinPrice,
		MaxPrice:      input.MaxPrice,
		Show:          input.Show,
		Sort:          input.Sort,
		UserId:        input.UserId,
		Liked:         input.Like,
		Sexes:         input.Sexes,
		Countries:     input.Countries,
		Limit:         input.Limit,
		Offset:        input.Offset,
	}
	products, err := u.repo.FetchProducts(params)
	if err != nil {
		fmt.Printf("Error to fetch products: %s\n", err.Error())
		return nil, nil, fmt.Errorf("error to fetch products")
	}

	count, err := u.repo.GetProductsCount(params)
	if err != nil {
		fmt.Printf("Error to get products count: %s\n", err.Error())
		return nil, nil, fmt.Errorf("error to fetch products")
	}

	return products, count, nil
}

func (u *uc) FindProducts(input product_model.FindProductsInput) ([]product_model.Product, *int64, error) {
	products, err := u.repo.FindProducts(input.NameTail, input.Limit, input.Offset)
	if err != nil {
		fmt.Printf("Error to find products: %s\n", err.Error())
		return nil, nil, fmt.Errorf("error to find products")
	}
	count, err := u.repo.FindProductsCount(input.NameTail)
	if err != nil {
		fmt.Printf("Error to find products count: %s\n", err.Error())
		return nil, nil, fmt.Errorf("error to find products")
	}

	return products, count, nil
}

func (u *uc) GetProduct(internalId string) (*product_model.Product, error) {
	productData, err := u.repo.GetProductByInternalId(internalId)
	if err != nil {
		fmt.Printf("Error to get product: %s\n", err.Error())
		return nil, fmt.Errorf("error to get product")
	}
	return productData, nil
}

func (u *uc) LikeProduct(internalId string, userId int64) error {
	productData, err := u.GetProduct(internalId)
	if err != nil {
		return err
	}

	liked, err := u.repo.CheckLiked(*productData.Id, userId)
	if err != nil {
		fmt.Printf("Error to check is product %d is liked by %d: %s", *productData.Id, userId, err.Error())
		return fmt.Errorf("error to check is product liked")
	}
	if *liked {
		return fmt.Errorf("product already liked")
	}

	err = u.repo.LikeProduct(*productData.Id, userId)
	if err != nil {
		fmt.Printf("Error to like product %s for user %d: %s", internalId, userId, err.Error())
		return fmt.Errorf("error to like product")
	}

	return nil
}

func (u *uc) UnlikeProduct(internalId string, userId int64) error {
	productData, err := u.GetProduct(internalId)
	if err != nil {
		return err
	}

	err = u.repo.UnlikeProduct(*productData.Id, userId)
	if err != nil {
		fmt.Printf("Error to unlike product %s for user %d: %s", internalId, userId, err.Error())
		return fmt.Errorf("error to unlike product")
	}

	return nil
}

func (u *uc) FetchProductStars(internalId string) ([]product_model.ProductStars, error) {
	productData, err := u.GetProduct(internalId)
	if err != nil {
		return nil, err
	}

	result, err := u.repo.FetchProductStars(*productData.Id)
	if err != nil {
		fmt.Printf("Error to fetch product %s stars: %s", internalId, err.Error())
		return nil, fmt.Errorf("error to fetch product stars")
	}

	return result, nil
}

func (u *uc) UpdateProductCount(internalId string, count int64) error {
	if count < 0 {
		return fmt.Errorf("wrong count")
	}
	err := u.repo.UpdateProductCount(internalId, count)
	if err != nil {
		fmt.Printf("Error to update product %s count: %s", internalId, err.Error())
		return fmt.Errorf("error to update product count")
	}
	return nil
}

func (u *uc) ViewProduct(userId int64, productInternalId string) error {
	productData, err := u.GetProduct(productInternalId)
	if err != nil {
		return err
	}

	err = u.repo.ViewProduct(userId, *productData.Id)
	if err != nil {
		fmt.Printf("Error to view product %s by user %d: %s", productInternalId, userId, err.Error())
		return fmt.Errorf("error to view product")
	}
	return nil
}

func (u *uc) FetchRecentlyViewedProductsInfo(userId int64, limit int64) ([]product_model.ProductInfo, error) {
	ids, err := u.repo.FetchRecentlyViewedIds(userId, limit)
	if err != nil {
		fmt.Printf("Error to fetch recently viewed by user %d products: %s", userId, err.Error())
		return nil, fmt.Errorf("error to fetch recently viewed products")
	}

	result, err := u.repo.GetProductsInfo(ids, &userId)
	if err != nil {
		fmt.Printf("Error to get products info: %s", err.Error())
		return nil, fmt.Errorf("error to get products info")
	}
	return result, nil
}

func (u *uc) FetchBoughtProductsInfo(userId int64, limit *int64, offset *int64) ([]product_model.ProductInfo, *int64, error) {
	ids, err := u.repo.FetchBoughtIds(userId, limit, offset)
	if err != nil {
		fmt.Printf("Error to fetch bought by user %d products: %s", userId, err.Error())
		return nil, nil, fmt.Errorf("error to fetch bought products")
	}

	count, err := u.repo.GetBoughtCount(userId)
	if err != nil {
		fmt.Printf("Error to get bought count by user %d products: %s", userId, err.Error())
		return nil, nil, fmt.Errorf("error to get bought products")
	}

	result, err := u.repo.GetProductsInfo(ids, &userId)
	if err != nil {
		fmt.Printf("Error to get products info: %s", err.Error())
		return nil, nil, fmt.Errorf("error to get products info")
	}
	return result, count, nil
}

func (u *uc) GetProductsInfo(input []product_model.Product, userId *int64) ([]product_model.ProductInfo, error) {
	var ids []int64
	for _, productData := range input {
		ids = append(ids, *productData.Id)
	}

	result, err := u.repo.GetProductsInfo(ids, userId)
	if err != nil {
		fmt.Printf("Error to get products info: %s", err.Error())
		return nil, fmt.Errorf("error to get products info")
	}
	return result, nil
}
