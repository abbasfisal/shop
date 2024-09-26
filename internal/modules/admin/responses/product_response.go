package responses

import "shop/internal/entities"

type Product struct {
	ID            uint
	CategoryID    uint
	BrandID       uint
	Title         string
	Slug          string
	Sku           string
	Status        bool
	OriginalPrice uint
	SalePrice     uint
	Description   string

	//relation
	Category                   Category
	Brand                      Brand
	Images                     ImageProducts
	ProductAttributes          ProductAttributes
	ProductInventories         ProductInventories
	ProductInventoryAttributes ProductInventoryAttributes
	Features                   Features
}

type Products struct {
	Data []Product
}

func ToProducts(products []entities.Product) Products {
	var pResponse Products
	for _, p := range products {
		pResponse.Data = append(pResponse.Data, ToProduct(p))
	}
	return pResponse
}
func ToProduct(p entities.Product) Product {
	return Product{
		ID:            p.ID,
		CategoryID:    p.CategoryID,
		BrandID:       p.BrandID,
		Title:         p.Title,
		Slug:          p.Slug,
		Sku:           p.Sku,
		Status:        p.Status,
		OriginalPrice: p.OriginalPrice,
		SalePrice:     p.SalePrice,
		Description:   p.Description,

		//relation
		Category:                   ToCategory(p.Category),
		Brand:                      ToBrand(p.Brand),
		Images:                     ToImageProducts(p.ProductImages),
		ProductAttributes:          ToProductAttributes(p.ProductAttributes),
		ProductInventories:         ToProductInventories(p.ProductInventories),
		ProductInventoryAttributes: ToProductInventoryAttributes(p.ProductInventoryAttributes),
		Features:                   ToFeatures(p.Features),
	}
}
