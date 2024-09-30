package responses

import "shop/internal/entities"

func ToMongoProductResponse(mongoProduct entities.MongoProduct) map[string]interface{} {
	return map[string]interface{}{
		"product":     mongoProduct.Product,
		"inventories": mongoProduct.Inventories,
	}
}
