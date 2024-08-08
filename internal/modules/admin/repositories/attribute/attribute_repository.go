package attribute

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
	"strings"
)

type AttributeRepository struct {
	db *gorm.DB
}

func NewAttributeRepository(db *gorm.DB) AttributeRepository {
	return AttributeRepository{db: db}
}

func (ar AttributeRepository) Store(ctx context.Context, attr entities.Attribute) (entities.Attribute, error) {
	err := ar.db.Create(&attr).Error
	return attr, err
}

func (ar AttributeRepository) GetByCategory(ctx context.Context, catID int) ([]entities.Attribute, error) {
	var attributes []entities.Attribute

	err := ar.db.Where("category_id = ? ", catID).Find(&attributes).Error
	return attributes, err
}

func (ar AttributeRepository) GetAll(c *gin.Context) ([]entities.Attribute, error) {
	var attributes []entities.Attribute
	err := ar.db.WithContext(c).Find(&attributes).Error
	return attributes, err

}

func (ar AttributeRepository) GetByID(c context.Context, attributeID int) (entities.Attribute, error) {
	var att entities.Attribute
	err := ar.db.WithContext(c).Preload("AttributeValues").First(&att, attributeID).Error
	return att, err
}

func (ar AttributeRepository) Update(c *gin.Context, attributeID int, req requests.CreateAttributeRequest) error {

	var att entities.Attribute
	err := ar.db.First(&att, attributeID).Error
	if err != nil {
		return err
	}

	uErr := ar.db.Model(&att).Update("title", strings.TrimSpace(req.Title)).Error

	return uErr
}
