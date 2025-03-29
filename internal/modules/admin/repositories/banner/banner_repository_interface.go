package banner

import (
	"github.com/gin-gonic/gin"
	"shop/internal/modules/admin/requests"
)

type BannerRepositoryInterface interface {
	Insert(c *gin.Context, req requests.CreateBannerRequest) error
}
