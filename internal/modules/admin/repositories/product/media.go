package product

import (
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
)

func (p *ProductRepository) GetImage(c *gin.Context, imageID int) (*entities.ProductImages, error) {
	var image entities.ProductImages
	err := p.db.WithContext(c).Find(&image, imageID).Error
	return &image, err
}

func (p *ProductRepository) DeleteImage(c *gin.Context, imageID int) error {

	var productImage entities.ProductImages
	if imgErr := p.db.WithContext(c).Where("id=?", imageID).First(&productImage).Error; imgErr != nil {
		return imgErr
	}
	if delImgErr := p.db.WithContext(c).Unscoped().Delete(&productImage).Error; delImgErr != nil {
		return delImgErr
	}

	_ = SyncMongo(c, p.db, productImage.ProductID)

	return nil
}

func (p *ProductRepository) StoreImages(c *gin.Context, productID int, imageStoredPath []string) error {
	var images []entities.ProductImages
	for _, image := range imageStoredPath {
		images = append(images, entities.ProductImages{
			ProductID: uint(productID),
			Path:      image,
		},
		)
	}

	if storeImgErr := p.db.WithContext(c).Create(&images).Error; storeImgErr != nil {
		return storeImgErr
	}

	_ = SyncMongo(c, p.db, uint(productID))
	return nil
}
