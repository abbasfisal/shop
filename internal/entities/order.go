package entities

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	CustomerID         uint
	OrderNumber        string `gorm:"type:varchar(10);unique"`
	PaymentStatus      uint
	TotalOriginalPrice uint
	TotalSalePrice     uint
	Discount           uint
	OrderStatus        uint

	//-- relation

	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;"`
	Payment    Payment     `gorm:"foreignKey:OrderID"`
}

// Order Status Constants
const (
	OrderPending     uint = 1  // در حال بررسی
	OrderConfirmed   uint = 2  // تایید شده
	OrderCancelled   uint = 3  // لغو شده
	OrderPreparing   uint = 4  // در حال آماده‌سازی
	OrderReadyToShip uint = 5  // آماده برای ارسال
	OrderShipped     uint = 6  // ارسال شده
	OrderInTransit   uint = 7  // در مسیر ارسال
	OrderDelivered   uint = 8  // تحویل داده شده
	OrderReturned    uint = 9  // مرجوع شده
	OrderCompleted   uint = 10 // تکمیل شده
	OrderUnderReview uint = 11 // اختلاف یا مشکل
)
