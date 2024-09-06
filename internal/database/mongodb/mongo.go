package mongodb

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var client *mongo.Client

// collection names

const (
	ProductsCollection = "products"
)

func Connect() error {

	var err error
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("mongo error creating client: ", err)
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("mongo error pinging database: ", err)
		return err
	}

	fmt.Println("\n[mongo] connected to MongoDB successfully")
	return nil
}
func Get() *mongo.Client {
	return client
}

func GetCollection(CollectionName string) *mongo.Collection {
	if client == nil {
		log.Fatal("MongoDB client is not initialized. Did you call Connect?")
	}

	// Get the database name from viper and return the collection
	dbName := viper.GetString("mongodb.name")
	if dbName == "" {
		log.Fatal("MongoDB database name not found in config")
	}

	return client.Database(dbName).Collection(CollectionName)
}
