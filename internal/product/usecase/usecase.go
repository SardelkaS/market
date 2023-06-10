package product_usecase

import (
	"fmt"
	"github.com/google/uuid"
	"market_auth/internal/product"
	product_model "market_auth/internal/product/model"
	"market_auth/pkg/secure"
)

type uc struct {
	repo product.Repository
}

func New(repo product.Repository) product.UC {
	return &uc{
		repo: repo,
	}
}

func (u *uc) InsertProduct(input product_model.InsertProductBody) error {
	internalId := secure.CalcInternalId(uuid.New().String())
	manufacturerId, err := u.repo.GetManufacturerIdByName(*input.Manufacturer)
	if err != nil {
		fmt.Printf("Error to get manufacturer (%s): %s\n", *input.Manufacturer, err.Error())
		return fmt.Errorf("error to get manufacturer")
	}

	buyCount := int64(0)
	productId, err := u.repo.InsertProduct(product_model.Product{
		InternalId:     &internalId,
		Name:           input.Name,
		Price:          input.Price,
		Count:          input.Count,
		ManufacturerId: manufacturerId,
		Description:    input.Description,
		Pictures:       input.Pictures,
		BuyCount:       &buyCount,
		Show:           input.Show,
	})
	if err != nil {
		fmt.Printf("Error to insert product: %s\n", err.Error())
		return fmt.Errorf("error to insert product")
	}

	for _, category := range input.Categories {
		categoryId, err := u.repo.GetCategoryIdByName(category)
		if err != nil {
			fmt.Printf("Error to get category %s: %s\n", category, err.Error())
			continue
		}
		_, err = u.repo.InsertProductCategory(*productId, *categoryId)
		if err != nil {
			fmt.Printf("Error to insert product category: %s\n", err.Error())
		}
	}

	return nil
}

func (u *uc) InsertManufacturer(input product_model.InsertManufacturerBody) error {
	_, err := u.repo.InsertManufacturer(*input.Name)
	if err != nil {
		fmt.Printf("Error to save manufacturer: %s\n", err.Error())
		return fmt.Errorf("error to save manufacturer")
	}
	return nil
}

func (u *uc) InsertCategory(input product_model.InsertCategoryBody) error {
	_, err := u.repo.InsertCategory(*input.Name)
	if err != nil {
		fmt.Printf("Error to save category: %s\n", err.Error())
		return fmt.Errorf("error to save category")
	}
	return nil
}

func (u *uc) FetchCategories() ([]string, error) {
	categories, err := u.repo.FetchCategories()
	if err != nil {
		fmt.Printf("Error to fecth categories: %s\n", err.Error())
		return nil, fmt.Errorf("error to fecth categories")
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

func (u *uc) FetchProducts(input product_model.FetchProductsInput) ([]product_model.Product, *int64, error) {
	var categoryId *int64
	var err error
	if input.Category != nil {
		categoryId, err = u.repo.GetCategoryIdByName(*input.Category)
		if err != nil {
			fmt.Printf("Error to get category %s: %s\n", *input.Category, err.Error())
			return nil, nil, fmt.Errorf("error to fetch products")
		}
	}
	params := product_model.FetchProductsGatewayInput{
		Category:      categoryId,
		Manufacturers: input.Manufacturers,
		MinPrice:      input.MinPrice,
		MaxPrice:      input.MaxPrice,
		Show:          input.Show,
		Sort:          input.Sort,
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

	//TODO: CHECK LIKE

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

func (u *uc) ShowProduct(internalId string) error {
	err := u.repo.ShowProduct(internalId)
	if err != nil {
		fmt.Printf("Error to show product %s: %s", internalId, err.Error())
		return fmt.Errorf("error to show product")
	}
	return nil
}

func (u *uc) HideProduct(internalId string) error {
	err := u.repo.HideProduct(internalId)
	if err != nil {
		fmt.Printf("Error to hide product %s: %s", internalId, err.Error())
		return fmt.Errorf("error to hide product")
	}
	return nil
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
