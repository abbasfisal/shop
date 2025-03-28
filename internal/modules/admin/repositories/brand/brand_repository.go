package brand

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
	"strings"
)

type BrandRepository struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) BrandRepositoryInterface {
	return &BrandRepository{db: db}
}

func (br *BrandRepository) FindBy(ctx context.Context, columnName string, value any) (*entities.Brand, error) {
	var brand entities.Brand
	condition := fmt.Sprintf("%s = ?", columnName)
	value = strings.TrimSpace(value.(string))

	err := br.db.WithContext(ctx).First(&brand, condition, value).Error
	return &brand, err
}

func (br *BrandRepository) Store(ctx context.Context, brand *entities.Brand) (*entities.Brand, error) {
	brand.Title = strings.TrimSpace(brand.Title)
	brand.Slug = strings.TrimSpace(brand.Slug)
	brand.Image = strings.TrimSpace(brand.Image)
	err := br.db.WithContext(ctx).Create(&brand).Error
	return brand, err
}

func (br *BrandRepository) GetAll(ctx context.Context) ([]*entities.Brand, error) {
	var brands []*entities.Brand
	err := br.db.WithContext(ctx).Find(&brands).Error
	return brands, err
}

func (br *BrandRepository) SelectBy(ctx context.Context, brandID int) (*entities.Brand, error) {
	var brand entities.Brand
	err := br.db.WithContext(ctx).First(&brand, "id = ?", brandID).Error
	return &brand, err
}

func (br *BrandRepository) Update(c *gin.Context, brandID int, req *requests.UpdateBrandRequest) (*entities.Brand, error) {
	var brand entities.Brand
	err := br.db.WithContext(c).First(&brand, brandID).Error
	if err != nil {
		return nil, err
	}

	err = br.db.WithContext(c).Model(&brand).Updates(entities.Brand{
		Title: strings.TrimSpace(req.Title),
		Slug:  strings.TrimSpace(req.Slug),
		Image: strings.TrimSpace(req.Image),
	}).Error

	return &brand, err
}
