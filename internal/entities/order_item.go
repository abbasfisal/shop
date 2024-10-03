package entities

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	CustomerID         uint
	OrderID            uint
	ProductID          uint
	InventoryID        uint
	Quantity           uint8 `gorm:"type:SmallInt"`
	OriginalPrice      uint
	SalePrice          uint
	TotalOriginalPrice uint
	TotalSalePrice     uint

	//-- relation

}
