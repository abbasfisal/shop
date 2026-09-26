package fee

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

// testCtx builds a bare gin context (the repository only needs it as a
// context.Context for gorm).
func testCtx() *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	return c
}

// testDB opens the migrated test database (same convention as the other
// repository integration tests: TEST_DATABASE_URL + transaction rollback).
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
	if err := db.Raw("SELECT count(*) FROM fee_rates").Scan(&n).Error; err != nil {
		t.Fatalf("schema missing (run `go run . migrate` first): %v", err)
	}
	return db
}

func uintRef(v uint) *uint { return &v }

func TestActiveRateFor_Resolution(t *testing.T) {
	db := testDB(t)
	c := testCtx()
	tx := db.Begin()
	defer tx.Rollback()

	repo := NewFeeRateRepository(tx)
	ctx := c
	now := time.Now()
	tag := fmt.Sprintf("t%d", time.Now().UnixNano())

	mk := func(kind, title string, amount uint, status bool, start, end *time.Time) uint {
		t.Helper()
		row, err := repo.Store(ctx, &entities.FeeRate{
			Kind: kind, Title: title, Amount: amount, Status: status,
			StartsAt: start, EndsAt: end,
		})
		if err != nil {
			t.Fatalf("store %s: %v", title, err)
		}
		return row.ID
	}

	dayAgo := now.AddDate(0, 0, -1)
	twoDaysAgo := now.AddDate(0, 0, -2)
	yesterday := now.AddDate(0, 0, -1)

	// two dated contenders: the newer start date must win (seeded rows carry
	// NULL starts_at and always sort last, so this holds on any database)
	newerID := mk(entities.FeeKindShipping, tag+"-newer", 111_000, true, &dayAgo, nil)
	mk(entities.FeeKindShipping, tag+"-older", 222_000, true, &twoDaysAgo, nil)

	// expired + inactive rows with the NEWEST start date must still lose
	mk(entities.FeeKindShipping, tag+"-expired", 1, true, &now, &yesterday)
	mk(entities.FeeKindShipping, tag+"-inactive", 2, false, &now, nil)

	got, err := repo.ActiveRateFor(ctx, entities.FeeKindShipping, now)
	if err != nil {
		t.Fatalf("active shipping: %v", err)
	}
	if got.ID != newerID || got.Amount != 111_000 {
		t.Fatalf("wrong row resolved: id=%d amount=%d (want id=%d amount=111000)",
			got.ID, got.Amount, newerID)
	}

	// kind isolation: my titled rows show up under shipping only
	rows, err := repo.GetByKind(ctx, entities.FeeKindShipping)
	if err != nil {
		t.Fatalf("GetByKind: %v", err)
	}
	mine := 0
	for _, r := range rows {
		if len(r.Title) >= len(tag) && r.Title[:len(tag)] == tag {
			mine++
		}
	}
	if mine != 4 {
		t.Fatalf("GetByKind: want 4 tagged shipping rows, got %d", mine)
	}
}

func TestFeeRate_CRUD(t *testing.T) {
	db := testDB(t)
	c := testCtx()
	tx := db.Begin()
	defer tx.Rollback()

	repo := NewFeeRateRepository(tx)
	ctx := c

	stored, err := repo.Store(ctx, &entities.FeeRate{
		Kind: entities.FeeKindPackaging, Title: fmt.Sprintf("crud-%d", time.Now().UnixNano()),
		Amount: 23_000, Status: true,
	})
	if err != nil || stored.ID == 0 {
		t.Fatalf("store: %v", err)
	}

	found, err := repo.FindByID(ctx, stored.ID)
	if err != nil || found.Title != stored.Title {
		t.Fatalf("find: %+v (%v)", found, err)
	}

	update := *found
	update.Amount = 25_000
	update.Status = false
	if err := repo.Update(ctx, stored.ID, &update); err != nil {
		t.Fatalf("update: %v", err)
	}
	found, _ = repo.FindByID(ctx, stored.ID)
	if found.Amount != 25_000 || found.Status {
		t.Fatalf("update not applied: %+v", found)
	}

	if err := repo.Delete(ctx, stored.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := repo.FindByID(ctx, stored.ID); err == nil {
		t.Fatal("expected record not found after delete")
	}
}
