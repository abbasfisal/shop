package util

import (
	"context"
	"github.com/typesense/typesense-go/v3/typesense/api"
	"log"
	"shop/infrastructure/database/typesenceclient"
)

type UpsertTypesenceProduct struct {
	ID    string
	Title string
	Slug  string
	Sku   string
}

// UpsertInTypesence upsert product in typesence search engine
func UpsertInTypesence(c context.Context, product UpsertTypesenceProduct) {

	doc := struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Slug  string `json:"slug"`
		Sku   string `json:"sku"`
	}{
		ID:    product.ID,
		Title: product.Title,
		Slug:  product.Slug,
		Sku:   product.Sku,
	}

	client := typesenceclient.GetTClient()
	if client == nil {
		log.Println("--- typesense client not initialized, skip product upsert")
		return
	}

	_, err := client.
		Collection("products").
		Documents().Upsert(c, doc, &api.DocumentIndexParameters{})

	if err != nil {
		log.Println("--- upsert product typesence failed:", err)
	}
	log.Println("--- upsert product typesence successfully")

}
