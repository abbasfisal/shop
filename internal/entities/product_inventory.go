package entities

import (
	"encoding/json"
	"gorm.io/gorm"
)

type ProductInventory struct {
	gorm.Model
	ProductID      uint
	Quantity       uint
	AttributesJson json.RawMessage
}
