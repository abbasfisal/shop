package category

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"shop/internal/database/mysql"
	"shop/internal/entities"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
	return CategoryRepository{db: mysql.Get()}
}

func (cr CategoryRepository) GetAll(ctx context.Context) ([]entities.Category, error) {
	var categories []entities.Category
	err := cr.db.Find(&categories).Error
	return categories, err
}

func (cr CategoryRepository) SelectBy(ctx context.Context, categoryID int) (entities.Category, error) {
	var category entities.Category
	err := cr.db.First(&category, "id=?", categoryID).Error

	return category, err
}

func (cr CategoryRepository) FindBy(ctx context.Context, columnName string, value any) (entities.Category, error) {
	var cat entities.Category
	condition := fmt.Sprintf("%s = ?", columnName)
	err := cr.db.First(&cat, condition, value).Error
	return cat, err
}

func (cr CategoryRepository) Store(ctx context.Context, cat entities.Category) (entities.Category, error) {
	err := cr.db.Create(&cat).Error
	return cat, err
}
