package responses

import "shop/domain/entities"

// Banner is the admin list view of a promotion banner.
type Banner struct {
	ID          uint
	Title       string
	Layout      string
	LayoutLabel string
	Link        string
	Status      bool
	StatusText  string
	StartsAt    string
	EndsAt      string
	Schedule    string
	SortOrder   int
	Priority    uint
	Type        uint
	TypeLabel   string
	Image       string
	CreatedAt   string
}

type Banners struct {
	Data []Banner
}

func ToBanner(b *entities.Banner) *Banner {
	statusText := "غیرفعال"
	if b.Status {
		statusText = "فعال"
	}
	startsAt, endsAt := "", ""
	if b.StartsAt != nil {
		startsAt = b.StartsAt.Format("2006-01-02")
	}
	if b.EndsAt != nil {
		endsAt = b.EndsAt.Format("2006-01-02")
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

	return &Banner{
		ID:          b.ID,
		Title:       b.Title,
		Layout:      b.Layout,
		LayoutLabel: entities.BannerLayoutLabel(b.Layout),
		Link:        b.Link,
		Status:      b.Status,
		StatusText:  statusText,
		StartsAt:    startsAt,
		EndsAt:      endsAt,
		Schedule:    schedule,
		SortOrder:   b.SortOrder,
		Priority:    b.Priority,
		Type:        b.Type,
		TypeLabel:   entities.BannerTypeLabel(b.Type),
		Image:       b.Image,
		CreatedAt:   b.CreatedAt.Format("2006-01-02 15:04"),
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
