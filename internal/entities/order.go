package entities

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	CustomerID         uint
	OrderNumber        string `gorm:"type:varchar(10),unique"`
	PaymentStatus      uint   `gorm:"type:smallInt"`
	TotalOriginalPrice uint
	TotalSalePrice     uint
	Discount           uint
	OrderStatus        uint8 `gorm:"type:smallInt"`

	//-- relation

	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
	Payments   []Payment   `gorm:"foreignKey:OrderID"`
}
