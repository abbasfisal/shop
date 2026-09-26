package requests

import (
	"net/url"
	"strings"

	"shop/domain/entities"
)

// CreateFeeRateRequest is the order-fee tariff form (create + edit).
// Amount/threshold stay strings here so "missing" and "zero" are
// distinguishable during validation (a 0 shipping fee is legitimate).
type CreateFeeRateRequest struct {
	Title         string `form:"title"`
	Amount        string `form:"amount"`
	FreeThreshold string `form:"free_threshold"`
	Status        string `form:"status"` // "1" = فعال, "0" = غیرفعال
	StartsAt      string `form:"starts_at"`
	EndsAt        string `form:"ends_at"`
}

// ValidateFeeRateForm returns Persian messages keyed for the fee form
// ({{.ERRORS.<key>}}).
func ValidateFeeRateForm(form url.Values, kind string) map[string]string {
	errs := map[string]string{}

	if !entities.ValidFeeKind(kind) {
		errs["kind"] = "نوع تعرفه نامعتبر است."
		return errs
	}

	if strings.TrimSpace(form.Get("title")) == "" {
		errs["title"] = "عنوان تعرفه الزامی است."
	}

	amount := strings.TrimSpace(form.Get("amount"))
	if amount == "" {
		errs["amount"] = "مبلغ الزامی است."
	} else if !Numeric(amount) {
		errs["amount"] = "مبلغ باید عددی باشد."
	}

	if v := strings.TrimSpace(form.Get("free_threshold")); v != "" && !Numeric(v) {
		errs["free_threshold"] = "سقف ارسال رایگان باید عددی باشد."
	}

	if status := strings.TrimSpace(form.Get("status")); status != "0" && status != "1" {
		errs["status"] = "وضعیت تعرفه نامعتبر است."
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

// ToFeeRate maps a validated request onto the entity.
func (r *CreateFeeRateRequest) ToFeeRate(kind string) *entities.FeeRate {
	startsAt, endsAt := parseDate(r.StartsAt), parseDate(r.EndsAt)
	rate := &entities.FeeRate{
		Kind:     kind,
		Title:    strings.TrimSpace(r.Title),
		Amount:   parseUint(r.Amount),
		Status:   strings.TrimSpace(r.Status) == "1",
		StartsAt: startsAt,
		EndsAt:   endsAt,
	}
	if threshold := strings.TrimSpace(r.FreeThreshold); threshold != "" {
		t := parseUint(threshold)
		rate.FreeThreshold = &t
	}
	return rate
}
