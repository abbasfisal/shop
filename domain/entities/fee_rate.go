package entities

import (
	"time"

	"gorm.io/gorm"
)

// Fee kinds (هزینه ارسال / هزینه بسته‌بندی) — two admin menus share one table.
const (
	FeeKindShipping  = "shipping"
	FeeKindPackaging = "packaging"
)

// FeeKindLabel is the Persian label of a fee kind.
func FeeKindLabel(kind string) string {
	switch kind {
	case FeeKindPackaging:
		return "هزینه بسته‌بندی"
	case FeeKindShipping:
		return "هزینه ارسال"
	default:
		return kind
	}
}

// ValidFeeKind reports whether the kind is a manageable fee type.
func ValidFeeKind(kind string) bool {
	return kind == FeeKindShipping || kind == FeeKindPackaging
}

// FeeRate is one time-windowed tariff row. A new row per period keeps the
// price history; the storefront always resolves the single active row.
type FeeRate struct {
	gorm.Model
	Kind   string `gorm:"type:varchar(16)"`
	Title  string
	Amount uint
	// FreeThreshold (shipping only): orders with an items total at/above this
	// amount ship free. 0 (or NULL) disables the rule.
	FreeThreshold *uint      `gorm:"column:free_threshold"`
	StartsAt      *time.Time `gorm:"type:date"`
	EndsAt        *time.Time `gorm:"type:date"`
	// no `+"`gorm:default:true`"+` on purpose: GORM replaces a struct zero value
	// carrying a default tag on Create, turning غیرفعال into فعال.
	// (the SQL column keeps DEFAULT TRUE for raw inserts)
	Status bool
}

// IsActiveAt is the visibility rule shared by the admin list and the
// storefront resolution: switched on and inside the optional window.
func (f *FeeRate) IsActiveAt(now time.Time) bool {
	if f == nil || !f.Status {
		return false
	}
	if f.StartsAt != nil && now.Before(startOfDay(*f.StartsAt)) {
		return false
	}
	if f.EndsAt != nil && !now.Before(startOfDay(*f.EndsAt).AddDate(0, 0, 1)) {
		// ends_at is inclusive
		return false
	}
	return true
}

// HasFreeShipping reports whether the threshold rule is configured.
func (f *FeeRate) HasFreeShipping() bool {
	return f != nil && f.Kind == FeeKindShipping && f.FreeThreshold != nil && *f.FreeThreshold > 0
}

// QuoteOrderFees splits one order total into its payable parts from the
// resolved tariff rows (pure function — unit tested):
//
//	grand = items + shipping + packaging
//
// shipping drops to 0 when the configured threshold is met.
func QuoteOrderFees(itemsTotal uint, ship, pack *FeeRate) (shipping, packaging, grand uint, freeShipping bool) {
	if ship != nil && ship.Status {
		shipping = ship.Amount
		if ship.HasFreeShipping() && itemsTotal >= *ship.FreeThreshold {
			shipping, freeShipping = 0, true
		}
	}
	if pack != nil && pack.Status {
		packaging = pack.Amount
	}
	grand = itemsTotal + shipping + packaging
	return shipping, packaging, grand, freeShipping
}
