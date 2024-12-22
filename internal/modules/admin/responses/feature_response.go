package responses

import (
	"shop/internal/entities"
)

type Feature struct {
	ID        uint
	ProductID uint
	Title     string
	Value     string
}
type Features struct {
	Data []Feature
}

func ToFeature(f *entities.Feature) *Feature {
	return &Feature{
		ID:        f.ID,
		ProductID: f.ProductID,
		Title:     f.Title,
		Value:     f.Value,
	}
}

func ToFeatures(f []*entities.Feature) *Features {
	if f == nil {
		return &Features{}
	}

	var fResponse Features
	for _, featureItem := range f {
		fResponse.Data = append(fResponse.Data, *ToFeature(featureItem))
	}
	return &fResponse
}
