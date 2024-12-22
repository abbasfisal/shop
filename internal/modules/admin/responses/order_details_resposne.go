package responses

import "shop/internal/entities"

type OrderDetail struct {
	Order    *AdminOrder
	Customer *Customer
}

func ToOrderDetail(orderEntity *entities.Order, customerEntity *entities.Customer) *OrderDetail {
	var orderDetail OrderDetail

	if orderEntity != nil {
		orderDetail.Order = ToAdminOrder(orderEntity)
	} else {
		orderDetail.Order = &AdminOrder{}
	}

	if customerEntity != nil {
		orderDetail.Customer = ToCustomer(customerEntity)
	} else {
		orderDetail.Customer = &Customer{}
	}

	return &orderDetail
}
