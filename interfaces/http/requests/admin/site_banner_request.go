package requests

import (
	"net/url"
	"strings"
	"time"

	"shop/domain/entities"
)

// CreateSiteBannerRequest is the fixed site-banner form (create + edit):
// header banner today, the reserved image-hero placement tomorrow.
type CreateSiteBannerRequest struct {
	Title     string `form:"title"`
	Placement string `form:"placement"`
	Link      string `form:"link"`
	Status    string `form:"status"` // "1" = فعال, "0" = غیرفعال
	StartsAt  string `form:"starts_at"`
	EndsAt    string `form:"ends_at"`
	SortOrder int    `form:"sort_order"`

	// uploaded file names (set by the handler, not posted)
	SiteBannerImage string
	MobileImage     string
	// clearing the mobile image is explicit (checkbox), not "no file uploaded"
	RemoveMobileImage bool `form:"remove_mobile_image"`

	// existing rows on edit (hidden fields in the form)
	CurrentImage       string `form:"current_image"`
	CurrentMobileImage string `form:"current_mobile_image"`
}

// ValidateSiteBannerForm returns Persian validation messages keyed the way the
// form reads them ({{.ERRORS.<key>}}).
func ValidateSiteBannerForm(form url.Values) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(form.Get("title")) == "" {
		errs["title"] = "عنوان بنر الزامی است."
	}

	placement := entities.NormalizeSiteBannerPlacement(strings.TrimSpace(form.Get("placement")))
	if !entities.ValidSiteBannerPlacement(placement) {
		errs["placement"] = "جایگاه بنر نامعتبر است."
	}

	status := strings.TrimSpace(form.Get("status"))
	if status != "0" && status != "1" {
		errs["status"] = "وضعیت بنر نامعتبر است."
	}

	if v := strings.TrimSpace(form.Get("starts_at")); v != "" && parseDate(v) == nil {
		errs["starts_at"] = "تاریخ شروع معتبر نیست."
	}
	if v := strings.TrimSpace(form.Get("ends_at")); v != "" && parseDate(v) == nil {
		errs["ends_at"] = "تاریخ پایان معتبر نیست."
	}
	startRaw := strings.TrimSpace(form.Get("starts_at"))
	endRaw := strings.TrimSpace(form.Get("ends_at"))
	if startRaw != "" && endRaw != "" {
		s, e := parseDate(startRaw), parseDate(endRaw)
		if s != nil && e != nil && e.Before(*s) {
			errs["ends_at"] = "تاریخ پایان نمی‌تواند قبل از تاریخ شروع باشد."
		}
	}

	return errs
}

// IsActive is the bound-request form of the posted select ("1"/"0").
func (r *CreateSiteBannerRequest) IsActive() bool {
	return strings.TrimSpace(r.Status) == "1"
}

// PlacementValue is the bound-request form of NormalizeSiteBannerPlacement.
func (r *CreateSiteBannerRequest) PlacementValue() string {
	return entities.NormalizeSiteBannerPlacement(strings.TrimSpace(r.Placement))
}

// SiteBannerDates parses the optional schedule window (both bounds optional).
func SiteBannerDates(startsAt, endsAt string) (start, end *time.Time) {
	return parseDate(startsAt), parseDate(endsAt)
}
