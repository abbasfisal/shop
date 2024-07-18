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

	ar.db.Where("id=", attr.AttributeID).Find(&attribute)

	attr.AttributeTitle = attribute.Title
	err := ar.db.Create(&attr).Error

	return attr, err
}
