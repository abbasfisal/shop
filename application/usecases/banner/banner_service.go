package banner

import (
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
