package repositories

import (
	"github.com/gin-gonic/gin"
	"shop/domain/entities"
	"shop/interfaces/http/requests/admin"
)

type BannerRepositoryInterface interface {
	Insert(c *gin.Context, req requests.CreateBannerRequest) error
	// GetAll lists every banner (newest first) for the admin index page.
	GetAll(c *gin.Context) ([]*entities.Banner, error)
}
