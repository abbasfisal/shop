package responses

import "shop/internal/entities"

type OrderDetail struct {
	Order    AdminOrder
	Customer Customer
}

func ToOrderDetail(orderEntity entities.Order, customerEntity entities.Customer) OrderDetail {
	return OrderDetail{
		Order:    ToAdminOrder(orderEntity),
		Customer: ToCustomer(customerEntity),
	}
}
