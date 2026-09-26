package requests

import (
	"net/url"
	"strings"
	"time"

	"shop/domain/entities"
)

// CreateProductSliderRequest is the homepage product-slider form
// (create + edit): title, position slot, optional catalog scope and up to
// entities.MaxSliderProducts products.
type CreateProductSliderRequest struct {
	Title      string `form:"title"`
	Slug       string `form:"slug"`
	Position   string `form:"position"`
	CategoryID string `form:"category_id"`
	Status     string `form:"status"`
	StartsAt   string `form:"starts_at"`
	EndsAt     string `form:"ends_at"`
	ProductIDs []uint `form:"-"` // parsed from products[] by the handler
	OldIDs     []uint `form:"-"` // ids currently attached (edit form)
}

// ValidateProductSliderForm returns Persian messages keyed for the template.
func ValidateProductSliderForm(form url.Values) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(form.Get("title")) == "" {
		errs["title"] = "عنوان اسلایدر الزامی است."
	}

	position := strings.TrimSpace(form.Get("position"))
	if position == "" {
		errs["position"] = "جایگاه اسلایدر الزامی است."
	} else if !entities.ValidSliderPosition(position) {
		errs["position"] = "جایگاه اسلایدر نامعتبر است."
	}

	status := strings.TrimSpace(form.Get("status"))
	if !entities.ValidProductStatus(status) {
		errs["status"] = "وضعیت اسلایدر نامعتبر است."
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

	if len(ParseSliderProductIDs(form)) == 0 {
		errs["products"] = "حداقل یک محصول انتخاب کنید."
	}

	return errs
}

// ParseSliderProductIDs reads the posted products[] list (deduplicated, capped).
func ParseSliderProductIDs(form url.Values) []uint {
	seen := map[uint]struct{}{}
	out := make([]uint, 0, entities.MaxSliderProducts)
	for _, raw := range form["products[]"] {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		id, err := strconvParseUint(raw)
		if err != nil || id == 0 {
			continue
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
		if len(out) >= entities.MaxSliderProducts {
			break
		}
	}
	return out
}

// SliderDates parses the optional schedule window.
func SliderDates(startsAt, endsAt string) (start, end *time.Time) {
	return parseDate(startsAt), parseDate(endsAt)
}

// SliderStatus returns the validated status (defaults to published).
func SliderStatus(form url.Values) string {
	status := strings.TrimSpace(form.Get("status"))
	if entities.ValidProductStatus(status) {
		return status
	}
	return entities.ProductStatusPublished
}

// SliderCategoryID returns the parsed catalog scope (0 = none).
func SliderCategoryID(form url.Values) uint {
	id, _ := strconvParseUint(strings.TrimSpace(form.Get("category_id")))
	return id
}

// CatalogCategoryID is the bound-request form of SliderCategoryID.
func (r *CreateProductSliderRequest) CatalogCategoryID() uint {
	id, _ := strconvParseUint(strings.TrimSpace(r.CategoryID))
	return id
}
