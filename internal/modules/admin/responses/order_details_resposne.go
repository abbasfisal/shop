package responses

import "shop/internal/entities"

type OrderDetail struct {
	Order    Order
	Customer Customer
}

type Address struct {
}

func ToOrderDetail(orderEntity entities.Order, customerEntity entities.Customer) OrderDetail {
	return OrderDetail{
		Order:    ToOrder(orderEntity),
		Customer: ToCustomer(customerEntity),
	}
}
