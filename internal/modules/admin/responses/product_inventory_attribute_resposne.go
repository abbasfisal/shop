package responses

import (
	"shop/internal/entities"
)

type ProductInventoryAttribute struct {
	ID                 uint
	ProductID          uint
	ProductInventoryID uint
	ProductAttributeID uint
	ProductAttribute   *ProductAttribute
}

type ProductInventoryAttributes struct {
	Data []ProductInventoryAttribute
}

func ToProductInventoryAttribute(pi *entities.ProductInventoryAttribute) *ProductInventoryAttribute {
	pia := ProductInventoryAttribute{
		ID:        pi.ID,
		ProductID: pi.ProductID,

		ProductInventoryID: pi.ProductInventoryID,
		ProductAttributeID: pi.ProductAttributeID,
	}

	if pi.ProductAttribute != nil {
		pia.ProductAttribute = ToProductAttribute(pi.ProductAttribute)
	}

	return &pia
}

func ToProductInventoryAttributes(pis []*entities.ProductInventoryAttribute) *ProductInventoryAttributes {

	if pis == nil {
		return &ProductInventoryAttributes{}
	}

	var pResponse ProductInventoryAttributes
	for _, p := range pis {
		pResponse.Data = append(pResponse.Data, *ToProductInventoryAttribute(p))
	}
	return &pResponse
}
