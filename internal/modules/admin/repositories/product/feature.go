package product

import (
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
)

func (p *ProductRepository) InsertFeature(c *gin.Context, productID int, req *requests.CreateProductFeatureRequest) error {

	if err := p.db.Create(&entities.Feature{
		ProductID: uint(productID),
		Title:     req.Title,
		Value:     req.Value,
	}).Error; err != nil {
		return err
	}

	_ = SyncMongo(c, p.db, uint(productID))

	return nil
}

func (p *ProductRepository) DeleteFeature(c *gin.Context, productID int, featureID int) error {
	if err := p.db.WithContext(c).Where("product_id=? ", productID).Where("id = ?", featureID).Unscoped().Delete(&entities.Feature{}).Error; err != nil {
		return err
	}

	_ = SyncMongo(c, p.db, uint(productID))

	return nil
}

func (p *ProductRepository) GetFeatureBy(c *gin.Context, productID int, featureID int) (*entities.Feature, error) {
	var feature entities.Feature
	if err := p.db.WithContext(c).Where("id=?", featureID).Where("product_id=?", productID).First(&feature).Error; err != nil {
		return nil, err
	}
	return &feature, nil
}

func (p *ProductRepository) EditFeature(c *gin.Context, productID int, featureID int, req *requests.UpdateProductFeatureRequest) error {
	if err := p.db.
		Where("id=?", featureID).
		Where("product_id=?", productID).
		Model(&entities.Feature{}).
		Update("title", req.Title).
		Update("value", req.Value).Error; err != nil {
		return err
	}

	_ = SyncMongo(c, p.db, uint(productID))

	return nil
}
