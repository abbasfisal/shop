package responses

import "shop/domain/entities"

// Banner is the admin list view of a banner (labelled type + status).
type Banner struct {
	ID         uint
	Type       uint
	TypeLabel  string
	Link       string
	Priority   uint
	Status     bool
	StatusText string
	Image      string
	CreatedAt  string
}

type Banners struct {
	Data []Banner
}

func ToBanner(b *entities.Banner) *Banner {
	statusText := "غیرفعال"
	if b.Status {
		statusText = "فعال"
	}
	return &Banner{
		ID:         b.ID,
		Type:       b.Type,
		TypeLabel:  entities.BannerTypeLabel(b.Type),
		Link:       b.Link,
		Priority:   b.Priority,
		Status:     b.Status,
		StatusText: statusText,
		Image:      b.Image,
		CreatedAt:  b.CreatedAt.Format("2006-01-02 15:04"),
	}
}

func ToBanners(banners []*entities.Banner) *Banners {
	if banners == nil {
		return &Banners{}
	}
	out := &Banners{Data: make([]Banner, 0, len(banners))}
	for _, b := range banners {
		out.Data = append(out.Data, *ToBanner(b))
	}
	return out
}
