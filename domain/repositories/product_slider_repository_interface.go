package repositories

import (
	"context"

	"github.com/gin-gonic/gin"
	"shop/domain/entities"
	"shop/interfaces/http/requests/admin"
	"shop/pkg/pagination"
)

// ProductSliderRepositoryInterface manages the curated homepage product
// sliders (max entities.MaxSliderProducts products each).
type ProductSliderRepositoryInterface interface {
	GetAll(c *gin.Context) ([]*entities.ProductSlider, error)
	FindByID(c *gin.Context, sliderID uint) (*entities.ProductSlider, error)
	FindBySlug(c *gin.Context, slug string) (*entities.ProductSlider, error)
	Store(c *gin.Context, req *requests.CreateProductSliderRequest) (*entities.ProductSlider, error)
	Update(c *gin.Context, sliderID uint, req *requests.CreateProductSliderRequest) error
	Delete(c *gin.Context, sliderID uint) error

	// ActiveSliders is the storefront feed: published + inside the date
	// window, each with its products (capped at MaxSliderProducts).
	ActiveSliders(ctx context.Context) ([]*entities.ProductSlider, error)
	// Catalog powers the «مشاهده همه» page of one slider.
	Catalog(c *gin.Context, slider *entities.ProductSlider) (pagination.Pagination, error)
}
