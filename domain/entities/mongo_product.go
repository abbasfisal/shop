package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MongoProduct struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Product     P                    `bson:"product"`
	Inventories map[string]Inventory `bson:"inventories"`
}

type P struct {
	ID int64 `bson:"id"`

	Category   C     `bson:"Category"`
	CategoryID int64 `bson:"category_id"`

	Brand   B     `bson:"Brand"`
	BrandID int64 `bson:"brand_id"`

	Title         string `bson:"title"`
	Slug          string `bson:"slug"`
	Sku           string `bson:"sku"`
	Status        bool   `bson:"status"`
	OriginalPrice int64  `bson:"original_price"`
	SalePrice     int64  `bson:"sale_price"`
	Discount      int64  `bson:"Discount"`
	Description   string `bson:"description"`

	//Images []string `bson:"images"`
	Images   Img `bson:"Images"`
	Features F   `bson:"Features"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type C struct {
	ID       int64  `bson:"id"`
	ParentID int64  `bson:"parent_id"`
	Title    string `bson:"title"`
	Slug     string `bson:"slug"`
}

type B struct {
	ID    int64  `bson:"id"`
	Title string `bson:"title"`
	Slug  string `bson:"slug"`
}

type Img struct {
	Data []ImgData `bson:"Data"`
}

type ImgData struct {
	ID           int64  `bson:"ID"` //todo:correct capital
	OriginalPath string `bson:"OriginalPath"`
	FullPath     string `bson:"FullPath"`
}

type F struct {
	Data []FData `bson:"Data"`
}

type FData struct {
	ID        int64  `bson:"ID"`
	ProductID int64  `bson:"ProductID"`
	Title     string `bson:"Title"`
	Value     string `bson:"Value"`
}

type Inventory struct {
	InventoryID int64                 `json:"inventory_id" bson:"inventory_id"`
	Quantity    int64                 `json:"quantity" bson:"quantity"`
	Attributes  []InventoryAttributes `json:"attributes" bson:"attributes"`
}
type InventoryAttributes struct {
	AttributeID                 int64  `json:"attribute_id" bson:"attribute_id"`
	AttributeTitle              string `json:"attribute_title" bson:"attribute_title"`
	AttributeValueID            int64  `json:"attribute_value_id" bson:"attribute_value_id"`
	AttributeValueTitle         string `json:"attribute_value_title" bson:"attribute_value_title"`
	ProductInventoryAttributeID int64  `json:"product_inventory_attribute_id" bson:"product_inventory_attribute_id"`
}
