package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Static assets must always be revalidated: without Cache-Control browsers
// heuristically cache JS/CSS for days, so a shipped storefront fix (e.g. in
// tsearch.js) would silently never reach users.
func TestNoCacheAssets_SetsRevalidateHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(noCacheAssets())
	r.GET("/assets/shop/js/tseaerch/tsearch.js", func(c *gin.Context) {
		c.String(http.StatusOK, "js")
	})
	r.GET("/search", func(c *gin.Context) {
		c.String(http.StatusOK, "html")
	})

	asset := httptest.NewRecorder()
	reqAsset, _ := http.NewRequest(http.MethodGet, "/assets/shop/js/tseaerch/tsearch.js", nil)
	r.ServeHTTP(asset, reqAsset)
	if got := asset.Header().Get("Cache-Control"); got != "no-cache, must-revalidate" {
		t.Fatalf("expected revalidate Cache-Control on /assets, got %q", got)
	}

	page := httptest.NewRecorder()
	reqPage, _ := http.NewRequest(http.MethodGet, "/search?q=test", nil)
	r.ServeHTTP(page, reqPage)
	if got := page.Header().Get("Cache-Control"); got != "" {
		t.Fatalf("expected no Cache-Control on dynamic pages, got %q", got)
	}
}
