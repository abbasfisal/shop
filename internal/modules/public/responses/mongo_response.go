package responses

import "shop/internal/entities"

func ToMongoProductResponse(mongoProduct entities.MongoProduct) map[string]interface{} {
	return map[string]interface{}{
		"_id":         mongoProduct.ID.Hex(),
		"product":     mongoProduct.Product,
		"inventories": mongoProduct.Inventories,
	}
}
