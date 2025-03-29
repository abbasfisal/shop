package banner

import (
	"gorm.io/gorm"
)

type BannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) BannerRepositoryInterface {
	return &BannerRepository{db: db}
}
