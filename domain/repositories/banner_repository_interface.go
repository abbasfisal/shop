package repositories

import (
	"github.com/gin-gonic/gin"
	"shop/interfaces/http/requests/admin"
)

type BannerRepositoryInterface interface {
	Insert(c *gin.Context, req requests.CreateBannerRequest) error
}
