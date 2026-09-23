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
