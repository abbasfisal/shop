package attributeValue

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
	"strings"
)

type AttributeValueRepository struct {
	db *gorm.DB
}

func NewAttributeRepository(db *gorm.DB) repositories.AttributeValueRepositoryInterface {
	return &AttributeValueRepository{db: db}
}

func (ar *AttributeValueRepository) Store(ctx context.Context, req *requests.CreateAttributeValueRequest) (*entities.AttributeValue, error) {
	var attribute entities.Attribute
	err := ar.db.WithContext(ctx).First(&attribute, req.AttributeID).Error
	if err != nil {
		return nil, err
	}

	attr := &entities.AttributeValue{
		AttributeID:    req.AttributeID,
		AttributeTitle: attribute.Title,
		Value:          strings.TrimSpace(req.Value),
	}
	if attribute.IsColor() {
		if hex := strings.TrimSpace(req.ColorHex); entities.ValidHexColor(hex) {
			attr.SetColorHex(hex)
		}
	}

	if attErr := ar.db.Create(&attr).Error; attErr != nil {
		return nil, attErr
	}
	return attr, nil
}

func (ar *AttributeValueRepository) GetAllAttribute(c *gin.Context) ([]*entities.Attribute, error) {
	var attributes []*entities.Attribute
	err := ar.db.WithContext(c).Preload("AttributeValues").Find(&attributes).Error
	return attributes, err
}

func (ar *AttributeValueRepository) Find(c *gin.Context, attributeValueID int) (*entities.AttributeValue, error) {
	var attValue entities.AttributeValue
	err := ar.db.WithContext(c).Preload("Attribute").First(&attValue, attributeValueID).Error

	return &attValue, err
}

func (ar *AttributeValueRepository) Update(c *gin.Context, attributeValueID int, req *requests.UpdateAttributeValueRequest) (*entities.AttributeValue, error) {
	var attributeValue entities.AttributeValue

	err := ar.db.WithContext(c).Preload("Attribute").First(&attributeValue, attributeValueID).Error
	if err != nil {
		return nil, err
	}

	meta := attributeValue.Meta
	if attributeValue.Attribute.IsColor() {
		if hex := strings.TrimSpace(req.ColorHex); entities.ValidHexColor(hex) {
			fresh := &entities.AttributeValue{}
			fresh.SetColorHex(hex)
			meta = fresh.Meta
		}
	} else {
		// a value moved off a color attribute must not keep a stale hex
		meta = datatypes.JSON([]byte("{}"))
	}

	updateErr := ar.db.
		Model(&attributeValue).
		Update("attribute_id", req.AttributeID).
		Update("attribute_title", attributeValue.Attribute.Title).
		Update("value", strings.TrimSpace(req.Value)).
		Update("meta", meta).
		Error

	if updateErr != nil {
		return nil, updateErr
	}

	return &attributeValue, nil
}
