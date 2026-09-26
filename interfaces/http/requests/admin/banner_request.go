package requests

import (
	"net/url"
	"strings"
	"time"

	"shop/domain/entities"
)

// CreateBannerRequest is the promotion banner form (create + edit).
type CreateBannerRequest struct {
	Title     string `form:"title"`
	Layout    string `form:"layout"`
	Link      string `form:"link"`
	Status    string `form:"status"` // "1" = فعال, "0" = غیرفعال
	StartsAt  string `form:"starts_at"`
	EndsAt    string `form:"ends_at"`
	SortOrder int    `form:"sort_order"`

	// legacy placement columns (kept for old rows, not part of the form)
	Type     uint `form:"type"`
	Priority uint `form:"priority"`

	// uploaded file name (set by the handler, not posted)
	BannerImage string
	// existing image on edit (hidden field in the form)
	CurrentImage string `form:"current_image"`
}

// ValidateBannerForm returns Persian validation messages keyed the way the
// banner form reads them ({{.ERRORS.<key>}}).
func ValidateBannerForm(form url.Values) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(form.Get("title")) == "" {
		errs["title"] = "عنوان بنر الزامی است."
	}

	layout := NormalizeBannerLayout(form.Get("layout"))
	if !entities.ValidBannerLayout(layout) {
		errs["layout"] = "چیدمان بنر نامعتبر است (دوتایی یا چهار تایی)."
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
	if start, end := strings.TrimSpace(form.Get("starts_at")), strings.TrimSpace(form.Get("ends_at")); start != "" && end != "" {
		s, e := parseDate(start), parseDate(end)
		if s != nil && e != nil && e.Before(*s) {
			errs["ends_at"] = "تاریخ پایان نمی‌تواند قبل از تاریخ شروع باشد."
		}
	}

	return errs
}

// NormalizeBannerLayout defaults an empty layout to the two-up shape.
func NormalizeBannerLayout(layout string) string {
	layout = strings.TrimSpace(layout)
	if layout == "" {
		return entities.BannerLayoutTwo
	}
	return layout
}

// BannerStatusValue turns the posted select ("1"/"0") into the boolean column.
func BannerStatusValue(form url.Values) bool {
	return strings.TrimSpace(form.Get("status")) == "1"
}

// IsActive is the bound-request form of BannerStatusValue.
func (r *CreateBannerRequest) IsActive() bool {
	return strings.TrimSpace(r.Status) == "1"
}

// LayoutValue is the bound-request form of NormalizeBannerLayout.
func (r *CreateBannerRequest) LayoutValue() string {
	return NormalizeBannerLayout(r.Layout)
}

// BannerDates parses the optional schedule window (both bounds optional).
func BannerDates(startsAt, endsAt string) (start, end *time.Time) {
	return parseDate(startsAt), parseDate(endsAt)
}
