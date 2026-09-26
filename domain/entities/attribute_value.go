package entities

import (
	"encoding/json"
	"regexp"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AttributeValue struct {
	gorm.Model
	AttributeID    uint
	AttributeTitle string
	Value          string
	SortOrder      int `gorm:"default:0"`
	// Meta carries dynamic extras per value; color values store
	// {"hex": "#rrggbb"} so the storefront can paint swatches.
	Meta datatypes.JSON `gorm:"type:jsonb;default:'{}'"`

	//relation
	Attribute Attribute `gorm:"foreignKey:AttributeID"`
}

var hexColorRe = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

// ValidHexColor reports whether s is a #rrggbb color.
func ValidHexColor(s string) bool {
	return hexColorRe.MatchString(s)
}

// ColorHex returns the stored hex color ("" when missing or malformed).
func (v *AttributeValue) ColorHex() string {
	if v == nil || len(v.Meta) == 0 {
		return ""
	}
	var meta map[string]string
	if err := json.Unmarshal(v.Meta, &meta); err != nil {
		return ""
	}
	hex, ok := meta["hex"]
	if !ok || !ValidHexColor(hex) {
		return ""
	}
	return hex
}

// SetColorHex writes {"hex": ...} into Meta (validated upstream).
func (v *AttributeValue) SetColorHex(hex string) {
	meta, _ := json.Marshal(map[string]string{"hex": hex})
	v.Meta = datatypes.JSON(meta)
}
