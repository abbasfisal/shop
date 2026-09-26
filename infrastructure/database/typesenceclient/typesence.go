package typesenceclient

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/typesense/typesense-go/v3/typesense"
	"github.com/typesense/typesense-go/v3/typesense/api"
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
		// typesense is optional at boot: search degrades gracefully and
		// product upserts skip until the engine is reachable again.
		log.Println("[typesense] health check failed (search unavailable): ", tErr)
		return
	}

	log.Println("Typesense health check passed: ", typesenceHealth)
}

func GetTClient() *typesense.Client {
	return tClient
}

func boolPtr(b bool) *bool { return &b }

// ProductsCollection is the Typesense collection that mirrors the shop catalog.
const ProductsCollection = "products"

// CreateSchema creates the rich products collection schema (title/sku search,
// category/brand facets, price/stock fields for realtime search UI).
func CreateSchema(client *typesense.Client) error {
	facet := boolPtr(true)
	create, err := client.Collections().Create(context.TODO(), &api.CollectionSchema{
		Name: ProductsCollection,
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "title", Type: "string"},
			{Name: "slug", Type: "string"},
			{Name: "sku", Type: "string"},
			{Name: "description", Type: "string"},
			{Name: "category", Type: "string", Facet: facet},
			{Name: "brand", Type: "string", Facet: facet},
			{Name: "original_price", Type: "int32"},
			{Name: "sale_price", Type: "int32"},
			{Name: "discount", Type: "int32"},
			{Name: "stock", Type: "int32"},
			{Name: "in_stock", Type: "bool", Facet: facet},
			{Name: "status", Type: "string", Facet: facet},
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

// RecreateSchema drops the products collection (if present) and creates the
// current schema again — used by `search:reindex --recreate` after schema changes.
func RecreateSchema(client *typesense.Client) error {
	if _, err := client.Collection(ProductsCollection).Delete(context.TODO()); err != nil {
		var httpErr *typesense.HTTPError
		if !errors.As(err, &httpErr) || httpErr.Status != 404 {
			log.Println("[typesense] drop collection failed: ", err)
		}
	}
	return CreateSchema(client)
}
