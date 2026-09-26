package entities

import "time"

// ProductReadModel the flattened product document persisted in products.read_model (JSONB).
// It replaces the old MongoDB products collection.
type ProductReadModel struct {
	Product     P                    `json:"product"`
	Inventories map[string]Inventory `json:"inventories"`
}

// P is the flattened product projection consumed by the storefront templates.
type P struct {
	ID            int64  `json:"id"`
	Category      C      `json:"Category"`
	CategoryID    int64  `json:"category_id"`
	Brand         B      `json:"Brand"`
	BrandID       int64  `json:"brand_id"`
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Sku           string `json:"sku"`
	Status        string `json:"status"`
	OriginalPrice int64  `json:"original_price"`
	SalePrice     int64  `json:"sale_price"`
	Discount      int64  `json:"Discount"`
	// PricingService aggregates mirrored for the storefront (effective price
	// display: "از MinPrice", single-inventory pages, recommendations).
	MinPrice    int64     `json:"min_price"`
	MaxPrice    int64     `json:"max_price"`
	Description string    `json:"description"`
	Images      Img       `json:"Images"`
	Features    F         `json:"Features"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type C struct {
	ID       int64  `json:"id"`
	ParentID int64  `json:"parent_id"`
	Title    string `json:"title"`
	Slug     string `json:"slug"`
}

type B struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type Img struct {
	Data []ImgData `json:"Data"`
}

type ImgData struct {
	ID           int64  `json:"ID"`
	OriginalPath string `json:"OriginalPath"`
	FullPath     string `json:"FullPath"`
}

type F struct {
	Data []FData `json:"Data"`
}

type FData struct {
	ID        int64  `json:"ID"`
	ProductID int64  `json:"ProductID"`
	Title     string `json:"Title"`
	Value     string `json:"Value"`
}

type Inventory struct {
	InventoryID int64                 `json:"inventory_id"`
	Quantity    int64                 `json:"quantity"`
	Attributes  []InventoryAttributes `json:"attributes"`

	// per-variant pricing (Laravel ProductVariant) — NULL variant columns
	// already resolved against the product price here
	Price           int64  `json:"price"`          // crossed-out list price
	SalePrice       int64  `json:"sale_price"`     // normal price
	DiscountPrice   int64  `json:"discount_price"` // 0 when none
	HasDiscount     bool   `json:"has_discount"`
	DiscountPercent int64  `json:"discount_percent"`
	EffectivePrice  int64  `json:"effective_price"` // what the customer pays
	Status          string `json:"status"`          // active | inactive
	Available       int64  `json:"available"`       // stock - reserved
}

type InventoryAttributes struct {
	AttributeID                 int64  `json:"attribute_id"`
	AttributeTitle              string `json:"attribute_title"`
	AttributeValueID            int64  `json:"attribute_value_id"`
	AttributeValueTitle         string `json:"attribute_value_title"`
	ProductInventoryAttributeID int64  `json:"product_inventory_attribute_id"`
	// dynamic value presentation: color attributes paint swatches
	IsColor  bool   `json:"is_color"`
	ColorHex string `json:"color_hex"`
}
