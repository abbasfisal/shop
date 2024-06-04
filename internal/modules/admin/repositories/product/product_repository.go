package product

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"shop/internal/database/mysql"
	"shop/internal/entities"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() ProductRepository {
	return ProductRepository{db: mysql.Get()}
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
