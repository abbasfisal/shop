package entities

import (
	"time"

	"gorm.io/gorm"
)

// site banner placements (جایگاه بنرهای ثابت سایت)
const (
	SiteBannerPlacementHeader = "header"
	SiteBannerPlacementMain   = "main"
)

// SiteBanner is a fixed (non-promotion) site banner: the strip above the
// header on /home2 plus the reserved slot for a future image hero banner.
type SiteBanner struct {
	gorm.Model
	// no `default` tag on Placement/Status on purpose: GORM replaces a zero
	// value with the tag default on Create, which would silently turn
	// غیرفعال into فعال.
	Placement   string `gorm:"type:varchar(16)"`
	Title       string
	Link        string
	Image       string // desktop image (required)
	MobileImage string // mobile image (optional → falls back to Image)
	Status      bool
	StartsAt    *time.Time `gorm:"type:date"`
	EndsAt      *time.Time `gorm:"type:date"`
	SortOrder   int        `gorm:"default:0"`
}

// ValidSiteBannerPlacement reports whether the placement exists.
func ValidSiteBannerPlacement(placement string) bool {
	return placement == SiteBannerPlacementHeader || placement == SiteBannerPlacementMain
}

// SiteBannerPlacementLabel is the Persian label of a placement.
func SiteBannerPlacementLabel(placement string) string {
	switch placement {
	case SiteBannerPlacementHeader:
		return "بنر هدر"
	case SiteBannerPlacementMain:
		return "بنر تصویری اسلایدر اصلی"
	default:
		return "—"
	}
}

// NormalizeSiteBannerPlacement defaults an empty placement to the header slot.
func NormalizeSiteBannerPlacement(placement string) string {
	if placement == "" {
		return SiteBannerPlacementHeader
	}
	return placement
}

// IsVisibleNow is the storefront visibility rule: status on and inside the
// optional [starts_at, ends_at] window (both bounds optional).
func (b *SiteBanner) IsVisibleNow(now time.Time) bool {
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

// EffectiveMobileImage returns the mobile image, falling back to the desktop
// one when no mobile variant was uploaded.
func (b *SiteBanner) EffectiveMobileImage() string {
	if b.MobileImage != "" {
		return b.MobileImage
	}
	return b.Image
}
