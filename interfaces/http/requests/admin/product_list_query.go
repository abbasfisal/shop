package requests

// ProductListQuery is the admin product list filter set (Laravel
// ProductController::index: q / category_id / in_stock / attr / sort).
type ProductListQuery struct {
	Q          string            // title / sku ILIKE search
	Status     string            // draft|published|archived — empty = every status
	CategoryID uint              // 0 = every category
	InStock    bool              // only products with available stock
	Attrs      map[string][]uint // attribute code → value ids (OR inside, AND across)
	Sort       string            // newest | price_asc | price_desc | name
}
