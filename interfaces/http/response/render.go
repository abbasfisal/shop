package response

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"shop/application/dto/admin"
	"shop/pkg/converters"
	"shop/pkg/helpers"
	"shop/pkg/sessions"
)

func Render(c *gin.Context, code int, name string, data gin.H) {
	data = WithGlobalData(c, data)

	format := c.DefaultQuery("format", "html")
	if format == "json" {
		c.JSON(code, data)
		return
	}
	c.HTML(code, name, data)
}

func WithGlobalData(c *gin.Context, data gin.H) gin.H {
	data["APP_NAME"] = viper.Get("APP.Name")
	data["ERRORS"] = converters.StringToMap(sessions.Flash(c, "errors"))
	data["OLDS"] = converters.StringToUrlValues(sessions.Flash(c, "olds"))
	data["MESSAGE"] = sessions.Flash(c, "message")
	// OLD is the flat (first value per key) view of OLDS — used by <select> /
	// <radio> state. It must exist on EVERY page: index on a missing key
	// aborts template execution halfway through the response.
	data["OLD"] = oldFields(c)

	// sidebar highlight: which admin menu (and group) owns the current page.
	// Computed centrally so every present and future page is covered without
	// touching individual handlers.
	activeMenu, activeParent := ActiveAdminMenu(c.Request.URL.Path, c.Query("today"))
	data["ACTIVE_MENU"] = activeMenu
	data["ACTIVE_PARENT"] = activeParent

	user := helpers.Auth(c)
	if user.ID != 0 {
		data["AUTH"] = responses.ToUserResponse(user)
	}

	return data
}

func Error500(c *gin.Context) {
	c.Redirect(http.StatusFound, "/500")
}

/**
*-----------------------------
|		customer render 🛍 مشتری
*-----------------------------
*/

func CustomerRender(c *gin.Context, code int, name string, data gin.H) {
	data = customerWithGlobalData(c, data)

	format := c.DefaultQuery("format", "html")
	if format == "json" {
		c.JSON(code, data)
		return
	}
	c.HTML(code, name, data)
}

func customerWithGlobalData(c *gin.Context, data gin.H) gin.H {
	data["APP_NAME"] = viper.Get("APP.Name")
	data["ERRORS"] = converters.StringToMap(sessions.Flash(c, "errors"))
	data["OLDS"] = converters.StringToUrlValues(sessions.Flash(c, "olds"))
	data["MESSAGE"] = sessions.Flash(c, "message")
	// OLD is the flat (first value per key) view of OLDS — used by <select> /
	// <radio> state. It must exist on EVERY page: index on a missing key
	// aborts template execution halfway through the response.
	data["OLD"] = oldFields(c)

	// sidebar highlight: which admin menu (and group) owns the current page.
	// Computed centrally so every present and future page is covered without
	// touching individual handlers.
	activeMenu, activeParent := ActiveAdminMenu(c.Request.URL.Path, c.Query("today"))
	data["ACTIVE_MENU"] = activeMenu
	data["ACTIVE_PARENT"] = activeParent

	menu, _ := c.Get("menu") //We load the menu using the LoadMenu() middleware and ignore the ok variable because if there is any error in LoadMenu(), a 500 error will be returned
	data["MENU"] = menu

	//check `auth` key from context
	customer, ok := helpers.GetAuthUser(c)
	if ok && customer.ID > 0 {
		data["AUTH"] = customer
	}

	return data
}

// oldFields flattens the flashed old input to one value per key.
func oldFields(c *gin.Context) map[string]string {
	out := map[string]string{}
	raw := sessions.GET(c, "olds")
	if raw == "" {
		return out
	}
	var form map[string][]string
	if err := json.Unmarshal([]byte(raw), &form); err != nil {
		return out
	}
	for key, vals := range form {
		if len(vals) > 0 {
			out[key] = vals[0]
		}
	}
	return out
}

// ActiveAdminMenu maps a request path to the admin sidebar entry that owns
// it: exact create-pages first, then their section (list/show/edit and the
// legacy per-product operation pages fall back to the list entry).
// The "today" query flag distinguishes سفارشات امروز from لیست سفارشات.
// Unknown paths return ""/"" (nothing highlighted).
func ActiveAdminMenu(path, today string) (menu, parent string) {
	switch {
	case path == "/admins/home":
		return "home", ""
	case path == "/admins/monitoring" || hasPathPrefix(path, "/admins/monitoring/"):
		return "monitoring", ""

	// banners
	case path == "/admins/banners/create":
		return "banners-create", "banners"
	case path == "/admins/banners" || hasPathPrefix(path, "/admins/banners/"):
		return "banners-list", "banners"

	// product sliders
	case path == "/admins/sliders/create":
		return "sliders-create", "sliders"
	case path == "/admins/sliders" || hasPathPrefix(path, "/admins/sliders/"):
		return "sliders-list", "sliders"

	// brands
	case path == "/admins/brands/create":
		return "brands-create", "brands"
	case path == "/admins/brands" || hasPathPrefix(path, "/admins/brands/"):
		return "brands-list", "brands"

	// categories
	case path == "/admins/categories/create":
		return "categories-create", "categories"
	case path == "/admins/categories" || hasPathPrefix(path, "/admins/categories/"):
		return "categories-list", "categories"

	// attributes
	case path == "/admins/attributes/create":
		return "attributes-create", "attributes"
	case path == "/admins/attributes" || hasPathPrefix(path, "/admins/attributes/"):
		return "attributes-list", "attributes"

	// attribute values
	case path == "/admins/attribute-values/create":
		return "attr-values-create", "attr-values"
	case path == "/admins/attribute-values" || hasPathPrefix(path, "/admins/attribute-values/") ||
		hasPathPrefix(path, "/admins/attribute/values/"):
		return "attr-values-list", "attr-values"

	// products (list owns every per-product page: show/edit/gallery/
	// features/inventory/images/recommendations and the legacy routes)
	case path == "/admins/products/create":
		return "products-create", "products"
	case path == "/admins/products" || hasPathPrefix(path, "/admins/products/") ||
		hasPathPrefix(path, "/admins/inventories/") ||
		hasPathPrefix(path, "/admins/product-inventory-attributes/") ||
		hasPathPrefix(path, "/admins/products-attributes/") ||
		hasPathPrefix(path, "/admins/products/images/"):
		return "products-list", "products"

	// customers
	case path == "/admins/customers" || hasPathPrefix(path, "/admins/customers/"):
		return "customers-list", "customers"

	// orders (?today=1 is the daily list)
	case (path == "/admins/orders" || hasPathPrefix(path, "/admins/orders/")) && today == "1":
		return "orders-today", "orders"
	case path == "/admins/orders" || hasPathPrefix(path, "/admins/orders/"):
		return "orders-list", "orders"

	// order fees (no separate create entry — every fee page highlights its kind)
	case hasPathPrefix(path, "/admins/fees/shipping"):
		return "fees-shipping", "fees"
	case hasPathPrefix(path, "/admins/fees/packaging"):
		return "fees-packaging", "fees"
	case path == "/admins/fees" || hasPathPrefix(path, "/admins/fees/"):
		return "fees-shipping", "fees"
	}
	return "", ""
}

// hasPathPrefix is strings.HasPrefix kept local so templates never need a
// custom func for it (the mapper above is unit-tested instead).
func hasPathPrefix(path, prefix string) bool {
	if len(path) < len(prefix) {
		return false
	}
	return path[:len(prefix)] == prefix
}
