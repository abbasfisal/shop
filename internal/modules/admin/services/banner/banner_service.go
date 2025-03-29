package banner

import "shop/internal/modules/admin/repositories/banner"

type BannerService struct {
	repo banner.BannerRepositoryInterface
}

func NewBannerService(repo banner.BannerRepositoryInterface) *BannerService {
	return &BannerService{repo: repo}
}
