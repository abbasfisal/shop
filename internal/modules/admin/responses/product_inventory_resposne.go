package responses

import (
	"shop/internal/entities"
)

type ProductInventory struct {
	ID        uint
	ProductID uint
	Quantity  uint
}
type ProductInventories struct {
	Data []ProductInventory
}

func ToProductInventory(pi *entities.ProductInventory) *ProductInventory {
	return &ProductInventory{
		ID:        pi.ID,
		ProductID: pi.ProductID,
		Quantity:  pi.Quantity,
	}
}

func ToProductInventories(pis []*entities.ProductInventory) *ProductInventories {
	var pResponse ProductInventories
	for _, p := range pis {
		pResponse.Data = append(pResponse.Data, *ToProductInventory(p))
	}
	return &pResponse
}
