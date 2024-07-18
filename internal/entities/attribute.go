package entities

import "gorm.io/gorm"

type Attribute struct {
	gorm.Model
	CategoryID uint
	Title      string
	Category   Category
}
