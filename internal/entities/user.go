package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName   string `gorm:"type:varchar(50);not null"`
	LastName    string `gorm:"type:varchar(50);not null"`
	PhoneNumber string `gorm:"type:varchar(11);not null"`
	Password    string
	Type        string `gorm:"type:enum('admin','client');not null;default:'client'"`

	//Address Address //fk (1:1)
	Cart  Cart //fk
	Order Order
}
