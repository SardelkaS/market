package order_repository

import (
	"github.com/jmoiron/sqlx"
	"market_auth/internal/order"
	order_model "market_auth/internal/order/model"
)

type postgres struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) order.Repository {
	return &postgres{
		db: db,
	}
}

func (p *postgres) CreateOrder(input order_model.Order) (*int64, error) {
	var id int64
	err := p.db.QueryRowx(`insert into "order"(user_id, status_id, address, contact_data, create_time, internal_id) 
			values($1, $2, $3, $4, now(), $5) returning id`,
		input.UserId, input.StatusId, input.Address, input.ContactData, input.InternalId).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (p *postgres) AttachProductToOrder(orderId int64, productId int64, count int64) (*int64, error) {
	var id int64
	err := p.db.QueryRowx(`insert into order_products(order_id, product_id, count) values($1, $2, $3) returning id`,
		orderId, productId, count).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (p *postgres) RemoveProductFromOrder(orderId int64, productId int64) error {
	_, err := p.db.Exec(`delete from order_products where order_id = $1 and product_id = $2`, orderId, productId)
	return err
}

func (p *postgres) UpdateProductsCount(orderId int64, productId int64, count int64) error {
	_, err := p.db.Exec(`update order_products set count = $1 where product_id = $2 and order_id = $3`, count, productId, orderId)
	return err
}

func (p *postgres) UpdateOrderStatus(orderId string, statusId int64) error {
	_, err := p.db.Exec(`update "order" set status_id = $1, update_time = now() where internal_id = $2`, statusId, orderId)
	return err
}

func (p *postgres) CompleteOrder(orderId string) error {
	_, err := p.db.Exec(`update "order" set status_id = $1, complete_time = now() where internal_id = $2`, order_model.CompletedStatus, orderId)
	return err
}

func (p *postgres) CancelOrder(orderId string) error {
	_, err := p.db.Exec(`update "order" set status_id = $1, complete_time = now() where internal_id = $2`, order_model.CancelledStatus, orderId)
	return err
}

func (p *postgres) FetchOrders(input order_model.FetchOrdersGatewayInput) ([]order_model.Order, error) {
	var result []order_model.Order
	err := p.db.Select(&result, `select * from "order"
         								left join order_status os on "order".status_id = os.id
										where (user_id = $1 or $1 is null)
											and (os.name = $2 or $2 is null)
										limit $3 offset $4`,
		input.UserId, input.Status, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []order_model.Order{}, nil
	}
	return result, nil
}

func (p *postgres) GetOrdersCount(input order_model.FetchOrdersGatewayInput) (*int64, error) {
	var result int64
	err := p.db.Get(&result, `select count(*) from "order"
                						left join order_status os on "order".status_id = os.id
										where (user_id = $1 or $1 is null)
											and (os.name = $2 or $2 is null)`,
		input.UserId, input.Status)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) GetOrderByInternalId(internalId string) (*order_model.Order, error) {
	var result order_model.Order
	err := p.db.Get(&result, `select * from "order" where internal_id = $1`, internalId)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) GetOrderById(id int64) (*order_model.Order, error) {
	var result order_model.Order
	err := p.db.Get(&result, `select * from "order" where id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) FetchOrderProducts(input order_model.FetchOrderProductsGatewayInput) ([]order_model.OrderProduct, error) {
	var result []order_model.OrderProduct
	err := p.db.Select(&result, `select * from order_products where order_id = $1 limit $2 offset $3`, input.OrderId, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []order_model.OrderProduct{}, nil
	}
	return result, nil
}

func (p *postgres) GetOrderProductsCount(input order_model.FetchOrderProductsGatewayInput) (*int64, error) {
	var result int64
	err := p.db.Get(&result, `select count(*) from order_products where order_id = $1`, input.OrderId)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *postgres) GetOrdersInfo(ids []int64) ([]order_model.OrderInfo, error) {
	var result []order_model.OrderInfo
	err := p.db.Select(&result, `
			select 
			    o.internal_id,
			    o.user_id,
			    os.name as status,
			    o.address,
			    o.contact_data,
			    o.create_time,
			    o.update_time,
			    o.complete_time
			    	from "order" o
					left join order_status os on o.status_id = os.id
						where o.id = any($1)`, ids)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []order_model.OrderInfo{}, nil
	}
	return result, nil
}
