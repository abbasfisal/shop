package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	customerResponse "shop/application/dto/web"
)

func TestCartInventoryIDsJSON_OnlyThisProductsLines(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set("auth", customerResponse.Customer{
		ID: 1,
		Cart: customerResponse.Cart{
			CartItem: customerResponse.CartItem{
				Data: []customerResponse.Item{
					{ProductID: 5, InventoryID: 7},
					{ProductID: 5, InventoryID: 9},
					{ProductID: 6, InventoryID: 3}, // another product: ignored
				},
			},
		},
	})

	got := cartInventoryIDsJSON(c, map[string]interface{}{"_id": "5"})
	if got != "[7,9]" {
		t.Fatalf("expected [7,9], got %s", got)
	}
}

func TestCartInventoryIDsJSON_GuestAndMissingProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	guest, _ := gin.CreateTestContext(httptest.NewRecorder())
	if got := cartInventoryIDsJSON(guest, map[string]interface{}{"_id": "5"}); got != "[]" {
		t.Errorf("guest must get [], got %s", got)
	}

	authed, _ := gin.CreateTestContext(httptest.NewRecorder())
	authed.Set("auth", customerResponse.Customer{
		ID:   1,
		Cart: customerResponse.Cart{CartItem: customerResponse.CartItem{Data: []customerResponse.Item{{ProductID: 5, InventoryID: 7}}}},
	})
	if got := cartInventoryIDsJSON(authed, nil); got != "[]" {
		t.Errorf("product without id must get [], got %s", got)
	}
}
