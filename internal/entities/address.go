package entities

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	CustomerID uint

	ZipCode string
	Unit    string
	Floor   uint8 `gorm:"not null"`
	Number  string

	Latitude  float64
	Longitude float64

	//Phase       string `gorm:"type:varchar(50)"`
	//Block       string `gorm:"type:varchar(50)"`

	Description string `gorm:"type:varchar(255)"`

	//------ Relation

}
