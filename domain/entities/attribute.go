package entities

import (
	"fmt"

	"gorm.io/gorm"
)

type Attribute struct {
	gorm.Model
	//CategoryID      uint
	Title string `gorm:"type:varchar(150);not null"`
	Code  string `gorm:"type:varchar(64);uniqueIndex"`
	// input_type decides the admin value-form widget and the storefront
	// rendering: text (default) or color (picker + swatches)
	InputType string `gorm:"type:varchar(16);default:'text'"`
	//Category        Category
	SortOrder       int               `gorm:"default:0"`
	AttributeValues []*AttributeValue `gorm:"foreignKey:AttributeID"` //1:M
}

// AfterCreate generates a stable code (attr_<id>) when the admin form only
// posted a title — attributes_json keys and filtering rely on code being set.
func (a *Attribute) AfterCreate(tx *gorm.DB) error {
	if a.Code == "" && a.ID > 0 {
		return tx.Model(a).Update("code", fmt.Sprintf("attr_%d", a.ID)).Error
	}
	return nil
}

// Attribute input types (dynamic value widgets).
const (
	AttributeInputText  = "text"
	AttributeInputColor = "color"
)

// IsColor reports whether values of this attribute carry a hex color.
func (a *Attribute) IsColor() bool {
	return a != nil && a.InputType == AttributeInputColor
}

// ValidAttributeInputType reports whether the input type is supported.
func ValidAttributeInputType(inputType string) bool {
	return inputType == AttributeInputText || inputType == AttributeInputColor
}

// AttributeInputTypeLabel is the Persian label of an input type.
func AttributeInputTypeLabel(inputType string) string {
	if inputType == AttributeInputColor {
		return "رنگ (کالرپیکر)"
	}
	return "متن ساده"
}
