package responses

import "shop/internal/entities"

type Category struct {
	ID           uint
	ParentID     *uint
	UintParentID uint //  این فیلد رو برای این اضافه کردیم که درد حالت مقایسه در صفحه اچ تی ام ال به مشکل نخوریم (در قسمت ؛سطح؛) چون اشاره گر با متغیری که اشاره گر نیست قابل مقایسه نیست :)
	Priority     *uint
	Title        string
	Slug         string
	Image        string
	Status       bool
}
type Categories struct {
	Data []Category
}

func ToCategory(category *entities.Category) *Category {

	// read UintParentID doc ☝
	var parentID uint
	if category.ParentID != nil {
		parentID = *category.ParentID
	} else {
		parentID = 0 // مقدار پیش‌فرض
	}

	return &Category{
		ID:           category.ID,
		ParentID:     category.ParentID,
		UintParentID: parentID,
		Priority:     category.Priority,
		Title:        category.Title,
		Slug:         category.Slug,
		Image:        category.Image,
		Status:       category.Status,
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
