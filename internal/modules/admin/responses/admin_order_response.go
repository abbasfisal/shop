package responses

import (
	"shop/internal/entities"
)

type AdminOrder struct {
	CustomerID         uint
	OrderNumber        string
	PaymentStatus      uint
	TotalOriginalPrice uint
	TotalSalePrice     uint
	Discount           uint
	OrderStatus        uint
	OrderStatusText    string
	OrderItems         AdminOrderItems
	Payment            Payment
}

type AdminOrders struct {
	Data []AdminOrder
}

func ToAdminOrders(ordersList []entities.Order) AdminOrders {
	var oResponse AdminOrders
	for _, o := range ordersList {
		oResponse.Data = append(oResponse.Data, ToAdminOrder(o))
	}
	return oResponse
}

func ToAdminOrder(o entities.Order) AdminOrder {
	orderResponse := AdminOrder{
		CustomerID:         o.CustomerID,
		OrderNumber:        o.OrderNumber,
		PaymentStatus:      o.PaymentStatus,
		TotalOriginalPrice: o.TotalOriginalPrice,
		TotalSalePrice:     o.TotalSalePrice,
		Discount:           o.Discount,
		OrderStatus:        o.OrderStatus,
		OrderStatusText:    AdminOrderStatusMap(o.OrderStatus),
		OrderItems:         ToAdminOrderItems(o.OrderItems),
	}

	// بررسی وجود Payment
	if o.Payment.ID != 0 {
		orderResponse.Payment = ToPayment(o.Payment)
	}

	return orderResponse
}

func AdminOrderStatusMap(status uint) string {
	switch status {
	case 0:
		return "در حال پرداخت"
	case 1:
		return "پرداخت شده"
	case 2:
		return "لغو شده"
	default:
		return "نامعلوم"
	}
}

//------------- order item

type AdminOrderItem struct {
	CustomerID         uint `json:"customer_id"`
	OrderID            uint `json:"order_id"`
	ProductID          uint `json:"product_id"`
	InventoryID        uint `json:"inventory_id"`
	Quantity           uint `json:"quantity"`
	OriginalPrice      uint `json:"original_price"`
	SalePrice          uint `json:"sale_price"`
	TotalOriginalPrice uint `json:"total_original_price"`
	TotalSalePrice     uint `json:"total_sale_price"`
}

type AdminOrderItems struct {
	Data []AdminOrderItem `json:"data"`
}

func ToAdminOrderItem(oItem entities.OrderItem) AdminOrderItem {
	return AdminOrderItem{
		CustomerID:         oItem.CustomerID,
		OrderID:            oItem.OrderID,
		ProductID:          oItem.ProductID,
		InventoryID:        oItem.InventoryID,
		Quantity:           oItem.Quantity,
		OriginalPrice:      oItem.OriginalPrice,
		SalePrice:          oItem.SalePrice,
		TotalOriginalPrice: oItem.TotalOriginalPrice,
		TotalSalePrice:     oItem.TotalSalePrice,
	}
}

func ToAdminOrderItems(oItems []entities.OrderItem) AdminOrderItems {
	var orderItems AdminOrderItems
	for _, oItem := range oItems {
		orderItems.Data = append(orderItems.Data, ToAdminOrderItem(oItem))
	}
	return orderItems
}
