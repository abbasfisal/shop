package entities

import "gorm.io/gorm"

type Brand struct {
	gorm.Model
	Title string `gorm:"type:varchar(150);not null"`
	Slug  string `gorm:"type:varchar(150);not null;unique"`
	Image string
	//add product fk
}
