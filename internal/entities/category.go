package entities

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Priority      *uint
	Title         string `gorm:"type:varchar(150);not null"`
	Slug          string `gorm:"type:varchar(150);not null;unique"`
	ParentID      *uint
	Image         string
	Status        bool
	SubCategories []*Category `gorm:"foreignKey:ParentID"`
	Products      []Product   `gorm:"foreignKey:CategoryID"`

	//Attribute []Attribute `gorm:"foreignKey:CategoryID"`
}
