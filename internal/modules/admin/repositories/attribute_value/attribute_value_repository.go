package attributeValue

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/internal/entities"
)

type AttributeValueRepository struct {
	db *gorm.DB
}

func NewAttributeRepository(db *gorm.DB) AttributeValueRepository {
	return AttributeValueRepository{db: db}
}

func (ar AttributeValueRepository) Store(ctx context.Context, attr entities.AttributeValue) (entities.AttributeValue, error) {

	var attribute entities.Attribute
	err := ar.db.First(&attribute, attr.AttributeID).Error
	if err != nil {
		return entities.AttributeValue{}, err
	}

	attr.AttributeTitle = attribute.Title
	attErr := ar.db.Create(&attr).Error

	return attr, attErr
}

func (ar AttributeValueRepository) GetAllAttribute(c *gin.Context) ([]entities.Attribute, error) {
	var attributes []entities.Attribute
	err := ar.db.WithContext(c).Preload("AttributeValues").Find(&attributes).Error
	return attributes, err
}

func (ar AttributeValueRepository) Find(c *gin.Context, attributeValueID int) (entities.AttributeValue, error) {
	var attValue entities.AttributeValue
	err := ar.db.First(&attValue, attributeValueID).Error

	return attValue, err
}
