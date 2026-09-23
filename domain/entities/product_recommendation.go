package entities

import "time"

// ProductRecommendation is a row in product_recommendations (replaces the MongoDB
// recommendations collection). Hard deletes keep the UNIQUE constraint reusable.
type ProductRecommendation struct {
	ID                   uint      `gorm:"primaryKey"`
	ProductID            uint      `gorm:"not null;index"`
	RecommendedProductID uint      `gorm:"not null"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (ProductRecommendation) TableName() string {
	return "product_recommendations"
}

// RecommendedProduct is the storefront shape used by single_product.html
// (template accesses $rec.Product.<field>).
type RecommendedProduct struct {
	Product P `json:"product"`
}
