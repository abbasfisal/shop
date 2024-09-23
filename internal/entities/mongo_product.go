package entities

import (
	"time"
)

type MongoProduct struct {
	ID            uint     `bson:"id"`
	CategoryID    uint     `bson:"category_id"`
	BrandID       uint     `bson:"brand_id"`
	Title         string   `bson:"title"`
	Slug          string   `bson:"slug"`
	Sku           string   `bson:"sku"`
	Status        bool     `bson:"status"`
	OriginalPrice uint     `bson:"original_price"`
	SalePrice     uint     `bson:"sale_price"`
	Description   string   `bson:"description"`
	Images        []string `bson:"images"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}
