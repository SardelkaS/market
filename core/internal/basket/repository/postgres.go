package basket_repository

import (
	"github.com/jmoiron/sqlx"
	"market_auth/internal/basket"
	basket_model "market_auth/internal/basket/model"
)

type postgres struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) basket.Repository {
	return &postgres{
		db: db,
	}
}

func (p *postgres) AddProduct(input basket_model.AddProductGatewayInput) (*int64, error) {
	var id int64
	err := p.db.QueryRowx(`INSERT INTO basket(user_id, product_id, count) values($1, $2, $3) returning id`,
		input.UserId, input.ProductId, input.Count).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (p *postgres) CheckRecordExists(userId int64, productId int64) (*bool, error) {
	var result []bool
	err := p.db.Select(&result, `select count(*) > 0 from basket where user_id = $1 and product_id = $2`, userId, productId)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return &result[0], nil
}

func (p *postgres) IncrementCount(userId int64, productId int64) error {
	_, err := p.db.Exec(`update basket set count = count + 1 where user_id = $1 and product_id = $2`, userId, productId)
	return err
}

func (p *postgres) DecrementCount(userId int64, productId int64) error {
	_, err := p.db.Exec(`update basket set count = count - 1 where user_id = $1 and product_id = $2`, userId, productId)
	return err
}

func (p *postgres) DeleteProduct(userId int64, productId int64) error {
	_, err := p.db.Exec(`update basket set count = 0 where user_id = $1 and product_id = $2`, userId, productId)
	return err
}

func (p *postgres) ClearBasket(userId int64) error {
	_, err := p.db.Exec(`update basket set count = 0 where user_id = $1`, userId)
	return err
}

func (p *postgres) GetBasket(userId int64) ([]basket_model.Basket, error) {
	var result []basket_model.Basket
	err := p.db.Select(&result, `select * from basket where user_id = $1 and count > 0`, userId)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []basket_model.Basket{}, nil
	}
	return result, nil
}
