package repositories

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/domain/entities"
)

// FeeRateRepositoryInterface manages the time-windowed order fee tariffs
// (shipping + packaging, one table behind two admin menus).
type FeeRateRepositoryInterface interface {
	GetByKind(c context.Context, kind string) ([]*entities.FeeRate, error)
	FindByID(c context.Context, id uint) (*entities.FeeRate, error)
	Store(c *gin.Context, rate *entities.FeeRate) (*entities.FeeRate, error)
	Update(c *gin.Context, id uint, rate *entities.FeeRate) error
	Delete(c *gin.Context, id uint) error
	// ActiveRateFor resolves the single tariff that applies right now.
	ActiveRateFor(ctx context.Context, kind string, now time.Time) (*entities.FeeRate, error)
}

// ActiveRateForDB is the shared resolution query usable with any gorm handle
// (plain DB or an open transaction, e.g. order creation).
func ActiveRateForDB(ctx context.Context, db *gorm.DB, kind string, now time.Time) (*entities.FeeRate, error) {
	var rate entities.FeeRate
	err := db.WithContext(ctx).
		Where("kind = ? AND status = TRUE", kind).
		Where("(starts_at IS NULL OR starts_at <= ?)", now.Format("2006-01-02")).
		Where("(ends_at IS NULL OR ends_at >= ?)", now.Format("2006-01-02")).
		Order("starts_at DESC NULLS LAST, id DESC").
		First(&rate).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}
