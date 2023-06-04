package product_model

import "github.com/lib/pq"

type Product struct {
	Id             *int64         `json:"id" db:"id"`
	InternalId     *string        `json:"internal_id" db:"internal_id"`
	Name           *string        `json:"name" db:"name"`
	Price          *float64       `json:"price" db:"price"`
	Count          *int64         `json:"count" db:"count"`
	ManufacturerId *int64         `json:"manufacturer_id" db:"manufacturer_id"`
	Description    *string        `json:"description" db:"description"`
	Pictures       pq.StringArray `json:"pictures" db:"pictures"`
	BuyCount       *int64         `json:"buy_count" db:"buy_count"`
	Show           *bool          `json:"show" db:"show"`
}

type ProductCategory struct {
	Id   *int64  `json:"id" db:"id"`
	Name *string `json:"name" db:"name"`
}

type ProductInfo struct {
	InternalId   *string        `json:"internal_id" db:"internal_id"`
	Name         *string        `json:"name" db:"name"`
	Price        *float64       `json:"price" db:"price"`
	Count        *int64         `json:"count" db:"count"`
	Manufacturer *string        `json:"manufacturer" db:"manufacturer"`
	Categories   pq.StringArray `json:"categories" db:"categories"`
	Description  *string        `json:"description" db:"description"`
	Pictures     pq.StringArray `json:"pictures" db:"pictures"`
	BuyCount     *int64         `json:"buy_count" db:"buy_count"`
	Show         *bool          `json:"show" db:"show"`
	Stars        *int64         `json:"stars" db:"stars"`
	Liked        *bool          `json:"liked" db:"liked"`
}
