package responses

import (
	"encoding/json"
	"shop/internal/entities"
	"shop/internal/pkg/util"
	"time"
)

type CustomerOrder struct {
	//ID                 uint
	//CustomerID         uint
	OrderNumber              string
	PaymentStatus            int
	TotalOriginalPrice       uint
	TotalSalePrice           uint
	PrettyTotalSalePrice     string
	PrettyTotalOriginalPrice string
	CreatedAt                time.Time
	Discount                 uint
	//OrderNote          string
	//OrderStatus     uint
	OrderStatusText string
	OrderItems      CustomerOrderItems
	//Payment         Payment
	Address Address
}

type CustomerOrders struct {
	Data []CustomerOrder
}

func ToCustomerOrders(ordersList []entities.Order) CustomerOrders {
	var oResponse CustomerOrders
	for _, o := range ordersList {
		oResponse.Data = append(oResponse.Data, ToCustomerOrder(&o))
	}
	return oResponse
}

func ToCustomerOrder(o *entities.Order) CustomerOrder {
	orderResponse := CustomerOrder{
		//ID:                 o.ID,
		//CustomerID:         o.CustomerID,
		OrderNumber:              o.OrderNumber,
		PaymentStatus:            o.PaymentStatus,
		TotalOriginalPrice:       o.TotalOriginalPrice,
		TotalSalePrice:           o.TotalSalePrice,
		PrettyTotalOriginalPrice: util.PrettyPrice(int(o.TotalOriginalPrice)),
		PrettyTotalSalePrice:     util.PrettyPrice(int(o.TotalSalePrice)),
		CreatedAt:                o.CreatedAt,
		Discount:                 o.Discount,
		//OrderStatus:        o.OrderStatus,
		//OrderNote:          o.Note,
		OrderStatusText: CustomerOrderStatusMap(o.OrderStatus),
		OrderItems:      ToCustomerOrderItems(o.OrderItems),
		Address:         ToAddress(o.Address),
	}

	// بررسی وجود Payment
	//if o.Payment.ID != 0 {
	//	orderResponse.Payment = ToPayment(o.Payment)
	//}

	return orderResponse
}

func CustomerOrderStatusMap(status uint) string {

	switch status {
	case entities.OrderPending: //0
		return "در حال پرداخت"
	case entities.OrderConfirmed:
		return "پرداخت شده"
	case entities.OrderCancelled:
		return "لغو شده"
	case entities.OrderPreparing:
		return "در حال آماده سازی"
	case entities.OrderReadyToShip:
		return "آماده برای ارسال"
	case entities.OrderShipped:
		return "ارسال شده"
	case entities.OrderInTransit:
		return "در مسیر ارسال"
	case entities.OrderDelivered:
		return "تحویل داده شده"
	case entities.OrderReturned:
		return "مرجوع شده"
	case entities.OrderCompleted:
		return "تکمیل شده"
	case entities.OrderUnderReview:
		return "اختلاف یا مشکل"
	default:
		return "نامعلوم"
	}
}

//------------- order item

type CustomerOrderItem struct {
	CustomerID         uint
	OrderID            uint
	InventoryID        uint
	Quantity           uint
	OriginalPrice      uint
	SalePrice          uint
	TotalOriginalPrice uint
	TotalSalePrice     uint

	ProductID            uint
	ProductTitle         string
	ProductOriginalPrice uint
	ProductSalePrice     uint
	ProductSku           string
	ProductSlug          string
	OrderItemAttributes  OrderItemAttributes
}

type CustomerOrderItems struct {
	Data []CustomerOrderItem
}

func ToCustomerOrderItem(oItem entities.OrderItem) CustomerOrderItem {
	return CustomerOrderItem{
		CustomerID:         oItem.CustomerID,
		OrderID:            oItem.OrderID,
		ProductID:          oItem.ProductID,
		InventoryID:        oItem.InventoryID,
		Quantity:           oItem.Quantity,
		OriginalPrice:      oItem.OriginalPrice,
		SalePrice:          oItem.SalePrice,
		TotalOriginalPrice: oItem.TotalOriginalPrice,
		TotalSalePrice:     oItem.TotalSalePrice,

		ProductTitle:         oItem.Product.Title,
		ProductOriginalPrice: oItem.OriginalPrice,
		ProductSalePrice:     oItem.SalePrice,
		ProductSku:           oItem.Product.Sku,
		ProductSlug:          oItem.Product.Slug,

		OrderItemAttributes: ToOrderItemAttributes(oItem.Product.ProductInventoryAttributes),
	}
}

func ToCustomerOrderItems(oItems []entities.OrderItem) CustomerOrderItems {
	var orderItems CustomerOrderItems
	for _, oItem := range oItems {
		orderItems.Data = append(orderItems.Data, ToCustomerOrderItem(oItem))
	}
	return orderItems
}

type OrderItemAttributes struct {
	Data []OrderItemAttribute
}
type OrderItemAttribute struct {
	Title string
	Value string
}

func ToOrderItemAttribute(oItemAttribute entities.ProductAttribute) OrderItemAttribute {
	return OrderItemAttribute{
		Title: oItemAttribute.AttributeTitle,
		Value: oItemAttribute.AttributeValueTitle,
	}
}
func ToOrderItemAttributes(oItemAttributes []entities.ProductInventoryAttribute) OrderItemAttributes {
	var a OrderItemAttributes
	for _, i2 := range oItemAttributes {
		a.Data = append(a.Data, ToOrderItemAttribute(i2.ProductAttribute))
	}
	return a
}

type address struct {
	ReceiverName       string
	ReceiverMobile     string
	ReceiverAddress    string
	ReceiverPostalCode string
}

func ToAddress(address string) Address {
	var add Address
	err := json.Unmarshal([]byte(address), &add)
	if err != nil {
		return add
	}
	return add
}
