package entities

import "gorm.io/gorm"

type AttributeValue struct {
	gorm.Model
	AttributeID    uint
	AttributeTitle string
	Value          string
	SortOrder      int `gorm:"default:0"`

	//relation
	Attribute Attribute `gorm:"foreignKey:AttributeID"`
}
