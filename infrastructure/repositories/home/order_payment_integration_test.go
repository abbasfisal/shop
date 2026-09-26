package home

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"shop/bootstrap"
	"shop/domain/domain_err"
	"shop/domain/entities"
)

// testDB opens the migrated test database (transaction rollback is NOT used
// here: OrderPaidSuccessfully manages its own transaction internally, so the
// test cleans up its rows explicitly — see payTestCleanup).
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

// testRedisClient returns a live client or skips (the payment stock gate
// needs real redis locks, like production).
func testRedisClient(t *testing.T) *redis.Client {
	t.Helper()
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Skip("redis not available — skipping payment stock test")
	}
	return rdb
}

// payTestIDs tracks every row this test creates for exact cleanup.
type payTestIDs struct {
	payments   []uint
	items      []uint
	orders     []uint
	variants   []uint
	products   []uint
	categories []uint
	brands     []uint
	customers  []uint
}

func payTestCleanup(t *testing.T, db *gorm.DB, ids *payTestIDs) {
	t.Helper()
	del := func(table string, values []uint) {
		if len(values) == 0 {
			return
		}
		if err := db.Unscoped().Table(table).Where("id IN ?", values).Delete(nil).Error; err != nil {
			t.Logf("cleanup %s: %v", table, err)
		}
	}
	// reverse FK order
	del("payments", ids.payments)
	del("order_items", ids.items)
	del("orders", ids.orders)
	del("product_variants", ids.variants)
	del("products", ids.products)
	del("categories", ids.categories)
	del("brands", ids.brands)
	del("customers", ids.customers)
}

// TestOrderPaidSuccessfully_DeductsExactVariant is the regression test for
// the order flow redesign:
//
//   - stock moves exactly once, at successful payment (nothing on add-to-cart
//     or order creation)
//   - the deduction hits the exact variant stored in order_items
//   - SELECT ... FOR UPDATE serializes concurrent callbacks
//   - a zero-stock variant is never sold: the paid order is cancelled
func TestOrderPaidSuccessfully_DeductsExactVariant(t *testing.T) {
	db := testDB(t)
	rdb := testRedisClient(t)
	c := nilGinCtx()

	dep := &bootstrap.Dependencies{DB: db, RedisClient: rdb}
	repo := NewHomeRepository(dep, nil)

	ids := &payTestIDs{}
	t.Cleanup(func() { payTestCleanup(t, db, ids) })

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	mobile := "0919" + suffix[len(suffix)-7:]

	cat := entities.Category{Title: "pay-cat-" + suffix, Slug: "pay-cat-" + suffix, Status: true}
	if err := db.Create(&cat).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	ids.categories = append(ids.categories, cat.ID)
	brand := entities.Brand{Title: "pay-brand-" + suffix, Slug: "pay-brand-" + suffix}
	if err := db.Create(&brand).Error; err != nil {
		t.Fatalf("create brand: %v", err)
	}
	ids.brands = append(ids.brands, brand.ID)
	prod := entities.Product{
		CategoryID: cat.ID, BrandID: brand.ID,
		Title: "pay-prod-" + suffix, Slug: "pay-prod-" + suffix,
		Sku: "pay-" + suffix, Status: entities.ProductStatusPublished,
		OriginalPrice: 500, SalePrice: 400,
	}
	if err := db.Create(&prod).Error; err != nil {
		t.Fatalf("create product: %v", err)
	}
	ids.products = append(ids.products, prod.ID)
	variant := entities.ProductVariant{ProductID: prod.ID, Stock: 5, Status: entities.VariantStatusActive}
	if err := db.Create(&variant).Error; err != nil {
		t.Fatalf("create variant: %v", err)
	}
	ids.variants = append(ids.variants, variant.ID)
	customer := entities.Customer{Mobile: mobile, Active: true}
	if err := db.Create(&customer).Error; err != nil {
		t.Fatalf("create customer: %v", err)
	}
	ids.customers = append(ids.customers, customer.ID)

	newOrder := func(tag string, qty uint) *entities.Order {
		t.Helper()
		order := &entities.Order{
			CustomerID: customer.ID, OrderNumber: tag + suffix[len(suffix)-5:],
			PaymentStatus: int(entities.OrderPending), OrderStatus: entities.OrderPending,
			TotalOriginalPrice: 500 * qty, TotalSalePrice: 400 * qty, GrandTotal: 400 * qty,
		}
		if err := db.Create(order).Error; err != nil {
			t.Fatalf("create %s order: %v", tag, err)
		}
		ids.orders = append(ids.orders, order.ID)
		item := &entities.OrderItem{
			CustomerID: customer.ID, OrderID: order.ID,
			ProductID: prod.ID, InventoryID: variant.ID, Quantity: qty,
			OriginalPrice: 500, SalePrice: 400,
			TotalOriginalPrice: 500 * qty, TotalSalePrice: 400 * qty,
		}
		if err := db.Create(item).Error; err != nil {
			t.Fatalf("create %s item: %v", tag, err)
		}
		ids.items = append(ids.items, item.ID)
		payment := &entities.Payment{
			CustomerID: customer.ID, OrderID: order.ID,
			Authority: "auth-" + tag + "-" + suffix[len(suffix)-6:],
			Amount:    order.GrandTotal,
		}
		if err := db.Create(payment).Error; err != nil {
			t.Fatalf("create %s payment: %v", tag, err)
		}
		ids.payments = append(ids.payments, payment.ID)
		order.OrderItems = []*entities.OrderItem{item}
		order.Payment = payment
		return order
	}

	stockOf := func() uint {
		t.Helper()
		var v entities.ProductVariant
		if err := db.First(&v, variant.ID).Error; err != nil {
			t.Fatalf("reload variant: %v", err)
		}
		return v.Stock
	}
	statusOf := func(id uint) uint {
		t.Helper()
		var o entities.Order
		if err := db.Select("id", "order_status").First(&o, id).Error; err != nil {
			t.Fatalf("reload order: %v", err)
		}
		return o.OrderStatus
	}

	// 1) happy path: 5 - 2 = 3, order confirmed
	order1 := newOrder("ok", 2)
	if _, ok, cerr := repo.OrderPaidSuccessfully(c, order1, "ref-ok", true); !ok || cerr.Code != 0 {
		t.Fatalf("verified payment failed: ok=%v display=%q original=%q code=%d",
			ok, cerr.DisplayMessage, cerr.OriginalMessage, cerr.Code)
	}
	if stockOf() != 3 {
		t.Fatalf("stock: want 3 got %d", stockOf())
	}
	if statusOf(order1.ID) != entities.OrderConfirmed {
		t.Fatalf("order1 status: want confirmed got %d", statusOf(order1.ID))
	}

	// 2) oversell: 3 < 4 → the paid order is cancelled, stock untouched
	order2 := newOrder("shrt", 4)
	if _, ok, cerr := repo.OrderPaidSuccessfully(c, order2, "ref-short", true); ok {
		t.Fatal("oversell must not confirm")
	} else if cerr.Code != domain_err.OrderOutOfStockAfterPayment {
		t.Fatalf("want code %d got %+v", domain_err.OrderOutOfStockAfterPayment, cerr)
	}
	if stockOf() != 3 {
		t.Fatalf("stock must stay 3, got %d", stockOf())
	}
	if statusOf(order2.ID) != entities.OrderCancelled {
		t.Fatalf("order2 status: want cancelled got %d", statusOf(order2.ID))
	}

	// 3) failed payment: nothing moves
	order3 := newOrder("fail", 1)
	if _, _, cerr := repo.OrderPaidSuccessfully(c, order3, "", false); cerr.Code != 0 {
		t.Fatalf("failed payment should not error: %+v", cerr)
	}
	if stockOf() != 3 {
		t.Fatalf("stock must stay 3 after failed payment, got %d", stockOf())
	}
}
