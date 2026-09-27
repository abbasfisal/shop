package util

import (
	"context"
	"log"
	"sync"

	"github.com/typesense/typesense-go/v3/typesense/api"
	"shop/infrastructure/database/typesenceclient"
)

// typesenseWG tracks fire-and-forget index writes so batch commands
// (e.g. `search:reindex`) can wait until every upsert/delete has really
// reached Typesense before the process exits. Without this wait the
// runtime kills in-flight goroutines on exit and the collection stays
// empty even though the command reported success.
var typesenseWG sync.WaitGroup

// UpsertInTypesenceAsync is the non-blocking twin of UpsertInTypesence for
// request paths (admin save, order flow): the HTTP response never waits
// for the search engine, but the write is registered in typesenseWG.
func UpsertInTypesenceAsync(c context.Context, product UpsertTypesenceProduct) {
	typesenseWG.Add(1)
	go func() {
		defer typesenseWG.Done()
		UpsertInTypesence(c, product)
	}()
}

// DeleteInTypesenceAsync is the non-blocking twin of DeleteInTypesence.
func DeleteInTypesenceAsync(c context.Context, productID string) {
	typesenseWG.Add(1)
	go func() {
		defer typesenseWG.Done()
		DeleteInTypesence(c, productID)
	}()
}

// WaitForTypesence blocks until every tracked async index write has
// finished. Batch commands MUST call it before exiting.
func WaitForTypesence() {
	typesenseWG.Wait()
}

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
