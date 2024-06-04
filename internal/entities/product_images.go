package entities

import "gorm.io/gorm"

type ProductImages struct {
	gorm.Model
	ProductID uint
	Path      string
	Product   Product `gorm:"foreignKey:ProductID"`
}
