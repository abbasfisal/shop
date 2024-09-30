package home_mongo

import "github.com/gin-gonic/gin"

type MongoHomeRepositoryInterface interface {
	GetProduct(c *gin.Context, productSku string, productSlug string) (map[string]interface{}, error)
}
