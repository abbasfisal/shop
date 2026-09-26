package responses

import (
	"testing"

	"shop/domain/entities"
)

func TestToProduct_Discount(t *testing.T) {
	p := &entities.Product{
		Title:         "test",
		OriginalPrice: 1000,
		SalePrice:     800,
	}
	out := ToProduct(p)
	if out.Discount != 20 {
		t.Fatalf("expected discount 20, got %d", out.Discount)
	}
	if out.Title != "test" {
		t.Fatalf("unexpected title %q", out.Title)
	}
}

func TestToProductInventories_MapsStockToQuantity(t *testing.T) {
	variants := []*entities.ProductVariant{
		{Stock: 15},
		{Stock: 30},
	}
	variants[0].ID = 11
	variants[1].ID = 22
	out := ToProductInventories(variants)
	if out == nil || len(out.Data) != 2 {
		t.Fatalf("expected 2 inventories, got %+v", out)
	}
	if out.Data[0].Quantity != 15 || out.Data[1].Quantity != 30 {
		t.Fatalf("stock not mapped to quantity: %+v", out.Data)
	}
	if out.Data[0].ID != 11 || out.Data[1].ID != 22 {
		t.Fatalf("variant id must be preserved in dto: %+v", out.Data)
	}
}

// TestToProduct_EffectivePriceAndDiscount verifies the storefront contract:
// the customer-facing price is the aggregate minimum (not the nominal sale
// price) and the discount percent is measured against it.
func TestToProduct_EffectivePriceAndDiscount(t *testing.T) {
	// a simple product with a variant-level discount, like DEMO-SIMPLE-01
	p := &entities.Product{
		Title:         "discounted",
		OriginalPrice: 1000,
		SalePrice:     900,
		MinPrice:      700,
		MaxPrice:      700,
	}
	out := ToProduct(p)
	if out.EffectivePrice != 700 {
		t.Fatalf("expected effective price 700, got %d", out.EffectivePrice)
	}
	// (1000-700)/1000 = 30% — not the 10% a nominal computation would give
	if out.Discount != 30 {
		t.Fatalf("expected discount 30, got %d", out.Discount)
	}
}

func TestToProduct_NoVariantsKeepsNominalPrice(t *testing.T) {
	// legacy rows (MinPrice == 0) keep the old nominal behaviour
	p := &entities.Product{
		Title:         "legacy",
		OriginalPrice: 1000,
		SalePrice:     800,
	}
	out := ToProduct(p)
	if out.EffectivePrice != 800 {
		t.Fatalf("expected effective price 800, got %d", out.EffectivePrice)
	}
	if out.Discount != 20 {
		t.Fatalf("expected discount 20, got %d", out.Discount)
	}
}
