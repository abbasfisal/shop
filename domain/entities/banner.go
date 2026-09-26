package entities

import (
	"time"

	"gorm.io/gorm"
)

// promotion banner layouts (بنرهای پروموشن: دوتایی / چهار تایی)
const (
	BannerLayoutTwo  = "two"
	BannerLayoutFour = "four"
)

type Banner struct {
	gorm.Model
	Type      uint // legacy placement type (kept for old rows)
	Title     string
	Layout    string `gorm:"type:varchar(8);default:'two'"`
	Link      string
	Priority  uint
	SortOrder int        `gorm:"default:0"`
	Status    bool       `gorm:"default:true"`
	StartsAt  *time.Time `gorm:"type:date"`
	EndsAt    *time.Time `gorm:"type:date"`
	Image     string
}

// ValidBannerLayout reports whether the layout is one of the two supported shapes.
func ValidBannerLayout(layout string) bool {
	return layout == BannerLayoutTwo || layout == BannerLayoutFour
}

// BannerLayoutLabel is the Persian label of a promotion layout.
func BannerLayoutLabel(layout string) string {
	switch layout {
	case BannerLayoutFour:
		return "چهار تایی"
	case BannerLayoutTwo:
		return "دوتایی"
	default:
		return "—"
	}
}

// IsVisibleNow is the storefront visibility rule: status on and inside the
// optional [starts_at, ends_at] window (both bounds are optional).
func (b *Banner) IsVisibleNow(now time.Time) bool {
	if !b.Status {
		return false
	}
	if b.StartsAt != nil && now.Before(startOfDay(*b.StartsAt)) {
		return false
	}
	if b.EndsAt != nil && !now.Before(startOfDay(*b.EndsAt).AddDate(0, 0, 1)) {
		// ends_at is inclusive
		return false
	}
	return true
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

const (
	BannerType = iota
	Slider
	LeftSide
	Widget2 // need tow image
	LongHorizontal
)

// BannerTypeLabel is the Persian label of a banner placement (must match the
// options of admin_create_banner.html).
func BannerTypeLabel(t uint) string {
	switch t {
	case Slider:
		return "اسلایدر اصلی (1780×909)"
	case LeftSide:
		return "بنر کناری سمت چپ (856×428)"
	case Widget2:
		return "ویجت ۲ (828×328)"
	case LongHorizontal:
		return "بنر افقی بلند (1656×210)"
	default:
		return "نامشخص"
	}
}

func IsValidBannerType(t uint) bool {
	switch t {
	case Slider, LeftSide, Widget2, LongHorizontal:
		return true
	default:
		return false
	}
}
