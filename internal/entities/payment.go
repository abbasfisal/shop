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
