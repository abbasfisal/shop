package typesenceclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/typesense/typesense-go/v3/typesense"
	"github.com/typesense/typesense-go/v3/typesense/api"
	"log"
	"os"
	"sync"
	"time"
)

var (
	once    sync.Once
	tClient *typesense.Client
)

func Connect() {
	once.Do(func() {
		url := fmt.Sprintf("http://%s:%s", os.Getenv("TYPESENCE_HOST"), os.Getenv("TYPESENCE_PORT"))

		tClient = typesense.NewClient(
			typesense.WithServer(url),
			typesense.WithAPIKey(os.Getenv("APP_KEY")),
		)

		if err := CreateSchema(tClient); err != nil {
			log.Fatal("Error creating schema: ", err)
		}
	})

	typesenceHealth, tErr := tClient.Health(context.TODO(), time.Second)
	if tErr != nil {
		log.Fatal("Typesense health check failed: ", tErr)
	}

	log.Println("Typesense health check passed: ", typesenceHealth)
}

func GetTClient() *typesense.Client {
	return tClient
}

func CreateSchema(client *typesense.Client) error {

	_, err := client.Collection("products").Retrieve(context.TODO())
	if err == nil {
		log.Println("Collection 'products' already exists")
		return nil
	}

	if err != nil && !errors.Is(err, &typesense.HTTPError{}) {
		return fmt.Errorf("unexpected error checking collection: %v", err)
	}

	create, err := client.Collections().Create(context.TODO(), &api.CollectionSchema{
		Name: "products",
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "title", Type: "string"},
			{Name: "slug", Type: "string"},
			{Name: "sku", Type: "string"},
		},
		TokenSeparators: &[]string{" ", "-", ".", ",", ":"},
	})

	if err != nil {
		var httpErr *typesense.HTTPError
		if errors.As(err, &httpErr) {
			switch httpErr.Status {
			case 404:
				log.Println("Typesense Collection not found")
			case 409:
				log.Println("Typesense Collection already exists")
			default:
				log.Println("Unexpected error while creating collection: ", err)
			}
		} else {
			log.Println("Unexpected error format while creating collection: ", err)
		}
		return err
	}

	log.Printf("Typesense collection '%s' created successfully", create.Name)
	return nil
}
