package responses

import "shop/internal/entities"

type Brand struct {
	ID    uint
	Title string
	Slug  string
	Image string
}
type Brands struct {
	Data []Brand
}

func ToBrand(brand entities.Brand) Brand {
	return Brand{
		ID:    brand.ID,
		Title: brand.Title,
		Slug:  brand.Slug,
		Image: brand.Image,
	}
}

func ToBrands(brands []entities.Brand) Brands {
	var response Brands

	for _, brand := range brands {
		response.Data = append(response.Data, ToBrand(brand))
	}

	return response
}
