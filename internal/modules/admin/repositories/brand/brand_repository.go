package brand

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"shop/internal/entities"
)

type BrandRepository struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) BrandRepository {
	return BrandRepository{db: db}
}

func (br BrandRepository) FindBy(ctx context.Context, columnName string, value any) (entities.Brand, error) {
	var brand entities.Brand
	condition := fmt.Sprintf("%s = ?", columnName)
	err := br.db.First(&brand, condition, value).Error
	return brand, err
}

func (br BrandRepository) Store(ctx context.Context, brand entities.Brand) (entities.Brand, error) {
	err := br.db.Create(&brand).Error
	return brand, err
}
