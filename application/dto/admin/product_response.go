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
	// EffectivePrice is the customer-facing price: the PricingService minimum
	// when the product has priced variants, otherwise the nominal sale price.
	EffectivePrice uint

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

		EffectivePrice: EffectivePriceOf(p),

		// golden rule: the discount is measured against the price the
		// customer actually pays (effective), not the nominal sale price —
		// otherwise a variant-level تخفیف shows a wrong (lower) percent.
		Discount: DiscountPercentOf(p.OriginalPrice, EffectivePriceOf(p)),
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

// EffectivePriceOf resolves the customer-facing price of a product:
// the PricingService minimum when the product has priced variants,
// otherwise the nominal sale price.
func EffectivePriceOf(p *entities.Product) uint {
	if p.MinPrice > 0 {
		return p.MinPrice
	}
	return p.SalePrice
}

// DiscountPercentOf is the تومان discount percent of a crossed-out list
// price against the price the customer pays (0 when there is no real
// discount or the inputs are degenerate).
func DiscountPercentOf(originalPrice, effectivePrice uint) uint {
	if originalPrice == 0 || effectivePrice == 0 || effectivePrice >= originalPrice {
		return 0
	}
	return uint(math.Round((float64(originalPrice) - float64(effectivePrice)) / float64(originalPrice) * 100))
}
