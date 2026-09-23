package entities

import "gorm.io/gorm"

type ProductAttribute struct {
	gorm.Model
	ProductID uint

	AttributeID    uint
	AttributeTitle string

	AttributeValueID    uint
	AttributeValueTitle string

	Attribute      Attribute      `gorm:"foreignKey:AttributeID"`
	AttributeValue AttributeValue `gorm:"foreignKey:AttributeValueID"`

	Product Product
}
