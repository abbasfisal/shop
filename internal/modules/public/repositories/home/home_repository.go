package home

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"shop/internal/database/mysql"
	"shop/internal/entities"
)

type HomeRepository struct {
	db *gorm.DB
}

func NewHomeRepository() HomeRepository {
	return HomeRepository{
		db: mysql.Get(),
	}
}
func (h HomeRepository) GetRandomProducts(ctx context.Context, limit int) ([]entities.Product, error) {
	var products []entities.Product
	//implement Me
	return products, nil
}

func (h HomeRepository) GetLatestProducts(ctx context.Context, limit int) ([]entities.Product, error) {
	var products []entities.Product
	//todo: just load data if category.status = true and product.status=true
	err := h.db.Preload("Category").Where("status=?", true).Limit(limit).Find(&products).Error
	return products, err
}

func (h HomeRepository) GetCategories(ctx context.Context, limit int) ([]entities.Category, error) {
	var categories []entities.Category
	err := h.db.Limit(limit).Find(&categories, "status=?", true).Error

	return categories, err
}
func (h HomeRepository) GetProduct(ctx context.Context, productSlug, sku string) (entities.Product, error) {
	var product entities.Product
	err := h.db.Where("slug=? and sku=? and status=true", productSlug, sku).First(&product).Error
	return product, err
}

func (h HomeRepository) GetProductsBy(ctx context.Context, columnName string, value any) ([]entities.Product, error) {
	var products []entities.Product
	condition := fmt.Sprintf("%s = ?", columnName)
	err := h.db.Where(condition, value).Find(&products).Error

	return products, err
}

func (h HomeRepository) GetCategoryBy(ctx context.Context, columnName string, value any) (entities.Category, error) {
	var category entities.Category
	err := h.db.Where(fmt.Sprintf("%s = ?", columnName), value).Find(&category).Error

	return category, err
}
