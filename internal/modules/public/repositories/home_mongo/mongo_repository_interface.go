package home_mongo

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"shop/internal/entities"
	"shop/internal/modules/public/requests"
)

type MongoHomeRepositoryInterface interface {
	GetProduct(c *gin.Context, productSku string, productSlug string) (map[string]interface{}, error)
	GetProductByObjectID(c *gin.Context, productObjectID primitive.ObjectID, req requests.AddToCartRequest) (entities.MongoProduct, error)
}
