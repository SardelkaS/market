package order_usecase

import (
	"core/internal/failure"
	"core/internal/order"
	order_model "core/internal/order/model"
	"core/internal/product"
	"core/internal/tg_bot"
	tg_bot_model "core/internal/tg_bot/model"
	"core/pkg/secure"
	"fmt"
	"github.com/google/uuid"
)

type uc struct {
	repo        order.Repository
	productRepo product.Repository
	tgBot       tg_bot.UC
}

func New(repo order.Repository, productRepo product.Repository, tgBot tg_bot.UC) order.UC {
	return &uc{
		repo:        repo,
		productRepo: productRepo,
		tgBot:       tgBot,
	}
}

func (u *uc) CreateOrder(input order_model.CreateOrderBody) (*order_model.Order, error) {
	statusId := order_model.CreatedStatus
	internalId := secure.CalcInternalId(uuid.New().String())
	orderId, err := u.repo.CreateOrder(order_model.Order{
		InternalId:  &internalId,
		UserId:      input.UserId,
		StatusId:    &statusId,
		Address:     input.Address,
		ContactData: input.ContactData,
	})
	if err != nil {
		fmt.Printf("Error to create order: %s\n", err.Error())
		return nil, failure.ErrCreateOrder.Wrap(err)
	}

	orderData, err := u.repo.GetOrderById(*orderId)
	if err != nil {
		fmt.Printf("Error to get order %d: %s\n", *orderId, err.Error())
		return nil, failure.ErrCreateOrder.Wrap(err)
	}

	var productsNames []string
	if len(input.Products) > 0 {
		productsNames, err = u.AttachProductToOrder(order_model.AttachProductBody{
			OrderId:  orderData.InternalId,
			Products: input.Products,
		})
		if err != nil {
			return nil, err
		}
	}

	go u.tgBot.NotifyNewOrder(tg_bot_model.NotifyNewOrderLogicInput{
		InternalId:  *orderData.InternalId,
		Products:    productsNames,
		ContactData: *input.ContactData,
	})

	return orderData, nil
}

func (u *uc) AttachProductToOrder(input order_model.AttachProductBody) ([]string, error) {
	if input.OrderId == nil {
		return nil, fmt.Errorf("wrong input")
	}
	orderData, err := u.repo.GetOrderByInternalId(*input.OrderId)
	if err != nil {
		fmt.Printf("Error to get order %s: %s\n", *input.OrderId, err.Error())
		return nil, failure.ErrGetOrder.Wrap(err)
	}

	var productsNames []string
	for _, inputProduct := range input.Products {
		if inputProduct.ProductId == nil {
			continue
		}
		productData, err := u.productRepo.GetProductByInternalId(*inputProduct.ProductId)
		if err != nil {
			fmt.Printf("Error to get product %s: %s\n", *inputProduct.ProductId, err.Error())
			return nil, failure.ErrGetProduct.Wrap(err)
		}
		productsNames = append(productsNames, *productData.Name)

		if inputProduct.Count == nil {
			count := int64(1)
			inputProduct.Count = &count
		}

		_, err = u.repo.AttachProductToOrder(*orderData.Id, *productData.Id, *inputProduct.Count)
		if err != nil {
			fmt.Printf("Error to attach produst %s to order %s: %s", *inputProduct.ProductId, *input.OrderId, err.Error())
			return nil, failure.ErrAttachProduct.Wrap(err)
		}
	}
	return productsNames, nil
}

func (u *uc) RemoveProductFromOrder(input order_model.RemoveProductFromOrderBody) error {
	if input.OrderId == nil || input.ProductId == nil || input.UserId == nil {
		return fmt.Errorf("wrong input")
	}
	orderData, err := u.repo.GetOrderByInternalId(*input.OrderId)
	if err != nil {
		fmt.Printf("Error to get order %s: %s\n", *input.OrderId, err.Error())
		return failure.ErrGetOrder.Wrap(err)
	}

	if *orderData.UserId != *input.UserId {
		return fmt.Errorf("order belongs to anither user")
	}

	productData, err := u.productRepo.GetProductByInternalId(*input.ProductId)
	if err != nil {
		fmt.Printf("Error to get product %s: %s\n", *input.ProductId, err.Error())
		return failure.ErrGetProduct.Wrap(err)
	}

	err = u.repo.RemoveProductFromOrder(*orderData.Id, *productData.Id)
	if err != nil {
		fmt.Printf("Error to remove product %s from order %s: %s\n", *input.ProductId, *input.OrderId, err.Error())
		return failure.ErrDetachProduct.Wrap(err)
	}
	return nil
}

func (u *uc) UpdateProductCount(input order_model.UpdateProductsCountBody) error {
	if input.OrderId == nil || input.ProductId == nil || input.UserId == nil || input.Count == nil {
		return failure.ErrInput
	}
	orderData, err := u.repo.GetOrderByInternalId(*input.OrderId)
	if err != nil {
		fmt.Printf("Error to get order %s: %s\n", *input.OrderId, err.Error())
		return failure.ErrGetOrder.Wrap(err)
	}

	if *orderData.UserId != *input.UserId {
		return fmt.Errorf("order belongs to anither user")
	}

	productData, err := u.productRepo.GetProductByInternalId(*input.ProductId)
	if err != nil {
		fmt.Printf("Error to get product %s: %s\n", *input.ProductId, err.Error())
		return failure.ErrGetProduct.Wrap(err)
	}

	if *input.Count <= 0 {
		err = u.repo.RemoveProductFromOrder(*orderData.Id, *productData.Id)
		if err != nil {
			fmt.Printf("Error to remove product %s from order %s: %s\n", *input.ProductId, *input.OrderId, err.Error())
			return failure.ErrDetachProduct.Wrap(err)
		}
		return nil
	}

	err = u.repo.UpdateProductsCount(*orderData.Id, *productData.Id, *input.Count)
	if err != nil {
		fmt.Printf("Error to update product %s in order %s count: %s\n", *input.ProductId, *input.OrderId, err.Error())
		return failure.ErrUpdateOrderProductsCount.Wrap(err)
	}

	return nil
}

func (u *uc) PendingOrder(orderId string, userId int64) error {
	orderData, err := u.repo.GetOrderByInternalId(orderId)
	if err != nil {
		fmt.Printf("Error to get order %s: %s\n", orderId, err.Error())
		return failure.ErrGetOrder.Wrap(err)
	}

	if *orderData.UserId != userId {
		return fmt.Errorf("order belongs to anither user")
	}

	err = u.repo.UpdateOrderStatus(orderId, order_model.PendingStatus)
	if err != nil {
		fmt.Printf("Error to update order %s status to pending: %s", orderId, err.Error())
		return failure.ErrUpdateOrderStatus.Wrap(err)
	}
	return nil
}

func (u *uc) FetchOrders(input order_model.FetchOrdersParams) (*order_model.FetchOrdersResult, error) {
	params := order_model.FetchOrdersGatewayInput{
		UserId: input.UserId,
		Status: input.Status,
		Limit:  input.Limit,
		Offset: input.Offset,
	}
	orders, err := u.repo.FetchOrders(params)
	if err != nil {
		fmt.Printf("Error to fetch orders: %s", err.Error())
		return nil, failure.ErrGetOrder.Wrap(err)
	}
	count, err := u.repo.GetOrdersCount(params)
	if err != nil {
		fmt.Printf("Error to fetch orders count: %s", err.Error())
		return nil, failure.ErrGetOrder.Wrap(err)
	}

	return &order_model.FetchOrdersResult{
		Orders: orders,
		Count:  count,
	}, nil
}

func (u *uc) GetOrder(orderId string, userId int64) (*order_model.Order, error) {
	orderData, err := u.repo.GetOrderByInternalId(orderId)
	if err != nil {
		fmt.Printf("Error to get order %s: %s\n", orderId, err.Error())
		return nil, failure.ErrGetOrder.Wrap(err)
	}

	if *orderData.UserId != userId {
		return nil, failure.ErrForbidden
	}

	return orderData, nil
}

func (u *uc) FetchOrderProducts(input order_model.FetchOrderProductsParams) (*order_model.FetchOrderProductsResult, error) {
	orderData, err := u.repo.GetOrderByInternalId(*input.OrderId)
	if err != nil {
		fmt.Printf("Error to get order %s: %s\n", *input.OrderId, err.Error())
		return nil, failure.ErrGetOrder.Wrap(err)
	}

	if *orderData.UserId != *input.UserId {
		return nil, failure.ErrForbidden
	}

	params := order_model.FetchOrderProductsGatewayInput{
		OrderId: orderData.Id,
		Limit:   input.Limit,
		Offset:  input.Offset,
	}
	products, err := u.repo.FetchOrderProducts(params)
	if err != nil {
		fmt.Printf("Error to fetch order %s products: %s\n", *input.OrderId, err.Error())
		return nil, failure.ErrGetOrderProducts.Wrap(err)
	}
	count, err := u.repo.GetOrderProductsCount(params)
	if err != nil {
		fmt.Printf("Error to fetch order %s products count: %s\n", *input.OrderId, err.Error())
		return nil, failure.ErrGetOrderProducts.Wrap(err)
	}

	return &order_model.FetchOrderProductsResult{
		Products: products,
		Count:    count,
	}, nil
}

func (u *uc) GetOrdersInfo(orders []order_model.Order, productsLimit *int64) ([]order_model.OrderInfo, error) {
	var ids []int64
	for _, orderData := range orders {
		ids = append(ids, *orderData.Id)
	}

	ordersInfo, err := u.repo.GetOrdersInfo(ids)
	if err != nil {
		fmt.Printf("Error to get orders (%v) info: %s", ids, err.Error())
		return nil, failure.ErrGetOrder.Wrap(err)
	}

	if productsLimit != nil && len(ordersInfo) > 0 {
		for i := range ordersInfo {
			products, err := u.repo.FetchOrderProducts(order_model.FetchOrderProductsGatewayInput{
				OrderId: ordersInfo[i].Id,
				Limit:   productsLimit,
			})
			if err != nil {
				fmt.Printf("Error to fetch order %d products: %s\n", *ordersInfo[i].Id, err.Error())
				return nil, failure.ErrGetOrderProducts.Wrap(err)
			}

			var pIds []int64
			for _, productData := range products {
				pIds = append(pIds, *productData.Id)
			}

			productsInfo, err := u.productRepo.GetProductsInfo(pIds, ordersInfo[0].UserId)
			if err != nil {
				fmt.Printf("Error to fetch order %d products: %s\n", *ordersInfo[i].Id, err.Error())
				return nil, failure.ErrGetOrderProducts.Wrap(err)
			}

			ordersInfo[i].Products = productsInfo
		}
	}

	return ordersInfo, nil
}
