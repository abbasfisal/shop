package responses

import (
	"math"
	"shop/internal/entities"
)

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
	Discount      uint

	//relation
	Category                   *Category
	Brand                      *Brand
	Images                     *ImageProducts
	ProductAttributes          *ProductAttributes
	ProductInventories         *ProductInventories
	ProductInventoryAttributes *ProductInventoryAttributes
	Features                   *Features
}

type Products struct {
	Data []Product
}

func ToProducts(products []*entities.Product) *Products {
	var pResponse Products
	for _, p := range products {
		pResponse.Data = append(pResponse.Data, *ToProduct(p))
	}
	return &pResponse
}

func ToProduct(p *entities.Product) *Product {
	var product = Product{
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
		Discount: func() uint {
			originalPrice := float64(p.OriginalPrice)
			salePrice := float64(p.SalePrice)
			dis := ((originalPrice - salePrice) / originalPrice) * 100

			return uint(math.Round(dis))
		}(),
	}

	if p.Features != nil {
		product.Features = ToFeatures(p.Features)
	}

	if p.ProductInventoryAttributes != nil {
		product.ProductInventoryAttributes = ToProductInventoryAttributes(p.ProductInventoryAttributes)
	}

	if p.ProductInventories != nil {
		product.ProductInventories = ToProductInventories(p.ProductInventories)
	}
	if p.ProductAttributes != nil {
		product.ProductAttributes = ToProductAttributes(p.ProductAttributes)
	}

	if p.Category != nil {
		product.Category = ToCategory(p.Category)
	}

	if p.Brand != nil {
		product.Brand = ToBrand(p.Brand)
	}

	if p.ProductImages != nil {
		product.Images = ToImageProducts(p.ProductImages)
	}

	return &product
}
