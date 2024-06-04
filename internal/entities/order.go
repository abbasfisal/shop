package entities

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID         uint
	TrackingNumber string `gorm:"type:varchar(10)"`
	ProductID      uint
	Count          uint `gorm:"type:SmallInt"`
	Status         uint `gorm:"type:smallInt"`
	//add bank url , etc
}
