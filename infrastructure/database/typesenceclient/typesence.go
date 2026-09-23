package typesenceclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
			typesense.WithAPIKey(os.Getenv("TYPESENCE_API_KEY")),
		)

		if err := CreateSchema(tClient); err != nil {
			log.Println("Error creating schema: ", err)
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

	createSampleDocument(client)

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

func createSampleDocument(client *typesense.Client) {
	//create a doc just for testing
	if false {

		p := struct {
			ID    string `json:"id,omitempty"`
			Title string `json:"title,omitempty"`
			Slug  string `json:"slug,omitempty"`
			SKU   string `json:"sku,omitempty"`
			Price int    `json:"price,omitempty"`
		}{
			ID:    uuid.New().String(),
			Title: "کتاب شب های برره محسن چاوشی ۳۰ 9898",
			Slug:  "کتاب-شب-های-برره-محسن-چاوشی-۳۰-9898",
			SKU:   "sku888",
			Price: 5000,
		}
		createDoc, err := client.Collection("products").Documents().
			Create(context.Background(), p, &api.DocumentIndexParameters{})
		if err != nil {
			log.Fatal("create doc failed:", err)
		}
		log.Println("succ create doc:", createDoc)
	}
}
