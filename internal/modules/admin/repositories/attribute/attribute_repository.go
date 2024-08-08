package attribute

import (
	"context"
	"github.com/gin-gonic/gin"
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
