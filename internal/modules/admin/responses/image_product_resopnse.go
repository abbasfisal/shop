package responses

import (
	"github.com/spf13/viper"
	"shop/internal/entities"
)

type ImageProduct struct {
	ID           uint
	OriginalPath string
	FullPath     string
}

type ImageProducts struct {
	Data []ImageProduct
}

func ToImageProducts(productImages []*entities.ProductImages) *ImageProducts {
	if productImages == nil {
		return &ImageProducts{}
	}

	var pResponse ImageProducts
	for _, p := range productImages {
		pResponse.Data = append(pResponse.Data, *ToImageProduct(p))
	}
	return &pResponse
}

func ToImageProduct(p *entities.ProductImages) *ImageProduct {
	return &ImageProduct{
		ID:           p.ID,
		OriginalPath: p.Path,
		FullPath:     viper.GetString("Upload.Products") + p.Path,
	}
}
