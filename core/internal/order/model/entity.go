package order_model

import (
	product_model "market_auth/internal/product/model"
	"time"
)

type Order struct {
	Id           *int64     `db:"id"`
	InternalId   *string    `db:"internal_id"`
	UserId       *int64     `db:"user_id"`
	StatusId     *int64     `db:"status_id"`
	Address      *string    `db:"address"`
	ContactData  *string    `db:"contact_data"`
	CreateTime   *time.Time `db:"create_time"`
	UpdateTime   *time.Time `db:"update_time"`
	CompleteTime *time.Time `db:"complete_time"`
}

type OrderProduct struct {
	Id        *int64 `db:"id"`
	OrderId   *int64 `db:"order_id"`
	ProductId *int64 `db:"product_id"`
	Count     *int64 `db:"count"`
}

type OrderProductInfo struct {
	Count   *int64                     `json:"count"`
	Product *product_model.ProductInfo `json:"product"`
}

type OrderInfo struct {
	Id            *int64                      `json:"-" db:"id"`
	InternalId    *string                     `json:"internal_id" db:"internal_id"`
	UserId        *int64                      `json:"user_id" db:"user_id"`
	Status        *string                     `json:"status" db:"status"`
	Address       *string                     `json:"address" db:"address"`
	Cost          *float64                    `json:"cost" db:"cost"`
	ProductsCount *int64                      `json:"products_count" db:"products_count"`
	ContactData   *string                     `json:"contact_data" db:"contact_data"`
	CreateTime    *time.Time                  `json:"create_time" db:"create_time"`
	UpdateTime    *time.Time                  `json:"update_time" db:"update_time"`
	CompleteTime  *time.Time                  `json:"complete_time" db:"complete_time"`
	Products      []product_model.ProductInfo `json:"products,omitempty" db:"-"`
}
