package responses

import (
	"shop/internal/entities"
)

type CategoryResponse struct {
	ID       uint
	Priority *uint
	Title    string
	Slug     string
	ParentID *uint
	Status   bool

	SubCategories []*CategoryResponse `json:"SubCategories,omitempty"`
}

func ToMenuResponse(category *entities.Category) *CategoryResponse {
	var subCategories []*CategoryResponse

	for _, subCategory := range category.SubCategories {
		subCategories = append(subCategories, ToMenuResponse(subCategory))
	}

	return &CategoryResponse{
		ID:            category.ID,
		Priority:      category.Priority,
		Title:         category.Title,
		Slug:          category.Slug,
		ParentID:      category.ParentID,
		Status:        category.Status,
		SubCategories: subCategories,
	}
}
