package slider

import (
	"github.com/gin-gonic/gin"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
	"shop/pkg/pagination"
)

type ProductSliderService struct {
	repo repositories.ProductSliderRepositoryInterface
}

func NewProductSliderService(repo repositories.ProductSliderRepositoryInterface) *ProductSliderService {
	return &ProductSliderService{repo: repo}
}

// Index lists every slider for the admin table.
func (s *ProductSliderService) Index(c *gin.Context) ([]*entities.ProductSlider, error) {
	return s.repo.GetAll(c)
}

// Show loads one slider with its ordered products (admin form).
func (s *ProductSliderService) Show(c *gin.Context, sliderID uint) (*entities.ProductSlider, error) {
	return s.repo.FindByID(c, sliderID)
}

// BySlug loads one slider for the storefront catalog page.
func (s *ProductSliderService) BySlug(c *gin.Context, slug string) (*entities.ProductSlider, error) {
	return s.repo.FindBySlug(c, slug)
}

// Store creates a slider (title/slot/products) in one transaction.
func (s *ProductSliderService) Store(c *gin.Context, req *requests.CreateProductSliderRequest) (*entities.ProductSlider, error) {
	return s.repo.Store(c, req)
}

// Update rewrites a slider and replaces its product set.
func (s *ProductSliderService) Update(c *gin.Context, sliderID uint, req *requests.CreateProductSliderRequest) error {
	return s.repo.Update(c, sliderID, req)
}

// Delete removes a slider from the admin panel and the storefront.
func (s *ProductSliderService) Delete(c *gin.Context, sliderID uint) error {
	return s.repo.Delete(c, sliderID)
}

// ActiveSliders is the homepage feed (published + inside the date window).
func (s *ProductSliderService) ActiveSliders(c *gin.Context) ([]*entities.ProductSlider, error) {
	return s.repo.ActiveSliders(c)
}

// Catalog is the «مشاهده همه» page of one slider.
func (s *ProductSliderService) Catalog(c *gin.Context, slider *entities.ProductSlider) (pagination.Pagination, error) {
	return s.repo.Catalog(c, slider)
}
