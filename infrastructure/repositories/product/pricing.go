package product

import (
	"context"
	"encoding/json"
	"log"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"shop/domain/entities"
)

// RefreshProductAggregates recomputes the aggregate cache columns on the
// products row from its active variants (Laravel PricingService pattern):
//
//   - min_price / max_price  = range of effective prices
//     effective = discount_price if it is a real discount, else
//     variant.sale_price when set, else products.sale_price
//   - total_stock / total_reserved / available_stock / in_stock
//   - variants_count
//   - product_type           = 'variable' when any variant has attribute values
//   - attributes_json (JSONB) = {attribute_code: [attribute_value_ids]} (GIN indexed)
//
// Repository implementation (called through application/usecases/pricing.PricingService
// and directly from infrastructure flows such as order payment).
func RefreshProductAggregates(ctx context.Context, db *gorm.DB, productID uint) error {
	var prod entities.Product
	if err := db.WithContext(ctx).
		Select("id, original_price, sale_price").
		First(&prod, productID).
		Error; err != nil {
		return err
	}

	type variantAgg struct {
		VariantsCount int
		TotalStock    uint
		TotalReserved uint
		MinPrice      *int64
		MaxPrice      *int64
	}

	var agg variantAgg
	aggSQL := `
		SELECT COUNT(*)::int                                    AS variants_count,
		       COALESCE(SUM(stock), 0)                          AS total_stock,
		       COALESCE(SUM(reserved_stock), 0)                 AS total_reserved,
		       COALESCE(MIN(eff), ?)::bigint                    AS min_price,
		       COALESCE(MAX(eff), ?)::bigint                    AS max_price
		FROM (
			SELECT v.stock,
			       v.reserved_stock,
			       CASE
			           WHEN v.discount_price IS NOT NULL AND v.discount_price > 0
			               AND v.discount_price < COALESCE(NULLIF(v.sale_price, 0), ?)
			               THEN v.discount_price
			           ELSE COALESCE(NULLIF(v.sale_price, 0), ?)
			       END AS eff
			FROM product_variants v
			         JOIN products p ON p.id = v.product_id
			WHERE v.product_id = ?
			  AND v.deleted_at IS NULL
			  AND v.status = 'active'
			  AND p.deleted_at IS NULL
		) agg`
	fallback := int64(prod.SalePrice)
	if err := db.WithContext(ctx).
		Raw(aggSQL, fallback, fallback, fallback, fallback, productID).
		Scan(&agg).
		Error; err != nil {
		log.Println("[pricing] aggregate scan error:", err)
		return err
	}

	// attributes_json = {code: [value ids]} across the product's active variants
	type attrRow struct {
		Code    string
		ValueID uint
	}
	var attrRows []attrRow
	attrSQL := `
		SELECT COALESCE(NULLIF(a.code, ''), 'attr_' || a.id::text) AS code,
		       vav.attribute_value_id                             AS value_id
		FROM variant_attribute_values vav
		         JOIN product_variants v ON v.id = vav.variant_id
		         JOIN attribute_values av ON av.id = vav.attribute_value_id
		         JOIN attributes a ON a.id = av.attribute_id
		WHERE v.product_id = ?
		  AND v.deleted_at IS NULL
		  AND v.status = 'active'
		  AND vav.deleted_at IS NULL
		  AND av.deleted_at IS NULL
		  AND a.deleted_at IS NULL`
	if err := db.WithContext(ctx).Raw(attrSQL, productID).Scan(&attrRows).Error; err != nil {
		log.Println("[pricing] attributes_json scan error:", err)
		return err
	}

	attrs := map[string][]uint{}
	for _, r := range attrRows {
		ids := attrs[r.Code]
		dup := false
		for _, id := range ids {
			if id == r.ValueID {
				dup = true
				break
			}
		}
		if !dup {
			attrs[r.Code] = append(ids, r.ValueID)
		}
	}
	attrsJSON, err := json.Marshal(attrs)
	if err != nil {
		return err
	}

	// product_type = variable when the product has any variant attribute links
	type typeRow struct {
		HasAttrs bool
	}
	var tr typeRow
	if err := db.WithContext(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM variant_attribute_values vav
			         JOIN product_variants v ON v.id = vav.variant_id
			WHERE v.product_id = ?
			  AND v.deleted_at IS NULL
			  AND v.status = 'active'
			  AND vav.deleted_at IS NULL
		) AS has_attrs`, productID).Scan(&tr).Error; err != nil {
		return err
	}
	productType := "simple"
	if tr.HasAttrs {
		productType = "variable"
	}

	totalStock := agg.TotalStock
	totalReserved := agg.TotalReserved
	available := uint(0)
	if totalStock > totalReserved {
		available = totalStock - totalReserved
	}
	minPrice := uint(0)
	maxPrice := uint(0)
	if agg.MinPrice != nil && *agg.MinPrice > 0 {
		minPrice = uint(*agg.MinPrice)
	}
	if agg.MaxPrice != nil && *agg.MaxPrice > 0 {
		maxPrice = uint(*agg.MaxPrice)
	}

	return db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"min_price":       minPrice,
			"max_price":       maxPrice,
			"total_stock":     totalStock,
			"total_reserved":  totalReserved,
			"available_stock": available,
			"in_stock":        available > 0,
			"variants_count":  agg.VariantsCount,
			"product_type":    productType,
			"attributes_json": datatypes.JSON(attrsJSON),
		}).
		Error
}

// RefreshProductAggregates is the ProductRepositoryInterface implementation.
func (p *ProductRepository) RefreshProductAggregates(ctx context.Context, productID uint) error {
	return RefreshProductAggregates(ctx, p.db, productID)
}
