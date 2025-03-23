package product

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"shop/internal/database/mongodb"
	"time"
)

func (p *ProductRepository) GetAllMongoProduct(c context.Context) ([]bson.M, error) {

	cc := mongodb.GetCollection(mongodb.ProductsCollection)

	// تنظیم فیلتر (مثلاً دریافت همه اسناد)
	//filter := bson.M{}

	// تنظیم Projection برای دریافت فقط فیلدهای موردنیاز
	projection := bson.M{
		"_id":                    1, // دریافت id
		"product.id":             1,
		"product.title":          1,
		"product.slug":           1,
		"product.original_price": 1,
		"product.sale_price":     1,
	}

	// اجرای کوئری با فیلتر و Projection
	cursor, err := cc.Find(c, bson.D{}, options.Find().SetProjection(projection))

	if err != nil {
		return []bson.M{}, err
	}
	defer cursor.Close(c)

	// خواندن داده‌ها
	var products []bson.M
	if err = cursor.All(c, &products); err != nil {
		return products, err
	}

	return products, nil
}

func (p *ProductRepository) InsertRecommendation(c *gin.Context, productID int, productRecommendationIDs []string) error {

	// get mongo document by product id
	collection := mongodb.GetCollection(mongodb.RecommendationCollection)

	filter := bson.M{"product_id": productID}
	update := bson.M{
		"$set": bson.M{
			"product_recommendations": productRecommendationIDs,
			"updated_at":              time.Now(),
		},
	}

	// گزینه upsert: اگر داکیومنت موجود نبود، یک داکیومنت جدید ایجاد کند.
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(c, filter, update, opts)
	if err != nil {
		log.Println("--- upsert recommendation products failed :", err)
		return err
	}

	log.Println("--- successfully upsert recommendation")
	return nil
}
