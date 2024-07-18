package entities

import "gorm.io/gorm"

type AttributeValue struct {
	gorm.Model
	AttributeID    uint
	AttributeTitle string
	Value          string
}
