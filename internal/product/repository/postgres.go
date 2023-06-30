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

func (p *postgres) InsertProduct(input product_model.Product) (*int64, error) {
	var id int64
	err := p.db.QueryRowx(`INSERT INTO product("name", internal_id, price, count, description, pictures, manufacturer_id, sex_id, country_id) 
		values($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`,
		input.Name, input.InternalId, input.Price, input.Count, input.Description, input.Pictures, input.ManufacturerId,
		input.SexId, input.CountryId).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (p *postgres) InsertCategory(name string) (*int64, error) {
	var id int64
	err := p.db.QueryRowx(`INSERT INTO "category"(name) values($1) returning id`, name).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (p *postgres) InsertManufacturer(name string) (*int64, error) {
	var id int64
	err := p.db.QueryRowx(`INSERT INTO "manufacturer"(name) values($1) returning id`, name).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (p *postgres) InsertProductCategory(productId int64, categoryId int64) (*int64, error) {
	var id int64
	err := p.db.QueryRowx(`INSERT INTO product_category(product_id, category_id) values($1, $2) returning id`, productId, categoryId).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
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

func (p *postgres) FetchCategories() ([]string, error) {
	var result []string
	err := p.db.Select(&result, `SELECT "name" FROM category`)
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
						where ($1 = any(select pc.category_id from product_category pc where pc.product_id = p.id) or $1 is null)
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
		}
	}
	query += ` LIMIT $10 OFFSET $11`
	err := p.db.Select(&result, query, input.Category, pq.Array(input.Manufacturers), input.MinPrice, input.MaxPrice,
		input.Show, input.UserId, input.Likes, pq.Array(input.Sexes), pq.Array(input.Countries), input.Limit, input.Offset)
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
						where ($1 = any(select pc.category_id from product_category pc where pc.product_id = p.id) or $1 is null)
							and (m.name = any($2) or $2 is null)
							and (p.price >= $3 or $3 is null)
							and (p.price <= $4 or $4 is null)
							and (p.show = $5 or $5 is null)`,
		input.Category, pq.Array(input.Manufacturers), input.MinPrice, input.MaxPrice, input.Show)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) ShowProduct(internalId string) error {
	_, err := p.db.Exec(`update product set "show" = true where internal_id = $1`, internalId)
	return err
}

func (p *postgres) HideProduct(internalId string) error {
	_, err := p.db.Exec(`update product set "show" = false where internal_id = $1`, internalId)
	return err
}

func (p *postgres) UpdateProductCount(internalId string, count int64) error {
	_, err := p.db.Exec(`update product set "count" = $2 where internal_id = $1`, internalId, count)
	return err
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
				    (select array_agg(c."name") from "category" c 
				                   left join product_category pc on c.id = pc.category_id
				                   		where pc.product_id = p.id) as categories,
				    p.description,
				    p.pictures,
				    p.buy_count,
				    p.show,
				    (select coalesce(avg(f.stars), 0) from feedback f where f.product_id = p.id)::int8 as stars,
				    (select count(*) from like_product where user_id = $2) > 0 as liked,
				    (select count(*) from feedback f where f.product_id = p.id) as feedbacks_count,
				    (select count(*) from basket b where b.user_id = $2 and b.product_id = p.id) > 0 as in_basket,
				    s.name as sex,
				    c.name as country
					from product p
						left join manufacturer m on p.manufacturer_id = m.id
						left join sex s on p.sex_id = s.id
						left join country c on p.country_id = c.id
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
