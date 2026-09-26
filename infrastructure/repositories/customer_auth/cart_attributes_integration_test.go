package customer_auth

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"shop/domain/entities"
)

// TestFindCustomerBySessionID_LoadsVariantAttributes verifies /checkout/cart
// receives the labels of the picked variant («رنگ: قرمز», «سایز: M») for every
// cart line, in attribute order.
func TestFindCustomerBySessionID_LoadsVariantAttributes(t *testing.T) {
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
	ginCtx, _ := gin.CreateTestContext(nil)
	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	tx := db.Begin()
	defer tx.Rollback()

	cat := entities.Category{Title: "ca-cat-" + suffix, Slug: "ca-cat-" + suffix, Status: true}
	if err := tx.Create(&cat).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	brand := entities.Brand{Title: "ca-brand-" + suffix, Slug: "ca-brand-" + suffix}
	if err := tx.Create(&brand).Error; err != nil {
		t.Fatalf("create brand: %v", err)
	}
	prod := entities.Product{
		CategoryID: cat.ID, BrandID: brand.ID,
		Title: "ca-prod-" + suffix, Slug: "ca-prod-" + suffix, Sku: "ca-" + suffix,
		Status: entities.ProductStatusPublished, OriginalPrice: 100, SalePrice: 90,
	}
	if err := tx.Create(&prod).Error; err != nil {
		t.Fatalf("create product: %v", err)
	}
	variant := entities.ProductVariant{ProductID: prod.ID, Stock: 3, Status: entities.VariantStatusActive}
	if err := tx.Create(&variant).Error; err != nil {
		t.Fatalf("create variant: %v", err)
	}

	color := entities.Attribute{Title: "رنگ", InputType: entities.AttributeInputColor, SortOrder: 1}
	if err := tx.Create(&color).Error; err != nil {
		t.Fatalf("create attribute: %v", err)
	}
	red := entities.AttributeValue{AttributeID: color.ID, AttributeTitle: color.Title, Value: "قرمز", SortOrder: 1}
	if err := tx.Create(&red).Error; err != nil {
		t.Fatalf("create attribute value: %v", err)
	}
	size := entities.Attribute{Title: "سایز", InputType: entities.AttributeInputText, SortOrder: 2}
	if err := tx.Create(&size).Error; err != nil {
		t.Fatalf("create attribute: %v", err)
	}
	m := entities.AttributeValue{AttributeID: size.ID, AttributeTitle: size.Title, Value: "M", SortOrder: 1}
	if err := tx.Create(&m).Error; err != nil {
		t.Fatalf("create attribute value: %v", err)
	}
	for _, v := range []*entities.AttributeValue{&red, &m} {
		link := entities.VariantAttributeValue{VariantID: variant.ID, AttributeValueID: v.ID}
		if err := tx.Create(&link).Error; err != nil {
			t.Fatalf("create variant_attribute_values: %v", err)
		}
	}

	customer := entities.Customer{Mobile: "0912" + suffix[len(suffix)-7:], Active: true}
	if err := tx.Create(&customer).Error; err != nil {
		t.Fatalf("create customer: %v", err)
	}
	cart := entities.Cart{CustomerID: customer.ID, Status: 0}
	if err := tx.Create(&cart).Error; err != nil {
		t.Fatalf("create cart: %v", err)
	}
	lines := []entities.CartItem{
		{
			CustomerID: customer.ID, CartID: cart.ID, ProductID: prod.ID,
			InventoryID: variant.ID, Quantity: 1, OriginalPrice: 100, SalePrice: 90,
			ProductSku: prod.Sku, ProductTitle: prod.Title, ProductSlug: prod.Slug,
		},
		{
			CustomerID: customer.ID, CartID: cart.ID, ProductID: prod.ID,
			Quantity: 1, OriginalPrice: 100, SalePrice: 90, // stock-only line
			ProductSku: prod.Sku, ProductTitle: prod.Title, ProductSlug: prod.Slug,
		},
	}
	if err := tx.Create(&lines).Error; err != nil {
		t.Fatalf("create cart items: %v", err)
	}

	sessionID := "sess-" + suffix
	sess := entities.Session{Mobile: customer.Mobile, CustomerID: customer.ID, SessionID: sessionID, IsActive: true}
	if err := tx.Create(&sess).Error; err != nil {
		t.Fatalf("create session: %v", err)
	}

	repo := NewAuthenticateRepository(tx)
	got, err := repo.FindCustomerBySessionID(ginCtx, sessionID)
	if err != nil {
		t.Fatalf("FindCustomerBySessionID: %v", err)
	}
	if len(got.Carts) != 1 || len(got.Carts[0].CartItems) != 2 {
		t.Fatalf("expected 1 cart with 2 lines, got %d carts", len(got.Carts))
	}

	withVariant := got.Carts[0].CartItems[0]
	if withVariant.InventoryID != variant.ID {
		t.Fatalf("expected the variant line first, got inventory %d", withVariant.InventoryID)
	}
	if len(withVariant.Attributes) != 2 {
		t.Fatalf("expected 2 variant attributes, got %+v", withVariant.Attributes)
	}
	want := [][2]string{{"رنگ", "قرمز"}, {"سایز", "M"}}
	for i, w := range want {
		if withVariant.Attributes[i].Title != w[0] || withVariant.Attributes[i].Value != w[1] {
			t.Errorf("attribute %d: got %q:%q, want %q:%q",
				i, withVariant.Attributes[i].Title, withVariant.Attributes[i].Value, w[0], w[1])
		}
	}

	stockOnly := got.Carts[0].CartItems[1]
	if len(stockOnly.Attributes) != 0 {
		t.Errorf("stock-only line must carry no attributes, got %+v", stockOnly.Attributes)
	}
}
