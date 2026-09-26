package util

import (
	"context"
	"log"

	"github.com/typesense/typesense-go/v3/typesense/api"
	"shop/infrastructure/database/typesenceclient"
)

// UpsertTypesenceProduct is the rich document indexed into Typesense
// (mirrors the products collection schema).
type UpsertTypesenceProduct struct {
	ID            string
	Title         string
	Slug          string
	Sku           string
	Description   string
	Category      string
	Brand         string
	OriginalPrice int64
	SalePrice     int64
	Discount      int64
	Stock         int64
	InStock       bool
	Status        string
}

// UpsertInTypesence upsert product in typesense search engine
func UpsertInTypesence(c context.Context, product UpsertTypesenceProduct) {
	doc := struct {
		ID            string `json:"id"`
		Title         string `json:"title"`
		Slug          string `json:"slug"`
		Sku           string `json:"sku"`
		Description   string `json:"description"`
		Category      string `json:"category"`
		Brand         string `json:"brand"`
		OriginalPrice int64  `json:"original_price"`
		SalePrice     int64  `json:"sale_price"`
		Discount      int64  `json:"discount"`
		Stock         int64  `json:"stock"`
		InStock       bool   `json:"in_stock"`
		Status        string `json:"status"`
	}{
		ID:            product.ID,
		Title:         product.Title,
		Slug:          product.Slug,
		Sku:           product.Sku,
		Description:   product.Description,
		Category:      product.Category,
		Brand:         product.Brand,
		OriginalPrice: product.OriginalPrice,
		SalePrice:     product.SalePrice,
		Discount:      product.Discount,
		Stock:         product.Stock,
		InStock:       product.InStock,
		Status:        product.Status,
	}

	client := typesenceclient.GetTClient()
	if client == nil {
		log.Println("--- typesense client not initialized, skip product upsert")
		return
	}

	_, err := client.
		Collection(typesenceclient.ProductsCollection).
		Documents().Upsert(c, doc, &api.DocumentIndexParameters{})

	if err != nil {
		log.Println("--- upsert product typesence failed:", err)
		return
	}
	log.Println("--- upsert product typesence successfully, id:", product.ID)
}

// DeleteInTypesence removes a product document from the search index
// (called when a product is unpublished or removed).
func DeleteInTypesence(c context.Context, productID string) {
	client := typesenceclient.GetTClient()
	if client == nil {
		log.Println("--- typesense client not initialized, skip product delete")
		return
	}

	_, err := client.
		Collection(typesenceclient.ProductsCollection).
		Document(productID).
		Delete(c)

	if err != nil {
		log.Println("--- delete product from typesence failed:", err)
		return
	}
	log.Println("--- delete product from typesence successfully, id:", productID)
}
