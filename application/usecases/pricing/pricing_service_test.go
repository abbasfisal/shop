package pricing

import (
	"context"
	"errors"
	"testing"

	"shop/domain/repositories"
)

// fakeProductRepo embeds the repository interface so only the method under
// test needs to be implemented (unimplemented calls panic, which fails loudly).
type fakeProductRepo struct {
	repositories.ProductRepositoryInterface
	refreshed []uint
	err       error
}

func (f *fakeProductRepo) RefreshProductAggregates(ctx context.Context, productID uint) error {
	f.refreshed = append(f.refreshed, productID)
	return f.err
}

func TestRefreshProductAggregates_DelegatesToRepository(t *testing.T) {
	repo := &fakeProductRepo{}
	svc := NewPricingService(repo)

	if err := svc.RefreshProductAggregates(context.Background(), 42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.refreshed) != 1 || repo.refreshed[0] != 42 {
		t.Fatalf("expected refresh for product 42, got %v", repo.refreshed)
	}
}

func TestRefreshProductAggregates_PropagatesError(t *testing.T) {
	repo := &fakeProductRepo{err: errors.New("db down")}
	svc := NewPricingService(repo)

	if err := svc.RefreshProductAggregates(context.Background(), 1); err == nil {
		t.Fatal("expected error to propagate")
	}
}
