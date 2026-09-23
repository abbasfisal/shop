package responses

import (
	"shop/domain/entities"
)

type ProductInventory struct {
	ID        uint
	ProductID uint
	Quantity  uint
}
type ProductInventories struct {
	Data []ProductInventory
}

func ToProductInventory(pv *entities.ProductVariant) *ProductInventory {
	return &ProductInventory{
		ID:        pv.ID,
		ProductID: pv.ProductID,
		Quantity:  pv.Stock,
	}
}

func ToProductInventories(variants []*entities.ProductVariant) *ProductInventories {
	var pResponse ProductInventories
	for _, pv := range variants {
		pResponse.Data = append(pResponse.Data, *ToProductInventory(pv))
	}
	return &pResponse
}
