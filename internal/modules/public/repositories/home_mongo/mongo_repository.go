package home_mongo

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"shop/internal/database/mongodb"
	"shop/internal/entities"
	"shop/internal/modules/public/requests"
	"shop/internal/modules/public/responses"
	"shop/internal/pkg/custom_error"
)

type MongoHomeRepository struct {
}

func NewMongoRepository() MongoHomeRepository {
	return MongoHomeRepository{}
}

func (m MongoHomeRepository) GetProduct(c *gin.Context, productSku string, productSlug string) (map[string]interface{}, error) {
	productionCollection := mongodb.GetCollection(mongodb.ProductsCollection)
	filter := bson.D{
		{"product.sku", productSku},
		{"product.slug", productSlug},
	}
	var mongoProduct entities.MongoProduct
	FindProductErr := productionCollection.FindOne(c, filter).Decode(&mongoProduct)
	if FindProductErr != nil {
		if FindProductErr == mongo.ErrNoDocuments {
			// محصولی یافت نشد
			fmt.Println("--هیچ سندی برای محصول با اسلاگ ", productSlug, " و SKU :", productSku, " یافت نشد.")
			return nil, errors.New(custom_error.RecordNotFound)
		}
		fmt.Println("~~~~ error while finding doc :÷÷÷÷÷", FindProductErr)
		return nil, errors.New(custom_error.SomethingWrongHappened)
	}

	fmt.Println("~~~~~~ document find succ ~~~~~~~~")
	fmt.Printf("%+v", mongoProduct)
	return responses.ToMongoProductResponse(mongoProduct), nil

}

func (m MongoHomeRepository) GetProductByObjectID(c *gin.Context, productObjectID primitive.ObjectID, req requests.AddToCartRequest) (entities.MongoProduct, error) {
	productionCollection := mongodb.GetCollection(mongodb.ProductsCollection)
	filter := bson.D{
		{"_id", productObjectID},
	}
	var mongoProduct entities.MongoProduct
	FindProductErr := productionCollection.FindOne(c, filter).Decode(&mongoProduct)

	if FindProductErr != nil {
		if FindProductErr == mongo.ErrNoDocuments {
			// محصولی یافت نشد
			fmt.Println("[Error Not Found Document ]-[GetProductByObjectID]:err:", custom_error.RecordNotFound)
			return mongoProduct, errors.New(custom_error.RecordNotFound)
		}
		fmt.Println("[Error]-[GetProductByObjectID]:err:", custom_error.SomethingWrongHappened)

		return mongoProduct, errors.New(custom_error.SomethingWrongHappened)
	}

	fmt.Println("[success]-[GetProductByObjectID]-msg:record found")

	return mongoProduct, nil
}
