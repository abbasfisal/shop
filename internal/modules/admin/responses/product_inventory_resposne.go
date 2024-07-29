package responses

import (
	"encoding/json"
	"shop/internal/entities"
)

type ProductInventory struct {
	ID             uint
	ProductID      uint
	Quantity       uint
	AttributesJson json.RawMessage
}
type ProductInventories struct {
	Data []ProductInventory
}

func ToProductInventory(pi entities.ProductInventory) ProductInventory {
	return ProductInventory{
		ID:             pi.ID,
		ProductID:      pi.ProductID,
		Quantity:       pi.Quantity,
		AttributesJson: pi.AttributesJson,
	}
}

func ToProductInventories(pis []entities.ProductInventory) ProductInventories {
	var response ProductInventories

	for _, productAtt := range pis {
		response.Data = append(response.Data, ToProductInventory(productAtt))
	}

	return response
}

///
