package requests

import (
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// VariantRow is one sellable combination posted by the product create/edit
// form. Field names mirror the Laravel form (variants[i][price], ...).
type VariantRow struct {
	Index             int    `json:"index"` // posting index (variants[N])
	ID                uint   `json:"id"`    // > 0 => existing variant (edit form)
	Sku               string `json:"sku"`
	AttributeValueIDs []uint `json:"attribute_value_ids"`
	Price             uint   `json:"price"`
	SalePrice         uint   `json:"sale_price"`
	DiscountPrice     *uint  `json:"discount_price"`
	Stock             uint   `json:"stock"`
	Status            string `json:"status"`
	ExpiresAt         string `json:"expires_at"`
}

// variantKeyRe matches  variants[0][price] / variants[0][attribute_value_ids][] / ...
var variantKeyRe = regexp.MustCompile(`^variants\[(\d+)\]\[([a-z0-9_]+)\](\[\])?$`)

// ParseVariantRows reads the raw url-encoded form and builds ordered variant
// rows. Go's form parser keeps bracket keys verbatim, so the Laravel-style
// naming works without a custom decoder.
func ParseVariantRows(form url.Values) []VariantRow {
	rows := map[int]*VariantRow{}

	for key, vals := range form {
		m := variantKeyRe.FindStringSubmatch(key)
		if m == nil {
			continue
		}
		idx, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		field := m[2]

		row, ok := rows[idx]
		if !ok {
			row = &VariantRow{Index: idx, Status: "active"}
			rows[idx] = row
		}

		switch field {
		case "attribute_value_ids":
			for _, raw := range vals {
				raw = strings.TrimSpace(raw)
				if raw == "" {
					continue
				}
				if id, err := strconv.ParseUint(raw, 10, 64); err == nil && id > 0 {
					row.AttributeValueIDs = append(row.AttributeValueIDs, uint(id))
				}
			}
		default:
			v := ""
			if len(vals) > 0 {
				v = strings.TrimSpace(vals[0])
			}
			switch field {
			case "id":
				if n, err := strconv.ParseUint(v, 10, 64); err == nil {
					row.ID = uint(n)
				}
			case "sku":
				row.Sku = v
			case "price":
				row.Price = parseUint(v)
			case "sale_price":
				row.SalePrice = parseUint(v)
			case "discount_price":
				if v != "" {
					n := parseUint(v)
					row.DiscountPrice = &n
				}
			case "stock":
				row.Stock = parseUint(v)
			case "status":
				if v != "" {
					row.Status = v
				}
			case "expires_at":
				row.ExpiresAt = v
			}
		}
	}

	out := make([]VariantRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, *r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Index < out[j].Index })
	return out
}

// parseDate parses an HTML date input (YYYY-MM-DD).
func parseDate(v string) *time.Time {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	for _, layout := range []string{"2006-01-02", "2006/01/02"} {
		if t, err := time.Parse(layout, v); err == nil {
			return &t
		}
	}
	return nil
}

func parseUint(v string) uint {
	if v == "" {
		return 0
	}
	n, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0
	}
	return uint(n)
}

// Numeric reports whether the raw input is a valid non-negative integer
// (validation needs to distinguish "missing/invalid" from zero).
func Numeric(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" {
		return false
	}
	_, err := strconv.ParseUint(v, 10, 64)
	return err == nil
}

// strconvParseUint is a tiny helper shared by the admin form parsers.
func strconvParseUint(v string) (uint, error) {
	n, err := strconv.ParseUint(v, 10, 64)
	return uint(n), err
}
