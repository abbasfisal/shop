package banner

import (
	"context"

	"github.com/gin-gonic/gin"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
)

type BannerService struct {
	repo repositories.BannerRepositoryInterface
}

func NewBannerService(repo repositories.BannerRepositoryInterface) *BannerService {
	return &BannerService{repo: repo}
}

func (b *BannerService) Create(c *gin.Context, req requests.CreateBannerRequest) error {
	err := b.repo.Insert(c, req)

	if err != nil {
		return err
	}
	return nil
}

// Index lists every banner for the admin index page.
func (b *BannerService) Index(c *gin.Context) ([]*entities.Banner, error) {
	return b.repo.GetAll(c)
}

// Show loads one banner for the admin edit form.
func (b *BannerService) Show(c *gin.Context, bannerID uint) (*entities.Banner, error) {
	return b.repo.FindByID(c, bannerID)
}

// Update writes the banner columns (and its image when a new file arrived).
func (b *BannerService) Update(c *gin.Context, bannerID uint, req requests.CreateBannerRequest) error {
	return b.repo.Update(c, bannerID, req)
}

// Delete soft-deletes a banner.
func (b *BannerService) Delete(c *gin.Context, bannerID uint) error {
	return b.repo.Delete(c, bannerID)
}

// Active is the storefront promotion feed (on + inside the date window).
func (b *BannerService) Active(ctx context.Context) ([]*entities.Banner, error) {
	return b.repo.GetActive(ctx)
}
