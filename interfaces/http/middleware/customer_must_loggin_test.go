package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gsessions "github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// newIntentRouter wires the login wall and a fake post-login endpoint that
// consumes the remembered intent, mirroring PostVerifyOtp.
func newIntentRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gsessions.Sessions("session", cookie.NewStore([]byte("secret"))))

	r.POST("/add-to-cart", CustomerMustLogin(), func(c *gin.Context) {
		c.String(http.StatusOK, "added")
	})
	r.GET("/checkout/cart", CustomerMustLogin(), func(c *gin.Context) {
		c.String(http.StatusOK, "cart")
	})
	r.GET("/resume", func(c *gin.Context) {
		intent := TakeLoginIntent(c)
		c.JSON(http.StatusOK, gin.H{
			"next":        SafeNext(intent.Next),
			"product_id":  intent.ProductID,
			"inventory":   intent.InventoryID,
			"has_pending": intent.HasPendingAddToCart(),
		})
	})
	return r
}

// sessionCookie mirrors what a browser keeps: one cookie per name, the last
// Set-Cookie winning (the session is saved once per write).
func sessionCookie(w *httptest.ResponseRecorder) string {
	byName := map[string]string{}
	order := []string{}
	for _, c := range w.Result().Cookies() {
		if _, seen := byName[c.Name]; !seen {
			order = append(order, c.Name)
		}
		byName[c.Name] = c.Value
	}
	parts := []string{}
	for _, name := range order {
		parts = append(parts, name+"="+byName[name])
	}
	return strings.Join(parts, "; ")
}

func TestCustomerMustLogin_RemembersAddToCartIntent(t *testing.T) {
	r := newIntentRouter()

	req := httptest.NewRequest(http.MethodPost, "/add-to-cart",
		strings.NewReader("product_id=123&inventory_id=7"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://shop.test/product/abc/blue-shirt")
	req.Host = "shop.test"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("expected redirect to login, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/login" {
		t.Fatalf("expected redirect to /login, got %q", loc)
	}

	cookie := sessionCookie(w)

	// after the OTP login the visitor resumes the product page and the
	// add-to-cart they posted
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/resume", nil)
	req2.Header.Set("Cookie", cookie)
	r.ServeHTTP(w2, req2)

	var got struct {
		Next       string `json:"next"`
		ProductID  uint   `json:"product_id"`
		Inventory  uint   `json:"inventory"`
		HasPending bool   `json:"has_pending"`
	}
	if err := json.Unmarshal(w2.Body.Bytes(), &got); err != nil {
		t.Fatalf("bad response %q: %v", w2.Body.String(), err)
	}
	if got.Next != "/product/abc/blue-shirt" {
		t.Errorf("expected the product page as next, got %q", got.Next)
	}
	if !got.HasPending || got.ProductID != 123 || got.Inventory != 7 {
		t.Errorf("expected pending add-to-cart 123/7, got %+v", got)
	}
}

func TestCustomerMustLogin_IntentIsConsumedOnce(t *testing.T) {
	r := newIntentRouter()

	req := httptest.NewRequest(http.MethodGet, "/checkout/cart", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusFound {
		t.Fatalf("expected redirect to login, got %d", w.Code)
	}
	cookie := sessionCookie(w)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/resume", nil)
	req2.Header.Set("Cookie", cookie)
	r.ServeHTTP(w2, req2)

	if !strings.Contains(w2.Body.String(), `"/checkout/cart"`) {
		t.Errorf("expected /checkout/cart as next, got %s", w2.Body.String())
	}

	// a second login must not be sent back to the old destination
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/resume", nil)
	req3.Header.Set("Cookie", sessionCookie(w2))
	r.ServeHTTP(w3, req3)

	var got struct {
		Next       string `json:"next"`
		HasPending bool   `json:"has_pending"`
	}
	_ = json.Unmarshal(w3.Body.Bytes(), &got)
	if got.HasPending {
		t.Errorf("pending add-to-cart must be consumed once, got %s", w3.Body.String())
	}
}

func TestCustomerMustLogin_IgnoresForeignReferer(t *testing.T) {
	r := newIntentRouter()

	req := httptest.NewRequest(http.MethodPost, "/add-to-cart",
		strings.NewReader("product_id=5&inventory_id=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://evil.test/phish")
	req.Host = "shop.test"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/resume", nil)
	req2.Header.Set("Cookie", sessionCookie(w))
	r.ServeHTTP(w2, req2)

	if !strings.Contains(w2.Body.String(), `"next":"/"`) {
		t.Errorf("cross-host referer must not be honoured, got %s", w2.Body.String())
	}
	// the cart add itself is still remembered (it is this site's own form)
	if !strings.Contains(w2.Body.String(), `"has_pending":true`) {
		t.Errorf("expected the add-to-cart to stay pending, got %s", w2.Body.String())
	}
}

func TestSafeNext(t *testing.T) {
	cases := map[string]string{
		"":                      "/",
		"/checkout/cart":        "/checkout/cart",
		"/product/a/b?x=1":      "/product/a/b?x=1",
		"//evil.test/x":         "/",
		"/\\evil.com/x":         "/",
		"http://evil.test/x":    "/",
		"javascript:alert(1)":   "/",
		"https://shop.test/x/y": "/",
	}
	for in, want := range cases {
		if got := SafeNext(in); got != want {
			t.Errorf("SafeNext(%q) = %q, want %q", in, got, want)
		}
	}
}
