package entities

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	CustomerID         uint
	OrderNumber        string `gorm:"type:varchar(10);unique"`
	PaymentStatus      int
	TotalOriginalPrice uint
	TotalSalePrice     uint
	Discount           uint
	OrderStatus        uint

	Address string
	Note    string

	//-- relation

	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;"`
	Payment    Payment     `gorm:"foreignKey:OrderID"`
}

// Order Status Constants
const (
	OrderPending     uint = iota // در حال پرداخت
	OrderConfirmed               // پرداخت شده
	OrderCancelled               // لغو شده
	OrderPreparing               // در حال آماده‌سازی
	OrderReadyToShip             // آماده برای ارسال
	OrderShipped                 // ارسال شده
	OrderInTransit               // در مسیر ارسال
	OrderDelivered               // تحویل داده شده
	OrderReturned                // مرجوع شده
	OrderCompleted               // تکمیل شده
	OrderUnderReview             // اختلاف یا مشکل
)
