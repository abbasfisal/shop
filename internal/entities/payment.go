package entities

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	CustomerID    uint
	OrderId       uint
	PaymentStatus uint8 `gorm:"type:smallInt"`
	Amount        uint
	PaymentMethod uint8
	TransactionID string
	BankResponse  string

	//-- relation
}
