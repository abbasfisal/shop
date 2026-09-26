package responses

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"shop/domain/entities"
)

// VariantAttributeValue is one attribute value linked to a variant
// (used by the admin edit page combination table).
type VariantAttributeValue struct {
	ID               uint
	AttributeValueID uint
	AttributeID      uint
	AttributeTitle   string
	Value            string
}

type ProductInventory struct {
	ID        uint
	ProductID uint
	Quantity  uint // stock
	Reserved  uint
	Available uint
	Sku       string

	Price           uint
	SalePrice       uint
	DiscountPrice   *uint
	EffectivePrice  uint
	HasDiscount     bool
	DiscountPercent int

	Status     string
	StatusText string
	ExpiresAt  *time.Time

	AttributeValues []VariantAttributeValue
	// Label is the human readable combination ("قرمز / M")
	Label string
	// ValueIDsKey is the sorted attribute_value id list ("1,2") — used by the
	// admin edit template to detect combinations that already exist.
	ValueIDsKey string
}

type ProductInventories struct {
	Data []ProductInventory
}

// ToProductInventory maps one variant with its pricing/status helpers.
// fallbackOriginal/fallbackSale are the parent product prices — a NULL variant
// price column inherits them (see PricingService).
func ToProductInventory(pv *entities.ProductVariant, fallbackOriginal, fallbackSale uint) *ProductInventory {
	original := fallbackOriginal
	if pv.Price != nil {
		original = *pv.Price
	}
	sale := fallbackSale
	if pv.SalePrice != nil {
		sale = *pv.SalePrice
	}

	// golden rule: discount counts only when > 0 and < effective sale price.
	// The badge percent is measured against the crossed-out list price
	// (original), consistent with the storefront display.
	effective := sale
	hasDiscount := false
	discountPercent := 0
	if pv.DiscountPrice != nil && *pv.DiscountPrice > 0 && *pv.DiscountPrice < sale {
		effective = *pv.DiscountPrice
		hasDiscount = true
		if original > 0 && effective < original {
			discountPercent = int(math.Round(float64(original-effective) / float64(original) * 100))
		}
	}

	out := &ProductInventory{
		ID:        pv.ID,
		ProductID: pv.ProductID,
		Quantity:  pv.Stock,
		Reserved:  pv.ReservedStock,
		Available: availableOf(pv),

		Price:           original,
		SalePrice:       sale,
		DiscountPrice:   pv.DiscountPrice,
		EffectivePrice:  effective,
		HasDiscount:     hasDiscount,
		DiscountPercent: discountPercent,

		Status:     pv.Status,
		StatusText: variantStatusLabel(pv.Status),
		ExpiresAt:  pv.ExpiresAt,
	}

	labels := make([]string, 0, len(pv.VariantAttributeValues))
	valueIDs := make([]uint, 0, len(pv.VariantAttributeValues))
	for _, link := range pv.VariantAttributeValues {
		valueID := uint(0)
		attrID := uint(0)
		attrTitle := ""
		value := ""
		if link.AttributeValue != nil {
			valueID = link.AttributeValue.ID
			attrID = link.AttributeValue.AttributeID
			attrTitle = link.AttributeValue.AttributeTitle
			value = link.AttributeValue.Value
		}
		out.AttributeValues = append(out.AttributeValues, VariantAttributeValue{
			ID:               link.ID,
			AttributeValueID: valueID,
			AttributeID:      attrID,
			AttributeTitle:   attrTitle,
			Value:            value,
		})
		if valueID > 0 {
			valueIDs = append(valueIDs, valueID)
		}
		if value != "" {
			labels = append(labels, value)
		}
	}
	out.Label = strings.Join(labels, " / ")

	sort.Slice(valueIDs, func(i, j int) bool { return valueIDs[i] < valueIDs[j] })
	ids := make([]string, 0, len(valueIDs))
	for _, id := range valueIDs {
		ids = append(ids, strconv.FormatUint(uint64(id), 10))
	}
	out.ValueIDsKey = strings.Join(ids, ",")

	return out
}

// ToProductInventories maps every variant of a product (no price fallback).
func ToProductInventories(variants []*entities.ProductVariant) *ProductInventories {
	return ToProductInventoriesWithFallback(variants, 0, 0)
}

// ToProductInventoriesWithFallback maps variants inheriting the product prices.
func ToProductInventoriesWithFallback(variants []*entities.ProductVariant, fallbackOriginal, fallbackSale uint) *ProductInventories {
	var pResponse ProductInventories
	for _, pv := range variants {
		pResponse.Data = append(pResponse.Data, *ToProductInventory(pv, fallbackOriginal, fallbackSale))
	}
	return &pResponse
}

func availableOf(pv *entities.ProductVariant) uint {
	if pv.Stock > pv.ReservedStock {
		return pv.Stock - pv.ReservedStock
	}
	return 0
}

func variantStatusLabel(status string) string {
	if status == entities.VariantStatusInactive {
		return "غیرفعال"
	}
	return "فعال"
}
