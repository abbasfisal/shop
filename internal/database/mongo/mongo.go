package mongo

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var client *mongo.Client

func Connect() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("mongo error creating client: ", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("mongo error pinging database: ", err)
	}

	fmt.Println("\n[mongo] connected to MongoDB successfully")

}
func Get() *mongo.Client {
	return client
}

func GetCollection(CollectionName string) *mongo.Collection {
	return client.Database(viper.GetString("mongodb.name")).Collection(CollectionName)
}
