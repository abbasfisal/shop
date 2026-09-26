package banner

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
	"shop/pkg/logging"
)

type BannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) repositories.BannerRepositoryInterface {
	return &BannerRepository{db: db}
}

func (b BannerRepository) Insert(c *gin.Context, req requests.CreateBannerRequest) error {
	err := b.db.WithContext(c).Create(&entities.Banner{
		Type:     req.Type,
		Link:     req.Link,
		Priority: req.Priority,
		Status: func() bool {
			if req.Status == "on" {
				return true
			}
			return false
		}(),
		Image: req.BannerImage,
	}).Error

	if err != nil {
		logging.Log.WithError(err).WithFields(logrus.Fields{"method": "Insert", "file": "banner_repository"})
		return err
	}
	return nil
}

// GetAll returns every banner ordered by priority then newest first
// (admin banner index).
func (b BannerRepository) GetAll(c *gin.Context) ([]*entities.Banner, error) {
	var banners []*entities.Banner
	err := b.db.WithContext(c).
		Order("priority ASC, id DESC").
		Find(&banners).Error
	if err != nil {
		logging.Log.WithError(err).WithFields(logrus.Fields{"method": "GetAll", "file": "banner_repository"})
		return nil, err
	}
	return banners, nil
}
