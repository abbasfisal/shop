package attributeValue

import (
	"context"
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
