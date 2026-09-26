package entities

import "gorm.io/gorm"

type Banner struct {
	gorm.Model
	Type     uint
	Link     string
	Priority uint
	Status   bool
	Image    string
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
