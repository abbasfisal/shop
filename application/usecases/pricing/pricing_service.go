// Package pricing implements the PricingService pattern: after every variant
// change (stock, reservation, prices, attribute links) the product-level
// aggregate cache is refreshed so listings/filters can rely on
// products.min_price / max_price / total_stock / in_stock / attributes_json.
package pricing

import (
	"context"

	"shop/domain/repositories"
)

type PricingService struct {
	repo repositories.ProductRepositoryInterface
}

func NewPricingService(repo repositories.ProductRepositoryInterface) *PricingService {
	return &PricingService{repo: repo}
}

// RefreshProductAggregates recomputes the aggregate cache columns on the
// products row for the given product (golden rule: call after every variant change).
func (s *PricingService) RefreshProductAggregates(ctx context.Context, productID uint) error {
	return s.repo.RefreshProductAggregates(ctx, productID)
}
