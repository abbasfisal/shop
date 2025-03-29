package banner

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
	"shop/internal/pkg/logging"
)

type BannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) BannerRepositoryInterface {
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
