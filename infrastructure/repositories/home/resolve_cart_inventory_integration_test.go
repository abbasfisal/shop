package home

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"shop/bootstrap"
	"shop/domain/domain_err"
	"shop/domain/entities"
)

// TestResolveCartInventory guards the add-to-cart gate: a plain stock product
// takes inventory 0, a variant product refuses an empty / foreign / unsellable
// selection instead of writing a cart line no variant can explain.
func TestResolveCartInventory(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set — skipping integration test")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	var n int64
	if err := db.Raw("SELECT count(*) FROM products").Scan(&n).Error; err != nil {
		t.Fatalf("schema missing (run `go run . migrate` first): %v", err)
	}

	ginCtx := nilGinCtx()
	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	tx := db.Begin()
	defer tx.Rollback()

	cat := entities.Category{Title: "rc-cat-" + suffix, Slug: "rc-cat-" + suffix, Status: true}
	if err := tx.Create(&cat).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	brand := entities.Brand{Title: "rc-brand-" + suffix, Slug: "rc-brand-" + suffix}
	if err := tx.Create(&brand).Error; err != nil {
		t.Fatalf("create brand: %v", err)
	}

	newProduct := func(title string) entities.Product {
		p := entities.Product{
			CategoryID: cat.ID, BrandID: brand.ID,
			Title: title, Slug: title, Sku: title,
			Status: entities.ProductStatusPublished, OriginalPrice: 100, SalePrice: 90,
		}
		if err := tx.Create(&p).Error; err != nil {
			t.Fatalf("create product: %v", err)
		}
		return p
	}
	plain := newProduct("rc-plain-" + suffix)
	varied := newProduct("rc-varied-" + suffix)
	other := newProduct("rc-other-" + suffix)

	repo := NewHomeRepository(&bootstrap.Dependencies{DB: tx}, nil)

	// stock-only product: no combination to pick
	if got, err := repo.ResolveCartInventory(ginCtx, plain.ID, 0); err != nil || got != 0 {
		t.Errorf("stock-only product: got (%d, %v), want (0, nil)", got, err)
	}

	// variant product: every line needs a real, sellable variant
	if _, err := repo.ResolveCartInventory(ginCtx, varied.ID, 0); !errors.Is(err, domain_err.VariantNotSelected) {
		t.Errorf("missing selection: got %v, want VariantNotSelected", err)
	}

	sellable := entities.ProductVariant{ProductID: varied.ID, Stock: 5, Status: entities.VariantStatusActive}
	if err := tx.Create(&sellable).Error; err != nil {
		t.Fatalf("create variant: %v", err)
	}
	soldOut := entities.ProductVariant{ProductID: varied.ID, Stock: 0, Status: entities.VariantStatusActive}
	if err := tx.Create(&soldOut).Error; err != nil {
		t.Fatalf("create variant: %v", err)
	}
	inactive := entities.ProductVariant{ProductID: varied.ID, Stock: 5, Status: entities.VariantStatusInactive}
	if err := tx.Create(&inactive).Error; err != nil {
		t.Fatalf("create variant: %v", err)
	}

	if got, err := repo.ResolveCartInventory(ginCtx, varied.ID, sellable.ID); err != nil || got != sellable.ID {
		t.Errorf("sellable variant: got (%d, %v), want (%d, nil)", got, err, sellable.ID)
	}
	if _, err := repo.ResolveCartInventory(ginCtx, varied.ID, soldOut.ID); !errors.Is(err, domain_err.OutOfStock) {
		t.Errorf("sold out: got %v, want OutOfStock", err)
	}
	if _, err := repo.ResolveCartInventory(ginCtx, varied.ID, inactive.ID); !errors.Is(err, domain_err.OutOfStock) {
		t.Errorf("inactive: got %v, want OutOfStock", err)
	}
	if _, err := repo.ResolveCartInventory(ginCtx, varied.ID, 999999); !errors.Is(err, domain_err.VariantNotSelected) {
		t.Errorf("unknown variant: got %v, want VariantNotSelected", err)
	}
	if _, err := repo.ResolveCartInventory(ginCtx, other.ID, sellable.ID); !errors.Is(err, domain_err.VariantNotSelected) {
		t.Errorf("variant of another product: got %v, want VariantNotSelected", err)
	}
}
