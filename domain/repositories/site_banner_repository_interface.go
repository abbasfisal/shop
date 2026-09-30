package repositories

import (
	"context"

	"shop/domain/entities"
	"shop/interfaces/http/requests/admin"

	"github.com/gin-gonic/gin"
)

type SiteBannerRepositoryInterface interface {
	Insert(c *gin.Context, req requests.CreateSiteBannerRequest) error
	// GetAll lists every site banner (admin index page), optionally filtered
	// by placement (?placement=header).
	GetAll(c *gin.Context, placement string) ([]*entities.SiteBanner, error)
	// FindByID loads one banner for the admin edit form.
	FindByID(c *gin.Context, bannerID uint) (*entities.SiteBanner, error)
	// Update / Delete manage a banner from the admin panel.
	Update(c *gin.Context, bannerID uint, req requests.CreateSiteBannerRequest) error
	Delete(c *gin.Context, bannerID uint) error
	// GetActive returns the banners of one placement that are on and inside
	// the date window — the storefront feed.
	GetActive(ctx context.Context, placement string) ([]*entities.SiteBanner, error)
}
