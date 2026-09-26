package banner

import (
	"context"

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

// bannerColumns is the value set shared by Insert/Update.
func bannerValues(req *requests.CreateBannerRequest) map[string]interface{} {
	startsAt, endsAt := requests.BannerDates(req.StartsAt, req.EndsAt)
	return map[string]interface{}{
		// "type" is the legacy placement column — never rewritten by the
		// promotion form (it stays whatever the old rows had).
		"title":      req.Title,
		"layout":     req.LayoutValue(),
		"link":       req.Link,
		"priority":   req.Priority,
		"sort_order": req.SortOrder,
		"status":     req.IsActive(),
		"starts_at":  startsAt,
		"ends_at":    endsAt,
	}
}

func (b BannerRepository) Insert(c *gin.Context, req requests.CreateBannerRequest) error {
	row := entities.Banner{
		Image:     req.BannerImage,
		Type:      req.Type,
		Title:     req.Title,
		Layout:    req.LayoutValue(),
		Link:      req.Link,
		Priority:  req.Priority,
		SortOrder: req.SortOrder,
		Status:    req.IsActive(),
	}
	row.StartsAt, row.EndsAt = requests.BannerDates(req.StartsAt, req.EndsAt)

	if err := b.db.WithContext(c).Create(&row).Error; err != nil {
		logging.Log.WithError(err).WithFields(logrus.Fields{"method": "Insert", "file": "banner_repository"})
		return err
	}
	return nil
}

// Update writes the editable columns of one banner.
func (b BannerRepository) Update(c *gin.Context, bannerID uint, req requests.CreateBannerRequest) error {
	values := bannerValues(&req)
	if req.BannerImage != "" {
		values["image"] = req.BannerImage
	}
	err := b.db.WithContext(c).
		Model(&entities.Banner{}).
		Where("id = ?", bannerID).
		Updates(values).Error
	if err != nil {
		logging.Log.WithError(err).WithFields(logrus.Fields{"method": "Update", "file": "banner_repository"})
	}
	return err
}

// Delete soft-deletes one banner (removes it from the storefront feed).
func (b BannerRepository) Delete(c *gin.Context, bannerID uint) error {
	err := b.db.WithContext(c).Delete(&entities.Banner{}, bannerID).Error
	if err != nil {
		logging.Log.WithError(err).WithFields(logrus.Fields{"method": "Delete", "file": "banner_repository"})
	}
	return err
}

// FindByID loads a single banner for the admin edit form.
func (b BannerRepository) FindByID(c *gin.Context, bannerID uint) (*entities.Banner, error) {
	var banner entities.Banner
	if err := b.db.WithContext(c).First(&banner, bannerID).Error; err != nil {
		return nil, err
	}
	return &banner, nil
}

// GetAll returns every banner ordered for the admin index page.
func (b BannerRepository) GetAll(c *gin.Context) ([]*entities.Banner, error) {
	var banners []*entities.Banner
	err := b.db.WithContext(c).
		Order("sort_order ASC, id DESC").
		Find(&banners).Error
	if err != nil {
		logging.Log.WithError(err).WithFields(logrus.Fields{"method": "GetAll", "file": "banner_repository"})
		return nil, err
	}
	return banners, nil
}

// GetActive returns the banners that are on and inside their [starts_at,
// ends_at] window — the storefront promotion feed.
func (b BannerRepository) GetActive(c context.Context) ([]*entities.Banner, error) {
	var banners []*entities.Banner
	err := b.db.WithContext(c).
		Where("status = TRUE").
		Where("(starts_at IS NULL OR starts_at <= CURRENT_DATE)").
		Where("(ends_at IS NULL OR ends_at >= CURRENT_DATE)").
		Order("sort_order ASC, id ASC").
		Find(&banners).Error
	if err != nil {
		logging.Log.WithError(err).WithFields(logrus.Fields{"method": "GetActive", "file": "banner_repository"})
		return nil, err
	}
	return banners, nil
}
