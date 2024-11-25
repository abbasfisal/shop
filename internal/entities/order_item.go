package entities

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	CustomerID         uint
	OrderID            uint
	ProductID          uint
	InventoryID        uint
	Quantity           uint
	OriginalPrice      uint
	SalePrice          uint
	TotalOriginalPrice uint
	TotalSalePrice     uint

	//-- relation
	Product Product
}
