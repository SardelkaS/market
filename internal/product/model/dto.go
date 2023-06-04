package product_model

import "github.com/lib/pq"

type InsertProductBody struct {
	Name         *string        `json:"name"`
	Price        *float64       `json:"price"`
	Count        *int64         `json:"count"`
	Manufacturer *string        `json:"manufacturer"`
	Description  *string        `json:"description"`
	Categories   pq.StringArray `json:"categories"`
	Pictures     pq.StringArray `json:"pictures"`
	Show         *bool          `json:"show"`
}

type InsertManufacturerBody struct {
	Name *string `json:"name"`
}

type InsertCategoryBody struct {
	Name *string `json:"name"`
}

type FetchProductsInput struct {
	UserId        *int64   `query:"user_id"`
	Category      *string  `query:"category"`
	Manufacturers []string `query:"manufacturers"`
	MinPrice      *float64 `query:"min_price"`
	MaxPrice      *float64 `query:"max_price"`
	Show          *bool    `query:"show"`
	Like          *bool    `query:"like"`
	Sort          *string  `query:"sort"`
	Limit         *int64   `query:"limit"`
	Offset        *int64   `query:"offset"`
}

type FetchProductsResponse struct {
	Products []ProductInfo
	Count    *int64
}

type UpdateProductCountBody struct {
	Count int64 `json:"count"`
}
