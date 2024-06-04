package entities

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	UserID      uint
	Floor       uint8 `gorm:"not null"`
	Number      uint8
	Phase       string `gorm:"type:varchar(50)"`
	Block       string `gorm:"type:varchar(50)"`
	Description string `gorm:"type:varchar(255)"`
}
