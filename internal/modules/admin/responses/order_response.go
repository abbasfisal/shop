package responses

import (
	"shop/internal/entities"
	"time"
)

type Order struct {
	CustomerID         uint
	OrderID            uint
	OrderNumber        string
	PaymentStatus      int
	TotalOriginalPrice uint
	TotalSalePrice     uint
	Discount           uint
	OrderStatus        uint
	OrderStatusText    string
	CreatedAt          time.Time
	OrderItems         OrderItems
	Payment            Payment
}

type Orders struct {
	Data []Order
}

func ToOrders(ordersList []entities.Order) Orders {

	var oResponse Orders
	for _, o := range ordersList {
		oResponse.Data = append(oResponse.Data, ToOrder(o))
	}
	return oResponse
}
func ToOrder(o entities.Order) Order {

	orderResponse := Order{
		CustomerID:         o.CustomerID,
		OrderID:            o.ID,
		OrderNumber:        o.OrderNumber,
		PaymentStatus:      o.PaymentStatus,
		TotalOriginalPrice: o.TotalOriginalPrice,
		TotalSalePrice:     o.TotalSalePrice,
		Discount:           o.Discount,
		OrderStatus:        o.OrderStatus,
		CreatedAt:          o.CreatedAt,
		OrderStatusText:    OrderStatusMap(o.OrderStatus),
		OrderItems:         ToOrderItems(o.OrderItems),
	}
	if o.Payment.ID != 0 {
		orderResponse.Payment = ToPayment(o.Payment)
	}
	return orderResponse
}

func OrderStatusMap(status uint) string {

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

type OrderItem struct {
	CustomerID         uint
	OrderID            uint
	ProductID          uint
	InventoryID        uint
	Quantity           uint
	OriginalPrice      uint
	SalePrice          uint
	TotalOriginalPrice uint
	TotalSalePrice     uint
}
type OrderItems struct {
	Data []OrderItem
}

func ToOrderItem(oItem entities.OrderItem) OrderItem {
	return OrderItem{
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

func ToOrderItems(oItems []entities.OrderItem) OrderItems {
	var orderItems OrderItems
	for _, oItem := range oItems {
		orderItems.Data = append(orderItems.Data, ToOrderItem(oItem))
	}
	return orderItems
}
