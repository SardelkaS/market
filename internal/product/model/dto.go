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
	Category      *string
	Manufacturers []string
	MinPrice      *float64
	MaxPrice      *float64
	Show          *bool
	Like          *bool
	Sort          *string
	Limit         *int64
	Offset        *int64
}
