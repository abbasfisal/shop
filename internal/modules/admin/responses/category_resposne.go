package responses

import "shop/internal/entities"

type Category struct {
	ID     uint
	Title  string
	Slug   string
	Image  string
	Status bool
}
type Categories struct {
	Data []Category
}

func ToCategory(category entities.Category) Category {
	return Category{
		ID:     category.ID,
		Title:  category.Title,
		Slug:   category.Slug,
		Image:  category.Image,
		Status: category.Status,
	}
}

func ToCategories(categories []entities.Category) Categories {
	var response Categories

	for _, cat := range categories {
		response.Data = append(response.Data, ToCategory(cat))
	}

	return response
}
