package responses

import (
	"math"
	"time"

	"shop/domain/entities"
)

type Product struct {
	ID            uint
	CategoryID    uint
	BrandID       uint
	Title         string
	Slug          string
	Sku           string
	Status        string
	StatusText    string
	OriginalPrice uint
	SalePrice     uint
	Description   string
	Discount      uint

	// aggregate cache (PricingService) — used by the admin list & storefront list
	ProductType    string
	MinPrice       uint
	MaxPrice       uint
	TotalStock     uint
	TotalReserved  uint
	AvailableStock uint
	InStock        bool
	VariantsCount  int
	ExpiresAt      *time.Time

	//relation
	Category           *Category
	Brand              *Brand
	Images             *ImageProducts
	ProductAttributes  *ProductAttributes
	ProductInventories *ProductInventories
	Features           *Features
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
		StatusText:    entities.ProductStatusLabel(p.Status),
		OriginalPrice: p.OriginalPrice,
		SalePrice:     p.SalePrice,
		Description:   p.Description,
		ExpiresAt:     p.ExpiresAt,

		ProductType:    p.ProductType,
		MinPrice:       p.MinPrice,
		MaxPrice:       p.MaxPrice,
		TotalStock:     p.TotalStock,
		TotalReserved:  p.TotalReserved,
		AvailableStock: p.AvailableStock,
		InStock:        p.InStock,
		VariantsCount:  p.VariantsCount,

		Discount: func() uint {
			if p.OriginalPrice == 0 {
				return 0
			}
			originalPrice := float64(p.OriginalPrice)
			salePrice := float64(p.SalePrice)
			dis := ((originalPrice - salePrice) / originalPrice) * 100

			return uint(math.Round(dis))
		}(),
	}

	if p.Features != nil {
		product.Features = ToFeatures(p.Features)
	}

	product.ProductInventories = ToProductInventoriesWithFallback(p.ProductVariants, p.OriginalPrice, p.SalePrice)
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
