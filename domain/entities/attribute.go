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
