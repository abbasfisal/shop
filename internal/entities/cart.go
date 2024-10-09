package entities

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	CustomerID uint
	Status     uint8

	//---- relation
	CartItems []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
}
