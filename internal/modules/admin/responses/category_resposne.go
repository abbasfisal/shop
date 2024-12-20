package responses

import "shop/internal/entities"

type Category struct {
	ID       uint
	ParentID *uint
	Priority *uint
	Title    string
	Slug     string
	Image    string
	Status   bool
}
type Categories struct {
	Data []Category
}

func ToCategory(category *entities.Category) *Category {
	return &Category{
		ID:       category.ID,
		ParentID: category.ParentID,
		Priority: category.Priority,
		Title:    category.Title,
		Slug:     category.Slug,
		Image:    category.Image,
		Status:   category.Status,
	}
}

func ToCategories(categories []*entities.Category) *Categories {

	if categories == nil {
		return nil
	}

	response := Categories{
		Data: make([]Category, 0, len(categories)),
	}

	for _, cat := range categories {
		response.Data = append(response.Data, *ToCategory(cat))
	}

	return &response
}
