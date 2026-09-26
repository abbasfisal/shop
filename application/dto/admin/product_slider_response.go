package responses

import (
	"shop/domain/entities"
)

// SliderPosition is a labelled homepage slot for the admin forms.
type SliderPosition struct {
	Value string
	Label string
}

// SliderPositions returns every slot in rendering order.
func SliderPositions() []SliderPosition {
	src := entities.SliderPositions()
	out := make([]SliderPosition, 0, len(src))
	for _, p := range src {
		out = append(out, SliderPosition{Value: p.Value, Label: p.Label})
	}
	return out
}

// SliderProduct is one product inside a slider (admin form + list view).
type SliderProduct struct {
	ProductID       uint   `json:"product_id"`
	Title           string `json:"title"`
	Sku             string `json:"sku"`
	Slug            string `json:"slug"`
	Image           string `json:"image"`
	Price           uint   `json:"price"`
	OriginalPrice   uint   `json:"original_price"`
	DiscountPercent uint   `json:"discount_percent"`
	InStock         bool   `json:"in_stock"`
	SortOrder       int    `json:"sort_order"`
}

// ProductURL is the storefront permalink of a slider product.
func (p SliderProduct) ProductURL() string {
	return "/product/" + p.Sku + "/" + p.Slug
}

// ProductSlider is the admin view of a homepage product slider.
type ProductSlider struct {
	ID            uint
	Title         string
	Slug          string
	Position      string
	PositionLabel string
	Status        string
	StatusText    string
	StartsAt      string
	EndsAt        string
	Schedule      string
	CategoryID    uint
	CategoryTitle string
	ProductsCount int
	Products      []SliderProduct
	SelectedIDs   []uint
	CatalogURL    string
	HomeURL       string
	// MediaPath is the uploaded-image prefix the storefront partial needs
	// (the partial only receives this struct as its data).
	MediaPath string
}

type ProductSliders struct {
	Data []ProductSlider
}

func ToSliderProduct(link *entities.SliderProduct) SliderProduct {
	out := SliderProduct{
		ProductID: link.ProductID,
		SortOrder: link.SortOrder,
	}
	if link.Product != nil {
		view := ToProduct(link.Product)
		out.Title = view.Title
		out.Sku = view.Sku
		out.Slug = view.Slug
		out.InStock = view.InStock
		out.OriginalPrice = view.OriginalPrice
		out.DiscountPercent = view.Discount

		// selling price: the aggregate minimum when variants exist
		out.Price = view.MinPrice
		if out.Price == 0 {
			out.Price = view.SalePrice
		}
		if view.Images != nil && len(view.Images.Data) > 0 {
			out.Image = view.Images.Data[0].OriginalPath
		}
	}
	return out
}

func ToSlider(s *entities.ProductSlider) *ProductSlider {
	if s == nil {
		return &ProductSlider{}
	}

	startsAt, endsAt := "", ""
	if s.StartsAt != nil {
		startsAt = s.StartsAt.Format("2006-01-02")
	}
	if s.EndsAt != nil {
		endsAt = s.EndsAt.Format("2006-01-02")
	}
	schedule := "بدون محدودیت"
	switch {
	case startsAt != "" && endsAt != "":
		schedule = startsAt + " تا " + endsAt
	case startsAt != "":
		schedule = "از " + startsAt
	case endsAt != "":
		schedule = "تا " + endsAt
	}

	statusText := "پیش‌نویس"
	switch s.Status {
	case entities.ProductStatusPublished:
		statusText = "منتشر شده"
	case entities.ProductStatusArchived:
		statusText = "بایگانی"
	}

	out := &ProductSlider{
		ID:            s.ID,
		Title:         s.Title,
		Slug:          s.Slug,
		Position:      s.Position,
		PositionLabel: entities.SliderPositionLabel(s.Position),
		Status:        s.Status,
		StatusText:    statusText,
		StartsAt:      startsAt,
		EndsAt:        endsAt,
		Schedule:      schedule,
		ProductsCount: len(s.Products),
		Products:      make([]SliderProduct, 0, len(s.Products)),
		CatalogURL:    "/sliders/" + s.Slug,
		HomeURL:       "/",
	}
	if s.Category != nil {
		out.CategoryID = s.Category.ID
		out.CategoryTitle = s.Category.Title
	}
	for _, link := range s.Products {
		out.Products = append(out.Products, ToSliderProduct(link))
		out.SelectedIDs = append(out.SelectedIDs, link.ProductID)
	}
	return out
}

func ToSliders(sliders []*entities.ProductSlider) *ProductSliders {
	out := &ProductSliders{Data: make([]ProductSlider, 0, len(sliders))}
	for _, s := range sliders {
		out.Data = append(out.Data, *ToSlider(s))
	}
	return out
}

// SlidersByPosition groups the storefront feed by homepage slot so the home
// template can render each section where it belongs.
func SlidersByPosition(sliders []*entities.ProductSlider) map[string][]ProductSlider {
	out := map[string][]ProductSlider{}
	for _, s := range sliders {
		view := ToSlider(s)
		out[s.Position] = append(out[s.Position], *view)
	}
	return out
}
