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
