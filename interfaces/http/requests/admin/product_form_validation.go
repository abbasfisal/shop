package requests

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"shop/domain/entities"
)

// ValidateProductForm validates the create/edit product form (Laravel
// StoreProductRequest rules) and returns Persian messages keyed the same way
// the templates read them: {{.ERRORS.<key>}}.
func ValidateProductForm(form url.Values) map[string]string {
	errs := map[string]string{}

	productType := strings.TrimSpace(form.Get("product_type"))
	switch productType {
	case "simple", "variable":
	default:
		errs["product_type"] = "نوع محصول الزامی است (ساده یا واریانت‌دار)."
	}

	status := strings.TrimSpace(form.Get("status"))
	if status == "" {
		errs["status"] = "وضعیت محصول الزامی است."
	} else if !entities.ValidProductStatus(status) {
		errs["status"] = "وضعیت محصول نامعتبر است."
	}

	if v := strings.TrimSpace(form.Get("expires_at")); v != "" && parseDate(v) == nil {
		errs["expires_at"] = "تاریخ انقضای محصول معتبر نیست."
	}

	switch productType {
	case "simple":
		validateSimple(form, errs)
	case "variable":
		validateVariants(form, errs)
	}

	return errs
}

func validateSimple(form url.Values, errs map[string]string) {
	stock := strings.TrimSpace(form.Get("simple_stock"))
	if stock == "" {
		errs["simple_stock"] = "موجودی محصول الزامی است."
	} else if !Numeric(stock) {
		errs["simple_stock"] = "موجودی محصول باید عددی باشد."
	}

	if v := strings.TrimSpace(form.Get("simple_discount_price")); v != "" && !Numeric(v) {
		errs["simple_discount_price"] = "قیمت با تخفیف باید عددی باشد."
	}

	status := strings.TrimSpace(form.Get("simple_status"))
	if status == "" {
		status = entities.VariantStatusActive
	}
	if status != entities.VariantStatusActive && status != entities.VariantStatusInactive {
		errs["simple_status"] = "وضعیت واریانت نامعتبر است."
	}

	if v := strings.TrimSpace(form.Get("simple_expires_at")); v != "" && parseDate(v) == nil {
		errs["simple_expires_at"] = "تاریخ انقضای واریانت معتبر نیست."
	}
}

func validateVariants(form url.Values, errs map[string]string) {
	rows := ParseVariantRows(form)
	if len(rows) == 0 {
		errs["variants"] = "برای محصول واریانت‌دار حداقل یک ترکیب لازم است."
		return
	}

	var messages []string
	for i, row := range rows {
		label := fmt.Sprintf("ردیف %d", i+1)
		var rowErrs []string

		if !Numeric(rawVariantField(form, row.Index, "price")) {
			rowErrs = append(rowErrs, "قیمت اصلی الزامی است")
		}
		if !Numeric(rawVariantField(form, row.Index, "sale_price")) {
			rowErrs = append(rowErrs, "قیمت فروش الزامی است")
		}
		if !Numeric(rawVariantField(form, row.Index, "stock")) {
			rowErrs = append(rowErrs, "موجودی الزامی است")
		}
		if v := rawVariantField(form, row.Index, "discount_price"); v != "" && !Numeric(v) {
			rowErrs = append(rowErrs, "قیمت تخفیف باید عددی باشد")
		}
		switch row.Status {
		case entities.VariantStatusActive, entities.VariantStatusInactive:
		default:
			rowErrs = append(rowErrs, "وضعیت واریانت نامعتبر است")
		}
		// new combinations need at least one attribute value; existing rows
		// (posted with an id) are edited in place and keep their links
		if len(row.AttributeValueIDs) == 0 && row.ID == 0 {
			rowErrs = append(rowErrs, "حداقل یک ویژگی لازم است")
		}
		if row.ExpiresAt != "" && parseDate(row.ExpiresAt) == nil {
			rowErrs = append(rowErrs, "تاریخ انقضا معتبر نیست")
		}

		if len(rowErrs) > 0 {
			messages = append(messages, label+": "+strings.Join(rowErrs, ", "))
		}
	}

	if len(messages) > 0 {
		errs["variants"] = strings.Join(messages, " | ")
	}
}

func rawVariantField(form url.Values, index int, field string) string {
	return strings.TrimSpace(form.Get(fmt.Sprintf("variants[%d][%s]", index, field)))
}

// ExpiresAtTime returns the parsed product-level expiry (nil when empty/invalid).
func (r *CreateProductRequest) ExpiresAtTime() *time.Time {
	return parseDate(r.ExpiresAt)
}

// UpdateExpiresAtTime is the update-form counterpart of ExpiresAtTime.
func (r *UpdateProductRequest) ExpiresAtTime() *time.Time {
	return parseDate(r.ExpiresAt)
}

// SimpleVariantRows builds the single variant row for product_type=simple.
func (r *CreateProductRequest) SimpleVariantRows() []VariantRow {
	return simpleRows(r.OriginalPrice, r.SalePrice, r.SimpleDiscountPrice,
		r.SimpleStock, r.SimpleStatus, r.SimpleExpiresAt)
}

// SimpleVariantRows is the update-form counterpart.
func (r *UpdateProductRequest) SimpleVariantRows() []VariantRow {
	return simpleRows(r.OriginalPrice, r.SalePrice, r.SimpleDiscountPrice,
		r.SimpleStock, r.SimpleStatus, r.SimpleExpiresAt)
}

// VariantRows returns the rows that must be stored for this request:
// variable → posted combination rows, simple → the single derived row.
func (r *CreateProductRequest) VariantRows() []VariantRow {
	if r.ProductType == "variable" {
		return r.Variants
	}
	return r.SimpleVariantRows()
}

// VariantRows is the update-form counterpart (rows with ID > 0 update, others create).
func (r *UpdateProductRequest) VariantRows() []VariantRow {
	if r.ProductType == "variable" {
		return r.Variants
	}
	return r.SimpleVariantRows()
}

func simpleRows(price, salePrice uint, discount *uint, stockRaw, status, expiresAt string) []VariantRow {
	if status == "" {
		status = entities.VariantStatusActive
	}
	row := VariantRow{
		Index:         0,
		Price:         price,
		SalePrice:     salePrice,
		DiscountPrice: discount,
		Stock:         parseUint(stockRaw),
		Status:        status,
		ExpiresAt:     expiresAt,
	}
	return []VariantRow{row}
}
