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

func IsValidBannerType(t uint) bool {
	switch t {
	case Slider, LeftSide, Widget2, LongHorizontal:
		return true
	default:
		return false
	}
}
