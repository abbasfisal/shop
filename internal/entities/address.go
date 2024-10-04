package entities

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	CustomerID         uint
	ReceiverName       string
	ReceiverMobile     string
	ReceiverAddress    string
	ReceiverPostalCode string

	//------ Relation

}
