package responses

import (
	"encoding/json"
	"fmt"
	"log"
	"shop/internal/entities"
)

type ProductInventory struct {
	ID             uint
	ProductID      uint
	Quantity       uint
	AttributesJson []ProductAttribute
}
type ProductInventories struct {
	Data []ProductInventory
}

func ToProductInventory(pi entities.ProductInventory) ProductInventory {
	var pa []ProductAttribute

	err := json.Unmarshal([]byte(pi.AttributesJson), &pa)
	if err != nil {
		fmt.Println("----- single to ToProductInventory unmarshal failed ", err)
		return ProductInventory{}
	}

	return ProductInventory{
		ID:             pi.ID,
		ProductID:      pi.ProductID,
		Quantity:       pi.Quantity,
		AttributesJson: pa,
	}
}

func ToProductInventories(pis []entities.ProductInventory) ProductInventories {
	var response ProductInventories

	for _, productAtt := range pis {
		var at []ProductAttribute

		err := json.Unmarshal([]byte(productAtt.AttributesJson), &at)
		if err != nil {
			log.Fatal("-------- fail ToProductInventories while unmarshalling --", err)
			return ProductInventories{}
		}

		response.Data = append(response.Data, ToProductInventory(productAtt))
	}

	return response
}

///
