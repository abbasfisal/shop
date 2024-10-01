package entities

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Mobile    string `gorm:"index:unique"`
	FirstName string
	LastName  string
	Active    bool
	//Gender    bool //true:male , false:female
	//DateOfBirth

	//--------- Relation
	Address Address   `gorm:"foreignKey:CustomerID"`
	Carts   []Cart    `gorm:"foreignKey:CustomerID"`
	Session []Session `gorm:"foreignKey:CustomerID"`
}
