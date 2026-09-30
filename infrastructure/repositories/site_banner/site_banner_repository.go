package site_banner

import (
	"context"
	"strings"

	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
	"shop/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SiteBannerRepository struct {
	db *gorm.DB
}

func NewSiteBannerRepository(db *gorm.DB) repositories.SiteBannerRepositoryInterface {
	return &SiteBannerRepository{db: db}
}

// siteBannerValues is the value set shared by Insert/Update.
func siteBannerValues(req *requests.CreateSiteBannerRequest) map[string]interface{} {
	startsAt, endsAt := requests.SiteBannerDates(req.StartsAt, req.EndsAt)
	return map[string]interface{}{
		"placement":    req.PlacementValue(),
		"title":        req.Title,
		"link":         req.Link,
		"mobile_image": req.MobileImage,
		"sort_order":   req.SortOrder,
		"status":       req.IsActive(),
		"starts_at":    startsAt,
		"ends_at":      endsAt,
	}
}

func (r SiteBannerRepository) Insert(c *gin.Context, req requests.CreateSiteBannerRequest) error {
	row := entities.SiteBanner{
		Placement:   req.PlacementValue(),
		Title:       req.Title,
		Link:        req.Link,
		Image:       req.SiteBannerImage,
		MobileImage: req.MobileImage,
		Status:      req.IsActive(),
		SortOrder:   req.SortOrder,
	}
	row.StartsAt, row.EndsAt = requests.SiteBannerDates(req.StartsAt, req.EndsAt)

	if err := r.db.WithContext(c).Create(&row).Error; err != nil {
		logging.Log.WithError(err).WithFields(logrus.Fields{"method": "Insert", "file": "site_banner_repository"})
		return err
	}
	return nil
}

// Update writes the editable columns of one site banner. The desktop image is
// only rewritten when a new file arrived; mobile_image is written as-is so the
// admin can clear it (empty → fall back to the desktop image on the storefront).
func (r SiteBannerRepository) Update(c *gin.Context, bannerID uint, req requests.CreateSiteBannerRequest) error {
	values := siteBannerValues(&req)
	if req.SiteBannerImage != "" {
		values["image"] = req.SiteBannerImage
	}
	err := r.db.WithContext(c).
		Model(&entities.SiteBanner{}).
		Where("id = ?", bannerID).
		Updates(values).Error
	if err != nil {
		logging.Log.WithError(err).WithFields(logrus.Fields{"method": "Update", "file": "site_banner_repository"})
	}
	return err
}

// Delete soft-deletes one banner (removes it from the storefront feed).
func (r SiteBannerRepository) Delete(c *gin.Context, bannerID uint) error {
	err := r.db.WithContext(c).Delete(&entities.SiteBanner{}, bannerID).Error
	if err != nil {
		logging.Log.WithError(err).WithFields(logrus.Fields{"method": "Delete", "file": "site_banner_repository"})
	}
	return err
}

// FindByID loads a single banner for the admin edit form.
func (r SiteBannerRepository) FindByID(c *gin.Context, bannerID uint) (*entities.SiteBanner, error) {
	var banner entities.SiteBanner
	if err := r.db.WithContext(c).First(&banner, bannerID).Error; err != nil {
		return nil, err
	}
	return &banner, nil
}

// GetAll returns every banner ordered for the admin index page. placement is
// an optional filter (?placement=header); empty means every placement.
func (r SiteBannerRepository) GetAll(c *gin.Context, placement string) ([]*entities.SiteBanner, error) {
	var banners []*entities.SiteBanner
	q := r.db.WithContext(c)
	if p := strings.TrimSpace(placement); p != "" {
		q = q.Where("placement = ?", p)
	}
	err := q.Order("sort_order ASC, id DESC").Find(&banners).Error
	if err != nil {
		logging.Log.WithError(err).WithFields(logrus.Fields{"method": "GetAll", "file": "site_banner_repository"})
		return nil, err
	}
	return banners, nil
}

// GetActive returns the banners of one placement that are on and inside their
// [starts_at, ends_at] window — the storefront feed.
func (r SiteBannerRepository) GetActive(ctx context.Context, placement string) ([]*entities.SiteBanner, error) {
	var banners []*entities.SiteBanner
	err := r.db.WithContext(ctx).
		Where("placement = ?", placement).
		Where("status = TRUE").
		Where("(starts_at IS NULL OR starts_at <= CURRENT_DATE)").
		Where("(ends_at IS NULL OR ends_at >= CURRENT_DATE)").
		Order("sort_order ASC, id ASC").
		Find(&banners).Error
	if err != nil {
		logging.Log.WithError(err).WithFields(logrus.Fields{"method": "GetActive", "file": "site_banner_repository"})
		return nil, err
	}
	return banners, nil
}
