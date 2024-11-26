package entities

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	CustomerID  uint
	OrderID     uint
	Authority   string `gorm:"unique"`
	Description string
	PaymentURL  string
	StatusCode  int
	Amount      uint
	RefID       string
	Status      int //payment status -> 0-pending,1-paid ,2-failed

	//-- relation
}

// Payment Status Constants
const (
	PaymentPending uint = 1 // منتظر پرداخت
	PaymentSuccess uint = 2 // پرداخت موفق
	PaymentFailed  uint = 3 // پرداخت ناموفق
	PaymentRetry   uint = 4 // در انتظار پرداخت مجدد
)
