package home_mongo

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

func (m MongoHomeRepository) GetProduct(c *gin.Context, productSku string, productSlug string) (map[string]interface{}, []entities.MongoProductRecommendation, error) {
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
			log.Println("product-slug : ", productSlug, " | SKU :", productSku, " not found.")
			return nil, nil, errors.New(custom_error.RecordNotFound)
		}
		return nil, nil, errors.New(custom_error.SomethingWrongHappened)
	}

	// get recommendation products
	var recommendedProducts []entities.MongoProductRecommendation
	if true {

		recommendationsCollection := mongodb.GetCollection(mongodb.RecommendationCollection)

		//1. findOne recommendation by product_id
		var product bson.M
		err := recommendationsCollection.FindOne(c, bson.M{"product_id": mongoProduct.Product.ID}).Decode(&product)
		if err != nil {
			log.Println("--- recommendation findOne query failed , productID: ", mongoProduct.Product.ID, " | err:", err)
			return responses.ToMongoProductResponse(mongoProduct), recommendedProducts, nil
		}
		log.Println("--- product :", product)

		//2. convert to bson.A
		recommendations, ok := product["product_recommendations"].(bson.A)
		if !ok || len(recommendations) == 0 {
			return responses.ToMongoProductResponse(mongoProduct), recommendedProducts, nil
		}
		log.Println("--- recommendation:", recommendations)

		//3. because we stored product ids in string format we have to convert to ObjectID
		// to work properly our Find query
		var objectIDs []primitive.ObjectID
		for _, id := range recommendations {
			strID, ok := id.(string)
			if !ok {
				continue
			}
			objID, err := primitive.ObjectIDFromHex(strID)
			if err != nil {
				continue
			}
			objectIDs = append(objectIDs, objID)
		}
		log.Println("--- objectIDS:", objectIDs)

		//4. Find products(in products mongoDB) by recommendation IDS from step(1)
		productColl := mongodb.GetCollection(mongodb.ProductsCollection)
		cursor, err := productColl.Find(c, bson.M{"_id": bson.M{"$in": objectIDs}},
			options.Find().SetProjection(bson.M{
				"product.Category":       1,
				"product.title":          1,
				"product.slug":           1,
				"product.sku":            1,
				"product.Images":         1,
				"product.original_price": 1,
				"product.sale_price":     1,
				"product.Discount":       1,
				"product.Status":         1,
				"_id":                    1, // to delete switch 1 with 0
			}),
		)
		if err != nil {
			log.Println("--- Find products by recommendation IDS err:", err)
			return responses.ToMongoProductResponse(mongoProduct), recommendedProducts, nil
		}
		defer cursor.Close(c)

		//5. decode data into recommendedProducts
		if err = cursor.All(c, &recommendedProducts); err != nil {
			log.Println("--- cursor error:", err)
			return responses.ToMongoProductResponse(mongoProduct), recommendedProducts, nil
		}
	}

	return responses.ToMongoProductResponse(mongoProduct), recommendedProducts, nil
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
