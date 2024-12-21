package attributeValue

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
	"strings"
)

type AttributeValueRepository struct {
	db *gorm.DB
}

func NewAttributeRepository(db *gorm.DB) AttributeValueRepositoryInterface {
	return &AttributeValueRepository{db: db}
}

func (ar *AttributeValueRepository) Store(ctx context.Context, attr *entities.AttributeValue) (*entities.AttributeValue, error) {

	var attribute entities.Attribute
	err := ar.db.WithContext(ctx).First(&attribute, attr.AttributeID).Error
	if err != nil {
		return nil, err
	}

	attr.AttributeTitle = attribute.Title
	attErr := ar.db.Create(&attr).Error

	return attr, attErr
}

func (ar *AttributeValueRepository) GetAllAttribute(c *gin.Context) ([]*entities.Attribute, error) {
	var attributes []*entities.Attribute
	err := ar.db.WithContext(c).Preload("AttributeValues").Find(&attributes).Error
	return attributes, err
}

func (ar *AttributeValueRepository) Find(c *gin.Context, attributeValueID int) (*entities.AttributeValue, error) {
	var attValue entities.AttributeValue
	err := ar.db.First(&attValue, attributeValueID).Error

	return &attValue, err
}

func (ar *AttributeValueRepository) Update(c *gin.Context, attributeValueID int, req *requests.UpdateAttributeValueRequest) (*entities.AttributeValue, error) {
	var attributeValue entities.AttributeValue

	err := ar.db.WithContext(c).Preload("Attribute").First(&attributeValue, attributeValueID).Error
	if err != nil {
		return nil, err
	}

	updateErr := ar.db.
		Model(&attributeValue).
		Update("attribute_id", req.AttributeID).
		Update("attribute_title", attributeValue.Attribute.Title).
		Update("value", strings.TrimSpace(req.Value)).
		Error

	if updateErr != nil {
		return nil, updateErr
	}

	return &attributeValue, nil
}
