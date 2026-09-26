package fee

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/domain/entities"
	"shop/domain/repositories"
)

type FeeRateRepository struct {
	db *gorm.DB
}

func NewFeeRateRepository(db *gorm.DB) repositories.FeeRateRepositoryInterface {
	return &FeeRateRepository{db: db}
}

func (r *FeeRateRepository) GetByKind(c context.Context, kind string) ([]*entities.FeeRate, error) {
	var rates []*entities.FeeRate
	err := r.db.WithContext(c).
		Where("kind = ?", kind).
		Order("id DESC").
		Find(&rates).Error
	return rates, err
}

func (r *FeeRateRepository) FindByID(c context.Context, id uint) (*entities.FeeRate, error) {
	var rate entities.FeeRate
	if err := r.db.WithContext(c).First(&rate, id).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *FeeRateRepository) Store(c *gin.Context, rate *entities.FeeRate) (*entities.FeeRate, error) {
	if err := r.db.WithContext(c).Create(rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (r *FeeRateRepository) Update(c *gin.Context, id uint, rate *entities.FeeRate) error {
	return r.db.WithContext(c).
		Model(&entities.FeeRate{}).
		Where("id = ? AND kind = ?", id, rate.Kind).
		Updates(map[string]interface{}{
			"title":          rate.Title,
			"amount":         rate.Amount,
			"free_threshold": rate.FreeThreshold,
			"starts_at":      rate.StartsAt,
			"ends_at":        rate.EndsAt,
			"status":         rate.Status,
		}).Error
}

func (r *FeeRateRepository) Delete(c *gin.Context, id uint) error {
	return r.db.WithContext(c).Delete(&entities.FeeRate{}, id).Error
}

func (r *FeeRateRepository) ActiveRateFor(ctx context.Context, kind string, now time.Time) (*entities.FeeRate, error) {
	return repositories.ActiveRateForDB(ctx, r.db, kind, now)
}
