package site_banner

import (
	"context"

	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"

	"github.com/gin-gonic/gin"
)

type SiteBannerService struct {
	repo repositories.SiteBannerRepositoryInterface
}

func NewSiteBannerService(repo repositories.SiteBannerRepositoryInterface) *SiteBannerService {
	return &SiteBannerService{repo: repo}
}

func (s *SiteBannerService) Create(c *gin.Context, req requests.CreateSiteBannerRequest) error {
	return s.repo.Insert(c, req)
}

// Index lists every banner for the admin index page (placement is an optional
// filter passed from ?placement=).
func (s *SiteBannerService) Index(c *gin.Context, placement string) ([]*entities.SiteBanner, error) {
	return s.repo.GetAll(c, placement)
}

// Show loads one banner for the admin edit form.
func (s *SiteBannerService) Show(c *gin.Context, bannerID uint) (*entities.SiteBanner, error) {
	return s.repo.FindByID(c, bannerID)
}

// Update writes the banner columns (and its images when new files arrived).
func (s *SiteBannerService) Update(c *gin.Context, bannerID uint, req requests.CreateSiteBannerRequest) error {
	return s.repo.Update(c, bannerID, req)
}

// Delete soft-deletes a banner.
func (s *SiteBannerService) Delete(c *gin.Context, bannerID uint) error {
	return s.repo.Delete(c, bannerID)
}

// Active is the storefront feed of one placement (on + inside the date window).
// Empty slice = «در حال حاضر بنری برای نمایش نداریم» → the template hides it.
func (s *SiteBannerService) Active(ctx context.Context, placement string) ([]*entities.SiteBanner, error) {
	return s.repo.GetActive(ctx, placement)
}
