package entities

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID    uint
	ProductID uint

	Count  uint8 `gorm:"type:smallInt"`
	Status uint8

	Product Product
}
