package fee

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/domain/entities"
	"shop/domain/repositories"
)

type FeeRateService struct {
	repo repositories.FeeRateRepositoryInterface
}

func NewFeeRateService(repo repositories.FeeRateRepositoryInterface) *FeeRateService {
	return &FeeRateService{repo: repo}
}

// FeeQuote is the checkout-time fee split for one cart total.
type FeeQuote struct {
	Items     uint   // cart items total (input)
	Shipping  uint   // payable shipping (0 when free or no tariff)
	Packaging uint   // payable packaging (0 when no tariff)
	Grand     uint   // Items + Shipping + Packaging
	Free      bool   // free-shipping threshold met
	Threshold uint   // configured threshold (0 = none)
	ShipTitle string // active tariff titles (receipt / tooltip)
	PackTitle string
}

// Index lists every tariff of one kind (admin table).
func (s *FeeRateService) Index(c context.Context, kind string) ([]*entities.FeeRate, error) {
	return s.repo.GetByKind(c, kind)
}

// Show loads one tariff (admin form).
func (s *FeeRateService) Show(c context.Context, id uint) (*entities.FeeRate, error) {
	return s.repo.FindByID(c, id)
}

// Store creates a tariff row.
func (s *FeeRateService) Store(c *gin.Context, rate *entities.FeeRate) (*entities.FeeRate, error) {
	return s.repo.Store(c, rate)
}

// Update rewrites a tariff row.
func (s *FeeRateService) Update(c *gin.Context, id uint, rate *entities.FeeRate) error {
	return s.repo.Update(c, id, rate)
}

// Delete soft-deletes a tariff row (history of other rows stays intact).
func (s *FeeRateService) Delete(c *gin.Context, id uint) error {
	return s.repo.Delete(c, id)
}

// Active returns the tariff that applies right now (nil, nil when none).
func (s *FeeRateService) Active(c *gin.Context, kind string) (*entities.FeeRate, error) {
	rate, err := s.repo.ActiveRateFor(c, kind, time.Now())
	if err != nil {
		return nil, err
	}
	return rate, nil
}

// Quote resolves the active tariffs for right now and splits itemsTotal.
// Missing/inactive tariffs degrade to 0 — checkout never blocks on fees.
func (s *FeeRateService) Quote(c *gin.Context, itemsTotal uint) FeeQuote {
	return QuoteAt(s.repo, c, itemsTotal, time.Now())
}

// QuoteAt is the testable core of Quote (injectable clock).
func QuoteAt(repo repositories.FeeRateRepositoryInterface, ctx context.Context, itemsTotal uint, now time.Time) FeeQuote {
	quote := FeeQuote{Items: itemsTotal}

	// a missing tariff is a normal state (fee stays 0); any other error
	// degrades the same way so checkout never blocks on fees
	ship, err := repo.ActiveRateFor(ctx, entities.FeeKindShipping, now)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return quote
		}
		ship = nil
	}
	pack, err := repo.ActiveRateFor(ctx, entities.FeeKindPackaging, now)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return quote
		}
		pack = nil
	}

	quote.Shipping, quote.Packaging, quote.Grand, quote.Free =
		entities.QuoteOrderFees(itemsTotal, ship, pack)
	if ship != nil {
		quote.ShipTitle = ship.Title
		if ship.FreeThreshold != nil {
			quote.Threshold = *ship.FreeThreshold
		}
	}
	if pack != nil {
		quote.PackTitle = pack.Title
	}
	return quote
}
