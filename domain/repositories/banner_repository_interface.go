package repositories

import (
	"context"

	"shop/domain/entities"
	"shop/interfaces/http/requests/admin"

	"github.com/gin-gonic/gin"
)

type BannerRepositoryInterface interface {
	Insert(c *gin.Context, req requests.CreateBannerRequest) error
	// GetAll lists every banner (admin index page).
	GetAll(c *gin.Context) ([]*entities.Banner, error)
	// FindByID loads one banner for the admin edit form.
	FindByID(c *gin.Context, bannerID uint) (*entities.Banner, error)
	// Update / Delete manage a banner from the admin panel.
	Update(c *gin.Context, bannerID uint, req requests.CreateBannerRequest) error
	Delete(c *gin.Context, bannerID uint) error
	// GetActive returns the promotion feed (on + inside the date window).
	GetActive(c context.Context) ([]*entities.Banner, error)
}
