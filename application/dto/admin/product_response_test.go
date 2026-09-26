package responses

import (
	"testing"

	"shop/domain/entities"
)

func TestToProduct_Discount(t *testing.T) {
	p := &entities.Product{
		Title: "test",
		// OriginalPrice = admin-only purchase cost (never drives the badge)
		OriginalPrice: 600,
		SalePrice:     800,
	}
	out := ToProduct(p)
	// only sale_price is set → no discount (user rule: تخفیف فقط با discount_price)
	if out.Discount != 0 {
		t.Fatalf("expected discount 0, got %d", out.Discount)
	}
	// crossed-out base exposed to the storefront is the sale base, not the cost
	if out.OriginalPrice != 800 {
		t.Fatalf("expected crossed-out base 800, got %d", out.OriginalPrice)
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
	// a simple product with a variant-level discount, like DEMO-SIMPLE-01:
	// sale=900 + discount(effective)=700 → 22% off the sale base
	p := &entities.Product{
		Title:         "discounted",
		OriginalPrice: 700, // admin-only purchase cost
		SalePrice:     900,
		MinPrice:      700,
		MaxPrice:      700,
	}
	out := ToProduct(p)
	if out.EffectivePrice != 700 {
		t.Fatalf("expected effective price 700, got %d", out.EffectivePrice)
	}
	// (900-700)/900 = 22% — measured against the sale base, never the cost
	if out.Discount != 22 {
		t.Fatalf("expected discount 22, got %d", out.Discount)
	}
}

func TestToProduct_NoVariantsKeepsNominalPrice(t *testing.T) {
	// legacy rows (MinPrice == 0) keep the nominal sale price, no discount
	// badge unless a variant discount exists
	p := &entities.Product{
		Title:         "legacy",
		OriginalPrice: 600, // admin-only purchase cost
		SalePrice:     800,
	}
	out := ToProduct(p)
	if out.EffectivePrice != 800 {
		t.Fatalf("expected effective price 800, got %d", out.EffectivePrice)
	}
	if out.Discount != 0 {
		t.Fatalf("expected discount 0, got %d", out.Discount)
	}
}
