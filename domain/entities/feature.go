package entities

import "gorm.io/gorm"

type Feature struct {
	gorm.Model
	ProductID uint
	Title     string
	Value     string
}
