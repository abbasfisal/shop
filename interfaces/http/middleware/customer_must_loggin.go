package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"shop/pkg/helpers"
	"shop/pkg/sessions"
)

const (
	// SessionAuthNext is the page the visitor wanted before the login wall
	// (the product page / cart / profile …) so login does not dump them on "/".
	SessionAuthNext = "auth_next"
	// SessionPendingAddToCart holds "productID|inventoryID": an add-to-cart
	// posted while logged out, replayed right after the OTP login succeeds.
	SessionPendingAddToCart = "pending_add_to_cart"
)

// LoginIntent is what the visitor was doing when the login wall stopped them.
type LoginIntent struct {
	Next        string // where to land after login ("" → "/")
	ProductID   uint   // pending add-to-cart product, 0 when there is none
	InventoryID uint   // selected variant of that product (0 = stock-only)
}

// HasPendingAddToCart reports whether an add-to-cart should be replayed.
func (i LoginIntent) HasPendingAddToCart() bool { return i.ProductID > 0 }

// CustomerMustLogin  this middleware just check is customer LoggedIn if not redirect
func CustomerMustLogin() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("[Middleware] : CustomerMustLogin")

		customer, ok := helpers.GetAuthUser(c)
		if !ok || customer.ID <= 0 {
			rememberLoginIntent(c)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}

// TakeLoginIntent reads the remembered intent and clears it, so a later
// login can only resume it once.
func TakeLoginIntent(c *gin.Context) LoginIntent {
	intent := LoginIntent{Next: sessions.GET(c, SessionAuthNext)}
	sessions.Remove(c, SessionAuthNext)

	raw := sessions.GET(c, SessionPendingAddToCart)
	sessions.Remove(c, SessionPendingAddToCart)
	if raw != "" {
		parts := strings.SplitN(raw, "|", 2)
		if pid, err := strconv.ParseUint(parts[0], 10, 64); err == nil {
			intent.ProductID = uint(pid)
		}
		if len(parts) == 2 && parts[1] != "" {
			if inv, err := strconv.ParseUint(parts[1], 10, 64); err == nil {
				intent.InventoryID = uint(inv)
			}
		}
	}
	return intent
}

// SafeNext turns a stored destination into a redirect target that can only
// point back into this site: a local path, never a scheme/host, never a
// protocol-relative "//" and never a backslash (browsers read "/\host" as
// "//host").
func SafeNext(raw string) string {
	if raw == "" ||
		!strings.HasPrefix(raw, "/") ||
		strings.HasPrefix(raw, "//") ||
		strings.Contains(raw, `\`) {
		return "/"
	}
	return raw
}

// LocalReferer is the same-host path+query of the request Referer, or "" when
// it is missing / points at another site. Handlers use it instead of the raw
// Referer so a redirect can never leave this site.
func LocalReferer(c *gin.Context) string {
	return refererTarget(c)
}

// rememberLoginIntent stores where the visitor came from and — for the
// add-to-cart POST — what they tried to add, so PostVerifyOtp can resume it.
func rememberLoginIntent(c *gin.Context) {
	if target := localTarget(c); target != "" {
		sessions.Set(c, SessionAuthNext, target)
	}

	if c.Request.Method != http.MethodPost || c.Request.URL.Path != "/add-to-cart" {
		return
	}
	_ = c.Request.ParseForm()
	productID := c.Request.PostFormValue("product_id")
	if _, err := strconv.ParseUint(productID, 10, 64); err != nil {
		return
	}
	sessions.Set(c, SessionPendingAddToCart,
		productID+"|"+c.Request.PostFormValue("inventory_id"))
}

// localTarget returns the same-host path+query the visitor should return to:
// the current URI for GETs, the referer (the page that submitted the form)
// for POSTs. Cross-host referers are ignored so login can never be turned
// into an open redirect.
func localTarget(c *gin.Context) string {
	if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead {
		return pathWithQuery(c.Request.URL.Path, c.Request.URL.RawQuery)
	}
	return refererTarget(c)
}

// refererTarget is the request Referer reduced to a local path+query.
func refererTarget(c *gin.Context) string {
	ref := c.Request.Referer()
	if ref == "" {
		return ""
	}
	u, err := url.Parse(ref)
	if err != nil {
		return ""
	}
	if u.Host != "" && !strings.EqualFold(u.Host, c.Request.Host) {
		return ""
	}
	return pathWithQuery(u.Path, u.RawQuery)
}

func pathWithQuery(path, rawQuery string) string {
	if path == "" || !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") {
		return ""
	}
	if rawQuery == "" {
		return path
	}
	return path + "?" + rawQuery
}
