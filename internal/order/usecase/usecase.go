package order_usecase

import (
	"fmt"
	"github.com/google/uuid"
	"market_auth/internal/order"
	order_model "market_auth/internal/order/model"
	"market_auth/internal/product"
	"market_auth/pkg/secure"
)

type uc struct {
	repo        order.Repository
	productRepo product.Repository
}

func New(repo order.Repository, productRepo product.Repository) order.UC {
	return &uc{
		repo:        repo,
		productRepo: productRepo,
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
		return nil, fmt.Errorf("error to create order")
	}

	orderData, err := u.repo.GetOrderById(*orderId)
	if err != nil {
		fmt.Printf("Error to get order %d: %s\n", *orderId, err.Error())
		return nil, fmt.Errorf("error to create order")
	}

	if len(input.Products) > 0 {
		err = u.AttachProductToOrder(order_model.AttachProductBody{
			OrderId:  orderData.InternalId,
			Products: input.Products,
		})
		if err != nil {
			return nil, err
		}
	}

	return orderData, nil
}

func (u *uc) AttachProductToOrder(input order_model.AttachProductBody) error {
	if input.OrderId == nil {
		return fmt.Errorf("wrong input")
	}
	orderData, err := u.repo.GetOrderByInternalId(*input.OrderId)
	if err != nil {
		fmt.Printf("Error to get order %s: %s\n", *input.OrderId, err.Error())
		return fmt.Errorf("order not found")
	}

	for _, inputProduct := range input.Products {
		if inputProduct.ProductId == nil {
			continue
		}
		productData, err := u.productRepo.GetProductByInternalId(*inputProduct.ProductId)
		if err != nil {
			fmt.Printf("Error to get product %s: %s\n", *inputProduct.ProductId, err.Error())
			return fmt.Errorf("product not found")
		}

		if inputProduct.Count == nil {
			count := int64(1)
			inputProduct.Count = &count
		}

		_, err = u.repo.AttachProductToOrder(*orderData.Id, *productData.Id, *inputProduct.Count)
		if err != nil {
			fmt.Printf("Error to attach produst %s to order %s: %s", *inputProduct.ProductId, *input.OrderId, err.Error())
			return fmt.Errorf("error to attach product")
		}
	}
	return nil
}

func (u *uc) RemoveProductFromOrder(input order_model.RemoveProductFromOrderBody) error {
	if input.OrderId == nil || input.ProductId == nil || input.UserId == nil {
		return fmt.Errorf("wrong input")
	}
	orderData, err := u.repo.GetOrderByInternalId(*input.OrderId)
	if err != nil {
		fmt.Printf("Error to get order %s: %s\n", *input.OrderId, err.Error())
		return fmt.Errorf("order not found")
	}

	if *orderData.UserId != *input.UserId {
		return fmt.Errorf("order belongs to anither user")
	}

	productData, err := u.productRepo.GetProductByInternalId(*input.ProductId)
	if err != nil {
		fmt.Printf("Error to get product %s: %s\n", *input.ProductId, err.Error())
		return fmt.Errorf("product not found")
	}

	err = u.repo.RemoveProductFromOrder(*orderData.Id, *productData.Id)
	if err != nil {
		fmt.Printf("Error to remove product %s from order %s: %s\n", *input.ProductId, *input.OrderId, err.Error())
		return fmt.Errorf("error to remove product from order")
	}
	return nil
}

func (u *uc) UpdateProductCount(input order_model.UpdateProductsCountBody) error {
	if input.OrderId == nil || input.ProductId == nil || input.UserId == nil || input.Count == nil {
		return fmt.Errorf("wrong input")
	}
	orderData, err := u.repo.GetOrderByInternalId(*input.OrderId)
	if err != nil {
		fmt.Printf("Error to get order %s: %s\n", *input.OrderId, err.Error())
		return fmt.Errorf("order not found")
	}

	if *orderData.UserId != *input.UserId {
		return fmt.Errorf("order belongs to anither user")
	}

	productData, err := u.productRepo.GetProductByInternalId(*input.ProductId)
	if err != nil {
		fmt.Printf("Error to get product %s: %s\n", *input.ProductId, err.Error())
		return fmt.Errorf("product not found")
	}

	if *input.Count <= 0 {
		err = u.repo.RemoveProductFromOrder(*orderData.Id, *productData.Id)
		if err != nil {
			fmt.Printf("Error to remove product %s from order %s: %s\n", *input.ProductId, *input.OrderId, err.Error())
			return fmt.Errorf("error to remove product from order")
		}
		return nil
	}

	err = u.repo.UpdateProductsCount(*orderData.Id, *productData.Id, *input.Count)
	if err != nil {
		fmt.Printf("Error to update product %s in order %s count: %s\n", *input.ProductId, *input.OrderId, err.Error())
		return fmt.Errorf("error to update product count")
	}

	return nil
}

func (u *uc) PendingOrder(orderId string, userId int64) error {
	orderData, err := u.repo.GetOrderByInternalId(orderId)
	if err != nil {
		fmt.Printf("Error to get order %s: %s\n", orderId, err.Error())
		return fmt.Errorf("order not found")
	}

	if *orderData.UserId != userId {
		return fmt.Errorf("order belongs to anither user")
	}

	err = u.repo.UpdateOrderStatus(orderId, order_model.PendingStatus)
	if err != nil {
		fmt.Printf("Error to update order %s status to pending: %s", orderId, err.Error())
		return fmt.Errorf("error to update order status")
	}
	return nil
}

func (u *uc) SendOrder(orderId string) error {
	err := u.repo.UpdateOrderStatus(orderId, order_model.SendedStatus)
	if err != nil {
		fmt.Printf("Error to update order %s status to sended: %s", orderId, err.Error())
		return fmt.Errorf("error to update order status")
	}
	return nil
}

func (u *uc) DeliveryOrder(orderId string) error {
	err := u.repo.UpdateOrderStatus(orderId, order_model.DeliveredStatus)
	if err != nil {
		fmt.Printf("Error to update order %s status to delivered: %s", orderId, err.Error())
		return fmt.Errorf("error to update order status")
	}
	return nil
}

func (u *uc) CompleteOrder(orderId string) error {
	err := u.repo.CompleteOrder(orderId)
	if err != nil {
		fmt.Printf("Error to complete order %s status: %s", orderId, err.Error())
		return fmt.Errorf("error to update order status")
	}
	return nil
}

func (u *uc) CancelOrder(orderId string) error {
	err := u.repo.CancelOrder(orderId)
	if err != nil {
		fmt.Printf("Error to cancel order %s status: %s", orderId, err.Error())
		return fmt.Errorf("error to update order status")
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
		return nil, fmt.Errorf("error to fetch orders")
	}
	count, err := u.repo.GetOrdersCount(params)
	if err != nil {
		fmt.Printf("Error to fetch orders count: %s", err.Error())
		return nil, fmt.Errorf("error to fetch orders")
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
		return nil, fmt.Errorf("order not found")
	}

	if *orderData.UserId != userId {
		return nil, fmt.Errorf("order belongs to anither user")
	}

	return orderData, nil
}

func (u *uc) FetchOrderProducts(input order_model.FetchOrderProductsParams) (*order_model.FetchOrderProductsResult, error) {
	orderData, err := u.repo.GetOrderByInternalId(*input.OrderId)
	if err != nil {
		fmt.Printf("Error to get order %s: %s\n", *input.OrderId, err.Error())
		return nil, fmt.Errorf("order not found")
	}

	if *orderData.UserId != *input.UserId {
		return nil, fmt.Errorf("order belongs to anither user")
	}

	params := order_model.FetchOrderProductsGatewayInput{
		OrderId: orderData.Id,
		Limit:   input.Limit,
		Offset:  input.Offset,
	}
	products, err := u.repo.FetchOrderProducts(params)
	if err != nil {
		fmt.Printf("Error to fetch order %s products: %s\n", *input.OrderId, err.Error())
		return nil, fmt.Errorf("error to fetch order products")
	}
	count, err := u.repo.GetOrderProductsCount(params)
	if err != nil {
		fmt.Printf("Error to fetch order %s products count: %s\n", *input.OrderId, err.Error())
		return nil, fmt.Errorf("error to fetch order products")
	}

	return &order_model.FetchOrderProductsResult{
		Products: products,
		Count:    count,
	}, nil
}
