package home

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"shop/bootstrap"
	"shop/domain/entities"
	"shop/infrastructure/repositories/product"
)

// TestGetProduct_ReadModelContract verifies the storefront single-product
// payload keeps the exact keys single_product.html consumes:
// _id / product / inventories, plus recommendations.
func TestGetProduct_ReadModelContract(t *testing.T) {
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

	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	tx := db.Begin()
	defer tx.Rollback()

	cat := entities.Category{Title: "sf-cat-" + suffix, Slug: "sf-cat-" + suffix, Status: true}
	if err := tx.Create(&cat).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	brand := entities.Brand{Title: "sf-brand-" + suffix, Slug: "sf-brand-" + suffix}
	if err := tx.Create(&brand).Error; err != nil {
		t.Fatalf("create brand: %v", err)
	}

	main := entities.Product{
		CategoryID: cat.ID, BrandID: brand.ID,
		Title: "storefront-" + suffix, Slug: "storefront-" + suffix,
		Sku: "sf-" + suffix, Status: entities.ProductStatusPublished, OriginalPrice: 100, SalePrice: 90,
	}
	if err := tx.Create(&main).Error; err != nil {
		t.Fatalf("create main: %v", err)
	}
	rec := entities.Product{
		CategoryID: cat.ID, BrandID: brand.ID,
		Title: "recommended-" + suffix, Slug: "recommended-" + suffix,
		Sku: "rec-" + suffix, Status: entities.ProductStatusPublished, OriginalPrice: 50, SalePrice: 40,
	}
	if err := tx.Create(&rec).Error; err != nil {
		t.Fatalf("create recommended: %v", err)
	}

	v := entities.ProductVariant{ProductID: main.ID, Stock: 4, Status: entities.VariantStatusActive}
	if err := tx.Create(&v).Error; err != nil {
		t.Fatalf("create variant: %v", err)
	}
	if err := tx.Create(&entities.ProductRecommendation{
		ProductID: main.ID, RecommendedProductID: rec.ID,
	}).Error; err != nil {
		t.Fatalf("create recommendation: %v", err)
	}

	if err := product.SyncReadModel(ctx, tx, main.ID); err != nil {
		t.Fatalf("sync main: %v", err)
	}
	if err := product.SyncReadModel(ctx, tx, rec.ID); err != nil {
		t.Fatalf("sync rec: %v", err)
	}

	repo := NewHomeRepository(&bootstrap.Dependencies{DB: tx}, nil)
	pm, recs, err := repo.GetProduct(nilGinCtx(), main.Sku, main.Slug)
	if err != nil {
		t.Fatalf("GetProduct: %v", err)
	}

	for _, key := range []string{"_id", "product", "inventories"} {
		if _, ok := pm[key]; !ok {
			t.Fatalf("missing template key %q (have %v)", key, keys(pm))
		}
	}
	if pm["_id"] != fmt.Sprintf("%d", main.ID) {
		t.Errorf("_id: want %q got %v", fmt.Sprintf("%d", main.ID), pm["_id"])
	}

	// product payload must expose the fields single_product.html reads
	prodPayload, ok := pm["product"].(entities.P)
	if !ok {
		t.Fatalf("product payload has unexpected type %T", pm["product"])
	}
	if prodPayload.Title != main.Title || prodPayload.Sku != main.Sku {
		t.Errorf("product payload mismatch: %+v", prodPayload)
	}

	inventories, ok := pm["inventories"].(map[string]entities.Inventory)
	if !ok || len(inventories) != 1 {
		t.Fatalf("expected 1 inventory, got %T len=%d", pm["inventories"], len(inventories))
	}

	if len(recs) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(recs))
	}
	if recs[0].Product.Title != rec.Title {
		t.Errorf("recommendation title: want %q got %q", rec.Title, recs[0].Product.Title)
	}

	// _id round-trip: this is what the cart form posts back
	var rm entities.ProductReadModel
	var stored entities.Product
	if err := tx.First(&stored, main.ID).Error; err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(stored.ReadModel, &rm); err != nil {
		t.Fatalf("read_model: %v", err)
	}
	_ = rm
}

func keys(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func nilGinCtx() *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	return c
}
