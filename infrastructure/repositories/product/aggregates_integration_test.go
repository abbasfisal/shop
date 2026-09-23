package product

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"shop/domain/entities"
)

// testDB opens a connection to a migrated PostgreSQL database.
// Run migrations first:  TEST_DATABASE_URL=... go run . migrate
func testDB(t *testing.T) *gorm.DB {
	t.Helper()
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
	return db
}

// TestRefreshProductAggregates computes min/max effective price, stock
// totals, product_type and attributes_json from active variants.
func TestRefreshProductAggregates(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()
	tx := db.Begin()
	defer tx.Rollback()

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())

	cat := entities.Category{Title: "cat-" + suffix, Slug: "cat-" + suffix, Status: true}
	if err := tx.Create(&cat).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	brand := entities.Brand{Title: "brand-" + suffix, Slug: "brand-" + suffix}
	if err := tx.Create(&brand).Error; err != nil {
		t.Fatalf("create brand: %v", err)
	}

	attr := entities.Attribute{Title: "size-" + suffix, Code: "it_" + suffix}
	if err := tx.Create(&attr).Error; err != nil {
		t.Fatalf("create attribute: %v", err)
	}
	av := entities.AttributeValue{AttributeID: attr.ID, AttributeTitle: attr.Title, Value: "L"}
	if err := tx.Create(&av).Error; err != nil {
		t.Fatalf("create attribute value: %v", err)
	}

	prod := entities.Product{
		CategoryID:    cat.ID,
		BrandID:       brand.ID,
		Title:         "agg-test-" + suffix,
		Slug:          "agg-test-" + suffix,
		Sku:           "agg-" + suffix,
		Status:        true,
		OriginalPrice: 120,
		SalePrice:     80,
	}
	if err := tx.Create(&prod).Error; err != nil {
		t.Fatalf("create product: %v", err)
	}

	// v1: explicit sale 100 → eff 100
	v1 := entities.ProductVariant{ProductID: prod.ID, Stock: 10, SalePrice: uintPtr(100), Status: entities.VariantStatusActive}
	// v2: sale 200 + discount 50 (< sale) → eff 50 (golden rule #3)
	v2 := entities.ProductVariant{ProductID: prod.ID, Stock: 20, SalePrice: uintPtr(200), DiscountPrice: uintPtr(50), ReservedStock: 5, Status: entities.VariantStatusActive}
	// v3: NULL prices → inherit product sale 80 → eff 80
	v3 := entities.ProductVariant{ProductID: prod.ID, Stock: 5, Status: entities.VariantStatusActive}
	// v4: discount 150 >= sale 100 → not a discount → eff 100
	v4 := entities.ProductVariant{ProductID: prod.ID, Stock: 7, SalePrice: uintPtr(100), DiscountPrice: uintPtr(150), Status: entities.VariantStatusActive}
	// v5: inactive — must be ignored
	v5 := entities.ProductVariant{ProductID: prod.ID, Stock: 999, Status: entities.VariantStatusInactive}

	for i, v := range []*entities.ProductVariant{&v1, &v2, &v3, &v4, &v5} {
		if err := tx.Create(v).Error; err != nil {
			t.Fatalf("create variant %d: %v", i, err)
		}
	}
	// link attribute value to v1 only → variable + attributes_json
	vav := entities.VariantAttributeValue{ProductID: prod.ID, VariantID: v1.ID, AttributeValueID: av.ID}
	if err := tx.Create(&vav).Error; err != nil {
		t.Fatalf("create vav: %v", err)
	}
	// hook must have filled ProductID when only VariantID was provided
	var vavRaw entities.VariantAttributeValue
	if err := tx.First(&vavRaw, vav.ID).Error; err != nil {
		t.Fatalf("reload vav: %v", err)
	}
	if vavRaw.ProductID != prod.ID {
		t.Fatalf("BeforeCreate hook did not fill ProductID: got %d want %d", vavRaw.ProductID, prod.ID)
	}

	if err := RefreshProductAggregates(ctx, tx, prod.ID); err != nil {
		t.Fatalf("RefreshProductAggregates: %v", err)
	}

	var got entities.Product
	if err := tx.First(&got, prod.ID).Error; err != nil {
		t.Fatalf("reload product: %v", err)
	}

	// stocks of active variants: 10+20+5+7 = 42 (inactive999 excluded)
	if got.TotalStock != 42 {
		t.Errorf("total_stock: want 42 got %d", got.TotalStock)
	}
	// reserved: 5 (only v2)
	if got.TotalReserved != 5 {
		t.Errorf("total_reserved: want 5 got %d", got.TotalReserved)
	}
	if got.AvailableStock != 37 {
		t.Errorf("available_stock: want 37 got %d", got.AvailableStock)
	}
	if !got.InStock {
		t.Error("in_stock should be true")
	}
	if got.VariantsCount != 4 {
		t.Errorf("variants_count: want 4 got %d", got.VariantsCount)
	}
	// effective prices: v1=100, v2=50, v3=80, v4=100 → min 50, max 100
	if got.MinPrice != 50 {
		t.Errorf("min_price: want 50 got %d", got.MinPrice)
	}
	if got.MaxPrice != 100 {
		t.Errorf("max_price: want 100 got %d", got.MaxPrice)
	}
	if got.ProductType != "variable" {
		t.Errorf("product_type: want variable got %q", got.ProductType)
	}

	var attrsJSON map[string][]uint
	raw := string(got.AttributesJSON)
	if raw == "" || raw == "{}" {
		t.Fatalf("attributes_json empty: %q", raw)
	}
	if err := json.Unmarshal(got.AttributesJSON, &attrsJSON); err != nil {
		t.Fatalf("attributes_json invalid: %v (%s)", err, raw)
	}
	if ids, ok := attrsJSON[attr.Code]; !ok || len(ids) != 1 || ids[0] != av.ID {
		t.Errorf("attributes_json[%s]: want [%d], got %v (json=%s)", attr.Code, av.ID, ids, raw)
	}
}

// TestSyncReadModelAndStorefrontGetProduct verifies the JSONB read model the
// storefront template consumes (_id / product / inventories) and recommendations.
func TestSyncReadModelAndStorefrontGetProduct(t *testing.T) {
	db := testDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	cat := entities.Category{Title: "rm-cat-" + suffix, Slug: "rm-cat-" + suffix, Status: true}
	if err := tx.Create(&cat).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	brand := entities.Brand{Title: "rm-brand-" + suffix, Slug: "rm-brand-" + suffix}
	if err := tx.Create(&brand).Error; err != nil {
		t.Fatalf("create brand: %v", err)
	}
	prod := entities.Product{
		CategoryID:    cat.ID,
		BrandID:       brand.ID,
		Title:         "read-model-" + suffix,
		Slug:          "read-model-" + suffix,
		Sku:           "rm-" + suffix,
		Status:        true,
		OriginalPrice: 500,
		SalePrice:     400,
	}
	if err := tx.Create(&prod).Error; err != nil {
		t.Fatalf("create product: %v", err)
	}
	v := entities.ProductVariant{ProductID: prod.ID, Stock: 3, Status: entities.VariantStatusActive}
	if err := tx.Create(&v).Error; err != nil {
		t.Fatalf("create variant: %v", err)
	}

	ctx := context.Background()
	if err := SyncReadModel(ctx, tx, prod.ID); err != nil {
		t.Fatalf("SyncReadModel: %v", err)
	}

	var stored entities.Product
	if err := tx.First(&stored, prod.ID).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	var rm entities.ProductReadModel
	if err := json.Unmarshal(stored.ReadModel, &rm); err != nil {
		t.Fatalf("read_model invalid: %v", err)
	}
	if rm.Product.ID != int64(prod.ID) || rm.Product.Title != prod.Title {
		t.Errorf("read model product mismatch: %+v", rm.Product)
	}
	if len(rm.Inventories) != 1 {
		t.Fatalf("expected 1 inventory in read model, got %d", len(rm.Inventories))
	}
	inv, ok := rm.Inventories[fmt.Sprintf("%d", v.ID)]
	if !ok || inv.Quantity != 3 {
		t.Errorf("inventory quantity mismatch: %+v", inv)
	}
}

// TestAttributeAfterCreateGeneratesCode verifies the admin create flow
// (title only) still ends up with a unique attributes.code.
func TestAttributeAfterCreateGeneratesCode(t *testing.T) {
	db := testDB(t)
	suffix := fmt.Sprintf("%d", time.Now().UnixNano())

	attr := entities.Attribute{Title: "hook-" + suffix} // Code left empty on purpose
	if err := db.Create(&attr).Error; err != nil {
		t.Fatalf("create attribute: %v", err)
	}
	defer db.Unscoped().Where("id = ?", attr.ID).Delete(&entities.Attribute{})

	if attr.Code == "" {
		t.Fatal("AfterCreate hook did not set code")
	}
	if want := fmt.Sprintf("attr_%d", attr.ID); attr.Code != want {
		t.Fatalf("code: want %q got %q", want, attr.Code)
	}
}

func uintPtr(v uint) *uint { return &v }
