package entities

import (
	"time"

	"gorm.io/gorm"
)

// Homepage slots a product slider can be attached to (جایگاه — the admin can
// move a slider between them at any time).
const (
	SliderPositionAfterSlider     = "after_slider"
	SliderPositionMidContent      = "mid_content"
	SliderPositionAfterCategories = "after_categories"
	SliderPositionBeforeFooter    = "before_footer"
)

// MaxSliderProducts is the hard cap of products shown in one slider (spec: 10).
const MaxSliderProducts = 10

// SliderPositionOption is one selectable slot for the admin forms.
type SliderPositionOption struct {
	Value string
	Label string
}

// SliderPositions returns every homepage slot in rendering order.
func SliderPositions() []SliderPositionOption {
	return []SliderPositionOption{
		{Value: SliderPositionAfterSlider, Label: "بعد از اسلایدر اصلی"},
		{Value: SliderPositionMidContent, Label: "وسط صفحه"},
		{Value: SliderPositionAfterCategories, Label: "بعد از دسته‌بندی‌ها"},
		{Value: SliderPositionBeforeFooter, Label: "قبل از فوتر"},
	}
}

// ValidSliderPosition reports whether the slot exists.
func ValidSliderPosition(position string) bool {
	for _, p := range SliderPositions() {
		if p.Value == position {
			return true
		}
	}
	return false
}

// SliderPositionLabel is the Persian label of a slot.
func SliderPositionLabel(position string) string {
	for _, p := range SliderPositions() {
		if p.Value == position {
			return p.Label
		}
	}
	return position
}

// ProductSlider is a curated strip of products on the homepage
// (حداکثر ۱۰ محصول + دکمه «مشاهده همه» به صفحه کاتالوگ).
type ProductSlider struct {
	gorm.Model
	Title      string
	Slug       string `gorm:"unique"`
	Position   string `gorm:"type:varchar(32);default:'after_categories'"`
	CategoryID *uint
	Status     string     `gorm:"type:varchar(16);default:'published'"`
	StartsAt   *time.Time `gorm:"type:date"`
	EndsAt     *time.Time `gorm:"type:date"`

	Category *Category        `gorm:"foreignKey:CategoryID"`
	Products []*SliderProduct `gorm:"foreignKey:SliderID"`
}

// SliderProduct is one product inside a slider (ordered by sort_order).
type SliderProduct struct {
	gorm.Model
	SliderID  uint `gorm:"index"`
	ProductID uint `gorm:"index"`
	SortOrder int  `gorm:"default:0"`

	Slider  *ProductSlider `gorm:"foreignKey:SliderID"`
	Product *Product       `gorm:"foreignKey:ProductID"`
}

// IsVisibleNow is the storefront visibility rule for a slider.
func (s *ProductSlider) IsVisibleNow(now time.Time) bool {
	if s.Status != ProductStatusPublished {
		return false
	}
	if s.StartsAt != nil && now.Before(startOfDay(*s.StartsAt)) {
		return false
	}
	if s.EndsAt != nil && !now.Before(startOfDay(*s.EndsAt).AddDate(0, 0, 1)) {
		return false
	}
	return true
}
