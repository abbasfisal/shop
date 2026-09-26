package response

import "testing"

func TestActiveAdminMenu(t *testing.T) {
	cases := []struct {
		path, today string
		menu, group string
	}{
		{"/admins/home", "", "home", ""},
		{"/admins/monitoring", "", "monitoring", ""},
		{"/admins/monitoring/", "", "monitoring", ""},

		{"/admins/banners", "", "banners-list", "banners"},
		{"/admins/banners/create", "", "banners-create", "banners"},
		{"/admins/banners/3/edit", "", "banners-list", "banners"},

		{"/admins/sliders", "", "sliders-list", "sliders"},
		{"/admins/sliders/create", "", "sliders-create", "sliders"},
		{"/admins/sliders/1/edit", "", "sliders-list", "sliders"},

		{"/admins/brands", "", "brands-list", "brands"},
		{"/admins/brands/create", "", "brands-create", "brands"},
		{"/admins/brands/1/edit", "", "brands-list", "brands"},
		{"/admins/brands/1/products", "", "brands-list", "brands"},

		{"/admins/categories", "", "categories-list", "categories"},
		{"/admins/categories/create", "", "categories-create", "categories"},
		{"/admins/categories/1/products", "", "categories-list", "categories"},

		{"/admins/attributes", "", "attributes-list", "attributes"},
		{"/admins/attributes/create", "", "attributes-create", "attributes"},
		{"/admins/attributes/1/edit", "", "attributes-list", "attributes"},

		{"/admins/attribute-values", "", "attr-values-list", "attr-values"},
		{"/admins/attribute-values/create", "", "attr-values-create", "attr-values"},
		{"/admins/attribute-values/6/edit", "", "attr-values-list", "attr-values"},
		{"/admins/attribute/values/2/show", "", "attr-values-list", "attr-values"},

		{"/admins/products", "", "products-list", "products"},
		{"/admins/products/create", "", "products-create", "products"},
		{"/admins/products/205", "", "products-list", "products"},
		{"/admins/products/205/edit", "", "products-list", "products"},
		{"/admins/products/205/show-gallery", "", "products-list", "products"},
		{"/admins/products/205/add-attributes", "", "products-list", "products"},
		{"/admins/products/205/add-inventory", "", "products-list", "products"},
		{"/admins/products/205/add-feature", "", "products-list", "products"},
		{"/admins/products/images/9/delete", "", "products-list", "products"},
		{"/admins/inventories/209/delete", "", "products-list", "products"},
		{"/admins/product-inventory-attributes/11/delete", "", "products-list", "products"},
		{"/admins/products-attributes/11/delete", "", "products-list", "products"},

		{"/admins/customers", "", "customers-list", "customers"},

		{"/admins/orders", "", "orders-list", "orders"},
		{"/admins/orders", "1", "orders-today", "orders"},
		{"/admins/orders/9/details", "", "orders-list", "orders"},
		{"/admins/orders/9/details", "1", "orders-today", "orders"},

		{"/admins/fees/shipping", "", "fees-shipping", "fees"},
		{"/admins/fees/shipping/create", "", "fees-shipping", "fees"},
		{"/admins/fees/shipping/20/edit", "", "fees-shipping", "fees"},
		{"/admins/fees/packaging", "", "fees-packaging", "fees"},
		{"/admins/fees/packaging/create", "", "fees-packaging", "fees"},

		{"/admins/login", "", "", ""},
		{"/something-else", "", "", ""},
		{"/", "", "", ""},
	}

	for _, tc := range cases {
		menu, group := ActiveAdminMenu(tc.path, tc.today)
		if menu != tc.menu || group != tc.group {
			t.Errorf("path %q today %q: got (%q,%q) want (%q,%q)",
				tc.path, tc.today, menu, group, tc.menu, tc.group)
		}
	}
}
