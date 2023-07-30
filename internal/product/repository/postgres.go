package product_repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"market_auth/internal/product"
	product_model "market_auth/internal/product/model"
)

type postgres struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) product.Repository {
	return &postgres{
		db: db,
	}
}

func (p *postgres) GetProductByInternalId(internalId string) (*product_model.Product, error) {
	var result product_model.Product
	err := p.db.Get(&result, `SELECT * FROM product WHERE internal_id = $1`, internalId)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) GetManufacturerIdByName(name string) (*int64, error) {
	var result int64
	err := p.db.Get(&result, `SELECT id from manufacturer where "name" = $1`, name)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) GetCategoryIdByName(name string) (*int64, error) {
	var result int64
	err := p.db.Get(&result, `SELECT id from category where "name" = $1`, name)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) GetSubcategoryIdByName(name string) (*int64, error) {
	var result int64
	err := p.db.Get(&result, `SELECT id from subcategory where "name" = $1`, name)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) GetSexIdByName(name string) (*int64, error) {
	var result int64
	err := p.db.Get(&result, `SELECT id from sex where "name" = $1`, name)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) GetCountryIdByName(name string) (*int64, error) {
	var result int64
	err := p.db.Get(&result, `SELECT id from country where "name" = $1`, name)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) FetchCategories() ([]product_model.CategoryInfo, error) {
	var result []product_model.CategoryInfo
	err := p.db.Select(&result, `SELECT id, "name" FROM category`)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []product_model.CategoryInfo{}, nil
	}
	return result, nil
}

func (p *postgres) FetchSubcategories(categoryId int64) ([]string, error) {
	var result []string
	err := p.db.Select(&result, `SELECT "name" FROM subcategory where category_id = $1`, categoryId)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []string{}, nil
	}
	return result, nil
}

func (p *postgres) FetchManufacturers() ([]string, error) {
	var result []string
	err := p.db.Select(&result, `SELECT "name" FROM manufacturer`)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []string{}, nil
	}
	return result, nil
}

func (p *postgres) FetchSexes() ([]string, error) {
	var result []string
	err := p.db.Select(&result, `SELECT "name" FROM sex`)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []string{}, nil
	}
	return result, nil
}

func (p *postgres) FetchCountries() ([]string, error) {
	var result []string
	err := p.db.Select(&result, `SELECT "name" FROM country`)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []string{}, nil
	}
	return result, nil
}

func (p *postgres) FetchProducts(input product_model.FetchProductsGatewayInput) ([]product_model.Product, error) {
	var result []product_model.Product
	query := `SELECT p.* FROM product p
         				left join manufacturer m on p.manufacturer_id = m.id
           				left join sex s on p.sex_id = s.id
						left join country c on p.country_id = c.id
						where (p.subcategory_id = $1 or $1 is null)
							and (m.name = any($2) or $2 is null)
							and (p.price >= $3 or $3 is null)
							and (p.price <= $4 or $4 is null)
							and (p.show = $5 or $5 is null)
							and (p.id = any(select lp.product_id from like_product lp where lp.user_id = $6) 
							         or cast($7 as bool) is null or cast($7 as bool) = false)
							and (s.name = any($8) or $8 is null)
							and (c.name = any($9) or $9 is null)`
	if input.Sort != nil {
		if *input.Sort == "price_asc" {
			query += ` order by p.price asc`
		} else if *input.Sort == "price_desc" {
			query += ` order by p.price desc`
		} else if *input.Sort == "stars" {
			query += ` order by (select avg(stars) from feedback where product_id = p.id) desc`
		} else if *input.Sort == "popularity_asc" {
			query += ` order by p.buy_count asc`
		} else if *input.Sort == "popularity_desc" {
			query += ` order by p.buy_count desc`
		} else if *input.Sort == "discount" {
			query += ` order by (1 - p.price/coalesce(p.old_price, p.price)) desc`
		} else if *input.Sort == "newest" {
			query += ` order by p.create_date desc`
		}
	}
	query += ` LIMIT $10 OFFSET $11`
	err := p.db.Select(&result, query, input.SubcategoryId, pq.Array(input.Manufacturers), input.MinPrice, input.MaxPrice,
		input.Show, input.UserId, input.Liked, pq.Array(input.Sexes), pq.Array(input.Countries), input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []product_model.Product{}, nil
	}
	return result, nil
}

func (p *postgres) GetProductsCount(input product_model.FetchProductsGatewayInput) (*int64, error) {
	var result int64
	err := p.db.Get(&result, `SELECT COUNT(p.*) FROM product p
         				left join manufacturer m on p.manufacturer_id = m.id
                  		left join sex s on p.sex_id = s.id
						left join country c on p.country_id = c.id
						where (p.subcategory_id = $1 or $1 is null)
							and (m.name = any($2) or $2 is null)
							and (p.price >= $3 or $3 is null)
							and (p.price <= $4 or $4 is null)
							and (p.show = $5 or $5 is null)
							and (p.id = any(select lp.product_id from like_product lp where lp.user_id = $6) 
							         or cast($7 as bool) is null or cast($7 as bool) = false)
							and (s.name = any($8) or $8 is null)
							and (c.name = any($9) or $9 is null)`,
		input.SubcategoryId, pq.Array(input.Manufacturers), input.MinPrice, input.MaxPrice, input.Show,
		input.UserId, input.Liked, pq.Array(input.Sexes), pq.Array(input.Countries))
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) FindProducts(nameTail *string, limit *int64, offset *int64) ([]product_model.Product, error) {
	var result []product_model.Product
	err := p.db.Select(&result, `
select * from product p
	where p.name like '%' || $1 || '%'
	order by p.buy_count desc
	limit $2 offset $3`,
		nameTail, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []product_model.Product{}, nil
	}
	return result, nil
}

func (p *postgres) FindProductsCount(nameTail *string) (*int64, error) {
	var result []int64
	err := p.db.Select(&result, `
select count(*) from product p
	where p.name like '%' || $1 || '%'`,
		nameTail)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return &result[0], nil
}

func (p *postgres) UpdateProductCount(internalId string, count int64) error {
	_, err := p.db.Exec(`update product set "count" = $2 where internal_id = $1`, internalId, count)
	return err
}

func (p *postgres) CheckLiked(productId int64, userId int64) (*bool, error) {
	var result bool
	err := p.db.Get(&result, `select count(*) > 0 from like_product where user_id = $1 and product_id = $2`, userId, productId)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) LikeProduct(productId int64, userId int64) error {
	_, err := p.db.Exec(`insert into like_product(product_id, user_id) values($1, $2)`, productId, userId)
	return err
}

func (p *postgres) UnlikeProduct(productId int64, userId int64) error {
	_, err := p.db.Exec(`delete from like_product where product_id = $1 and user_id = $2`, productId, userId)
	return err
}

func (p *postgres) GetProductsInfo(ids []int64, userId *int64) ([]product_model.ProductInfo, error) {
	var result []product_model.ProductInfo
	err := p.db.Select(&result, `
				select 
				    p.internal_id,
				    p.name,
				    p.price,
				    p.old_price,
				    p.count,
				    m.name as manufacturer,
				    p.description,
				    p.pictures,
				    p.buy_count,
				    p.show,
				    (select coalesce(avg(f.stars), 0) from feedback f where f.product_id = p.id)::int8 as stars,
				    (select count(*) from like_product lp where lp.user_id = $2 and lp.product_id = p.id) > 0 as liked,
				    (select count(*) from feedback f where f.product_id = p.id) as feedbacks_count,
				    (select count(*) from basket b where b.user_id = $2 and b.product_id = p.id and b.count > 0) > 0 as in_basket,
				    s.name as sex,
				    c.name as country,
				    s2.name as subcategory
					from product p
						left join manufacturer m on p.manufacturer_id = m.id
						left join sex s on p.sex_id = s.id
						left join country c on p.country_id = c.id
						left join subcategory s2 on s2.id = p.subcategory_id
							where p.id = any($1)
						order by array_position($1, p.id)`, pq.Array(ids), userId)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []product_model.ProductInfo{}, nil
	}
	return result, nil
}

func (p *postgres) ViewProduct(userId int64, productId int64) error {
	_, err := p.db.Exec(`insert into recently_viewed(user_id, product_id, view_date) values($1, $2, now())`,
		userId, productId)
	return err
}

func (p *postgres) FetchRecentlyViewedIds(userId int64, limit int64) ([]int64, error) {
	var result []int64
	err := p.db.Select(&result, `select product_id from recently_viewed where user_id = $1 order by view_date desc limit $2`,
		userId, limit)
	if err != nil {
		return nil, err
	}
	return result, nil
}
