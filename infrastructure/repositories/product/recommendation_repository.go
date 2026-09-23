package product

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/domain/entities"
)

// GetAllProductBriefs returns lightweight product summaries used by the admin
// recommendation picker. Shape is required by admin_edit_product.html:
//
//	$product.id            -> posted back as recommendation id (string product id)
//	$product.product.id    -> numeric id, compared against recommendations
//	$product.product.title -> display title
func (p *ProductRepository) GetAllProductBriefs(c context.Context) ([]map[string]interface{}, error) {
	var products []entities.Product
	err := p.db.WithContext(c).
		Select("id, title, slug, original_price, sale_price").
		Order("id").
		Find(&products).
		Error
	if err != nil {
		return nil, err
	}

	briefs := make([]map[string]interface{}, 0, len(products))
	for _, prod := range products {
		briefs = append(briefs, map[string]interface{}{
			"id": strconv.FormatUint(uint64(prod.ID), 10),
			"product": map[string]interface{}{
				"id":             int64(prod.ID),
				"title":          prod.Title,
				"slug":           prod.Slug,
				"original_price": int64(prod.OriginalPrice),
				"sale_price":     int64(prod.SalePrice),
			},
		})
	}
	return briefs, nil
}

// InsertRecommendation replaces the recommendation set of a product.
// ids are string product ids coming from the admin form.
func (p *ProductRepository) InsertRecommendation(c *gin.Context, productID int, productRecommendationIDs []string) error {
	return p.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("product_id = ?", productID).
			Delete(&entities.ProductRecommendation{}).
			Error; err != nil {
			return err
		}

		// resolve requested ids, dropping self references and ids that do not
		// exist (the FK on recommended_product_id would otherwise abort the tx)
		var existingIDs []uint
		if err := tx.Model(&entities.Product{}).
			Where("id IN ?", parsedIDs(productRecommendationIDs)).
			Pluck("id", &existingIDs).
			Error; err != nil {
			return err
		}
		valid := make(map[uint]struct{}, len(existingIDs))
		for _, id := range existingIDs {
			if id == uint(productID) {
				continue // no self recommendation
			}
			valid[id] = struct{}{}
		}

		for id := range valid {
			if err := tx.Create(&entities.ProductRecommendation{
				ProductID:            uint(productID),
				RecommendedProductID: id,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetAllRecommendation returns the recommendation products of a product in the
// same shape as GetAllProductBriefs (consumed by admin_edit_product.html).
func (p *ProductRepository) GetAllRecommendation(c *gin.Context, productID int) ([]map[string]interface{}, error) {
	var recs []entities.ProductRecommendation
	err := p.db.WithContext(c).
		Where("product_id = ?", productID).
		Find(&recs).
		Error
	if err != nil {
		return nil, err
	}
	if len(recs) == 0 {
		return nil, nil
	}

	ids := make([]uint, 0, len(recs))
	for _, r := range recs {
		ids = append(ids, r.RecommendedProductID)
	}

	var products []entities.Product
	err = p.db.WithContext(c).
		Select("id, title, slug, original_price, sale_price").
		Where("id IN ?", ids).
		Find(&products).
		Error
	if err != nil {
		return nil, err
	}

	briefs := make([]map[string]interface{}, 0, len(products))
	for _, prod := range products {
		briefs = append(briefs, map[string]interface{}{
			"id": strconv.FormatUint(uint64(prod.ID), 10),
			"product": map[string]interface{}{
				"id":             int64(prod.ID),
				"title":          prod.Title,
				"slug":           prod.Slug,
				"original_price": int64(prod.OriginalPrice),
				"sale_price":     int64(prod.SalePrice),
			},
		})
	}
	return briefs, nil
}

// parsedIDs converts string form ids to uints, skipping invalid entries.
func parsedIDs(ids []string) []uint {
	out := make([]uint, 0, len(ids))
	for _, s := range ids {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil || v == 0 {
			continue
		}
		out = append(out, uint(v))
	}
	return out
}
