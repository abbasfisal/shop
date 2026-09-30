package responses

import "shop/domain/entities"

// SiteBanner is the admin/storefront view of a fixed site banner (header and
// the reserved image-hero placement).
type SiteBanner struct {
	ID             uint
	Placement      string
	PlacementLabel string
	Title          string
	Link           string
	Image          string
	MobileImage    string
	Status         bool
	StatusText     string
	StartsAt       string
	EndsAt         string
	Schedule       string
	SortOrder      int
	CreatedAt      string
}

type SiteBanners struct {
	Data []SiteBanner
}

// MobileSrc is the image the mobile slot must show: the mobile upload, or the
// desktop one when the admin never uploaded a mobile variant.
// Value receiver on purpose — templates call it on a map value.
func (b SiteBanner) MobileSrc() string {
	if b.MobileImage != "" {
		return b.MobileImage
	}
	return b.Image
}

// HasLink reports whether the banner is clickable (empty link → no href).
func (b SiteBanner) HasLink() bool {
	return b.Link != ""
}

func ToSiteBanner(b *entities.SiteBanner) *SiteBanner {
	if b == nil {
		return nil
	}

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

	return &SiteBanner{
		ID:             b.ID,
		Placement:      b.Placement,
		PlacementLabel: entities.SiteBannerPlacementLabel(b.Placement),
		Title:          b.Title,
		Link:           b.Link,
		Image:          b.Image,
		MobileImage:    b.MobileImage,
		Status:         b.Status,
		StatusText:     statusText,
		StartsAt:       startsAt,
		EndsAt:         endsAt,
		Schedule:       schedule,
		SortOrder:      b.SortOrder,
		CreatedAt:      b.CreatedAt.Format("2006-01-02 15:04"),
	}
}

// ToSiteBanners is nil-safe (an empty feed must render, not panic).
func ToSiteBanners(banners []*entities.SiteBanner) *SiteBanners {
	if banners == nil {
		return &SiteBanners{}
	}
	out := &SiteBanners{Data: make([]SiteBanner, 0, len(banners))}
	for _, b := range banners {
		out.Data = append(out.Data, *ToSiteBanner(b))
	}
	return out
}
