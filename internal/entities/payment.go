package entities

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	CustomerID    uint
	OrderID       uint
	PaymentStatus uint

	Amount        uint
	PaymentMethod uint
	TransactionID string
	BankResponse  string

	//-- relation
}
