package mongodb

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"sync"
)

var (
	mongoClient *mongo.Client
	clientOnce  sync.Once
)

// collection names

const (
	ProductsCollection = "products"
)

func Connect() {
	clientOnce.Do(func() {
		var err error
		uri := fmt.Sprintf("mongodb://%s:%s/", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))
		clientOptions := options.Client().ApplyURI(uri)
		mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatalln("[x] Error creating MongoDB client:", err)
		}

		err = mongoClient.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatalln("[x] Error pinging MongoDB: ", err)
		}
	})

	fmt.Println("\n[x] connected to MongoDB successfully")

}
func Get() *mongo.Client {
	return mongoClient
}

// GetCollection retrieves the specified MongoDB collection.
func GetCollection(CollectionName string) *mongo.Collection {
	if mongoClient == nil {
		log.Fatal("[x] MongoDB client is not initialized. Did you call Connect?")
	}

	// Get the database name from viper and return the collection
	dbName := viper.GetString("mongodb.name")
	if dbName == "" {
		log.Fatal("MongoDB database name not found in config")
	}

	return mongoClient.Database(dbName).Collection(CollectionName)
}

// Disconnect closes the MongoDB connection gracefully.
func Disconnect() error {
	if mongoClient == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}

	err := mongoClient.Disconnect(context.TODO())
	if err != nil {
		log.Println("Error disconnecting MongoDB: ", err)
		return err
	}

	fmt.Println("[mongo] successfully disconnected from MongoDB")
	return nil
}
