package category

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"shop/internal/entities"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return CategoryRepository{db: db}
}

func (cr CategoryRepository) GetAll(ctx context.Context) ([]entities.Category, error) {
	var categories []entities.Category
	err := cr.db.WithContext(ctx).Find(&categories).Error
	return categories, err
}

func (cr CategoryRepository) GetAllParent(ctx context.Context) ([]entities.Category, error) {
	var categories []entities.Category
	err := cr.db.WithContext(ctx).Where("parent_id IS NULL").Find(&categories).Error
	return categories, err
}

func (cr CategoryRepository) SelectBy(ctx context.Context, categoryID int) (entities.Category, error) {
	var category entities.Category
	err := cr.db.WithContext(ctx).First(&category, "id = ?", categoryID).Error
	return category, err
}

func (cr CategoryRepository) FindBy(ctx context.Context, columnName string, value any) (entities.Category, error) {
	var category entities.Category
	condition := fmt.Sprintf("%s = ?", columnName)
	err := cr.db.WithContext(ctx).First(&category, condition, value).Error
	return category, err
}

func (cr CategoryRepository) Store(ctx context.Context, category entities.Category) (entities.Category, error) {
	err := cr.db.WithContext(ctx).Create(&category).Error
	return category, err
}
