package banner

import (
	"github.com/gin-gonic/gin"
	"shop/internal/modules/admin/repositories/banner"
	"shop/internal/modules/admin/requests"
)

type BannerService struct {
	repo banner.BannerRepositoryInterface
}

func NewBannerService(repo banner.BannerRepositoryInterface) *BannerService {
	return &BannerService{repo: repo}
}

func (b *BannerService) Create(c *gin.Context, req requests.CreateBannerRequest) error {
	err := b.repo.Insert(c, req)

	if err != nil {
		return err
	}
	return nil
}
