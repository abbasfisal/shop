package product

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/internal/entities"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return ProductRepository{db: db}
}

func (p ProductRepository) GetAll(ctx context.Context) ([]entities.Product, error) {
	var products []entities.Product
	err := p.db.Preload("Category").Find(&products).Error
	return products, err
}

func (p ProductRepository) FindBy(ctx context.Context, columnName string, value any) (entities.Product, error) {
	var product entities.Product
	condition := fmt.Sprintf("%s=?", columnName)
	err := p.db.Preload("Category").Preload("ProductImage").First(&product, condition, value).Error
	return product, err
}

func (p ProductRepository) FindByID(ctx context.Context, ID int) (entities.Product, error) {
	var product entities.Product

	err := p.db.Where("id = ? ", ID).First(&product).Error
	return product, err
}

func (p ProductRepository) Store(ctx context.Context, product entities.Product) (entities.Product, error) {

	err := p.db.Create(&product).Error
	//for _, Ip := range product.ProductImage {
	//	fmt.Println("insed ", Ip, "| product id : ", product.ID)
	//	p.db.Create(entities.ProductImages{
	//		ProductID: product.ID,
	//		Path:      Ip.Path,
	//	})
	//}

	return product, err
}

func (p ProductRepository) GetRootAttributes(c *gin.Context, productID int) ([]entities.Attribute, error) {
	var category entities.Category
	var attributes []entities.Attribute

	var product entities.Product
	pErr := p.db.Where("id = ? ", productID).First(&product).Error

	if pErr != nil {
		return attributes, pErr
	}

	err := p.db.Raw(
		` WITH RECURSIVE CategoryHierarchy AS (
            SELECT id, title, parent_id
            FROM categories
            WHERE id = ?

            UNION ALL

            SELECT c.id, c.title, c.parent_id
            FROM categories c
            INNER JOIN CategoryHierarchy ch ON c.id = ch.parent_id
        )
        SELECT *
        FROM CategoryHierarchy
        WHERE parent_id IS NULL
        LIMIT 1;`, product.CategoryID,
	).Scan(&category).Error

	if err != nil {
		fmt.Println("proudct repository _ root category not found")
		return attributes, err
	}

	aErr := p.db.Preload("AttributeValues").Where("category_id = ? ", category.ID).Find(&attributes).Error

	return attributes, aErr
}
