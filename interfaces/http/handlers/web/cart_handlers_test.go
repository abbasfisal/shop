package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	bannerUseCase "shop/application/usecases/banner"
	feeUseCase "shop/application/usecases/fee"
	"shop/application/usecases/home"
	sliderUseCase "shop/application/usecases/product_slider"
	"shop/bootstrap"
	webreq "shop/interfaces/http/requests/web"
)

// fakeHomeService embeds the interface so only AddToCart needs implementing.
type fakeHomeService struct {
	home.HomeServiceInterface
	called bool
	gotID  uint
	gotInv uint
}

func (f *fakeHomeService) AddToCart(c *gin.Context, productID uint, req webreq.AddToCartRequest) {
	f.called = true
	f.gotID = productID
	f.gotInv = req.InventoryID
}

func postAddToCart(t *testing.T, body string, referer string) (*fakeHomeService, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	svc := &fakeHomeService{}
	h := NewPublicHandler(svc,
		bannerUseCase.NewBannerService(nil),
		sliderUseCase.NewProductSliderService(nil),
		feeUseCase.NewFeeRateService(nil),
		&bootstrap.Dependencies{},
	)

	r := gin.New()
	r.POST("/add-to-cart", h.AddToCart)

	req := httptest.NewRequest(http.MethodPost, "/add-to-cart", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", referer)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return svc, w
}

func TestAddToCart_ParsesNumericProductID(t *testing.T) {
	svc, w := postAddToCart(t, "product_id=123&inventory_id=7", "/product/sku/slug")

	if !svc.called {
		t.Fatal("AddToCart service was not called")
	}
	if svc.gotID != 123 {
		t.Fatalf("expected product id 123, got %d", svc.gotID)
	}
	if svc.gotInv != 7 {
		t.Fatalf("expected inventory id 7, got %d", svc.gotInv)
	}
	if w.Code != http.StatusFound {
		t.Fatalf("expected redirect, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/product/sku/slug" {
		t.Fatalf("expected redirect to referer, got %q", loc)
	}
}

func TestAddToCart_RejectsLegacyObjectID(t *testing.T) {
	// 24-hex mongo ObjectID — must no longer be accepted
	svc, w := postAddToCart(t, "product_id=507f1f77bcf86cd799439011&inventory_id=1", "/somewhere")

	if svc.called {
		t.Fatal("service must not be called for a non-numeric product id")
	}
	if w.Code != http.StatusFound {
		t.Fatalf("expected redirect, got %d", w.Code)
	}
}
