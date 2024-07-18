package attribute

import (
	"context"
	"gorm.io/gorm"
	"shop/internal/entities"
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
