package product_model

type FetchProductsInput struct {
	UserId        *int64   `query:"user_id"`
	Subcategory   *string  `query:"subcategory"`
	Manufacturers []string `query:"manufacturers"`
	MinPrice      *float64 `query:"min_price"`
	MaxPrice      *float64 `query:"max_price"`
	Show          *bool    `query:"show"`
	Like          *bool    `query:"like"`
	Sort          *string  `query:"sort"`
	Sexes         []string `query:"sexes"`
	Countries     []string `query:"countries"`
	Limit         *int64   `query:"limit"`
	Offset        *int64   `query:"offset"`
}

type FetchProductsResponse struct {
	Products []ProductInfo
	Count    *int64
}
