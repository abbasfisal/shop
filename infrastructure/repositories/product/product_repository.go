package product

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
	"strings"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repositories.ProductRepositoryInterface {
	return &ProductRepository{db: db}
}

func (p *ProductRepository) GetAll(ctx context.Context) ([]*entities.Product, error) {
	var products []*entities.Product
	err := p.db.WithContext(ctx).Preload("Category").Preload("Brand").Find(&products).Error
	return products, err
}

func (p *ProductRepository) FindBy(ctx context.Context, columnName string, value any) (*entities.Product, error) {
	var product entities.Product
	condition := fmt.Sprintf("%s=?", columnName)
	err := p.db.
		WithContext(ctx).
		Preload("Category").Preload("Brand").
		Preload("ProductAttributes").Preload("ProductInventories").
		Preload("ProductImages").Preload("Features").
		First(&product, condition, value).
		Error

	return &product, err
}

func (p *ProductRepository) FindByID(ctx context.Context, ID int) (*entities.Product, error) {
	var product entities.Product
	err := p.db.WithContext(ctx).Preload("ProductAttributes").First(&product, ID).Error
	return &product, err
}

func (p *ProductRepository) Store(ctx context.Context, product *entities.Product) (*entities.Product, error) {

	err := p.db.WithContext(ctx).Create(&product).Error

	if err == nil {
		_ = SyncReadModel(ctx, p.db, product.ID)
	}
	return product, err
}

func (p *ProductRepository) Update(c *gin.Context, productID int, req *requests.UpdateProductRequest) (*entities.Product, error) {

	var product entities.Product
	pErr := p.db.WithContext(c).First(&product, productID).Error
	if pErr != nil {
		fmt.Println("---- repo product find err : ", pErr)
		return nil, pErr
	}

	updateErr := p.db.WithContext(c).Model(&product).Update("category_id", req.CategoryID).
		Update("brand_id", req.BrandID).
		Update("title", strings.TrimSpace(req.Title)).
		Update("slug", strings.TrimSpace(req.Slug)).
		Update("sku", strings.TrimSpace(req.Sku)).
		Update("status", func() bool {
			if req.Status == "" {
				return false
			}
			return true
		}()).
		Update("original_price", req.OriginalPrice).
		Update("sale_price", req.SalePrice).
		Update("description", req.Description).Error

	if updateErr != nil {
		return nil, pErr
	}

	_ = SyncReadModel(c, p.db, uint(productID))

	return &product, nil
}
